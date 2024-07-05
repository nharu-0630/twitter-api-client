package cmd

import (
	"os"
	"time"

	"github.com/nharu-0630/twitter-api-client/api"
	"github.com/nharu-0630/twitter-api-client/model"
	"github.com/nharu-0630/twitter-api-client/tools"
	"go.uber.org/zap"
)

type GroupUsersProps struct {
	UserIDs                map[string][]string `yaml:"UserIDs"`
	MaxConversationRequest int                 `yaml:"MaxConversationRequest"`
	RetryOnGuestFail       bool                `yaml:"RetryOnGuestFail"`
	StatusUpdateSec        int                 `yaml:"StatusUpdateSec"`
}

func (props GroupUsersProps) Validate() {
	if props.UserIDs == nil {
		zap.L().Fatal("User ids are required")
	}
	if props.MaxConversationRequest < 1 {
		zap.L().Fatal("Max conversation request must be greater than 0")
	}
	if props.StatusUpdateSec < 1 {
		zap.L().Fatal("Status update sec must be greater than 0")
	}
}

type GroupUsersCmd struct {
	CmdName     string `yaml:"CmdName"`
	Props       GroupUsersProps
	GuestClient *api.Client
	Client      *api.Client
	TweetIDs    map[string]string
}

func (cmd *GroupUsersCmd) Execute() {
	cmd.GuestClient = api.NewClient(
		api.ClientConfig{
			IsGuestTokenEnabled: true,
		},
	)
	cmd.Client = api.NewClient(
		api.ClientConfig{
			IsGuestTokenEnabled: false,
			AuthToken:           os.Getenv("AUTH_TOKEN"),
			CsrfToken:           os.Getenv("CSRF_TOKEN"),
		},
	)
	cmd.TweetIDs = make(map[string]string)
	cmd.CmdName = ""
	for groupName := range cmd.Props.UserIDs {
		cmd.CmdName += groupName + "_"
	}
	cmd.CmdName += time.Now().Format("20060102150405")

	startDateTime := time.Now()
	zap.L().Info("Start of the process", zap.String("CmdName", cmd.CmdName))

	cmd.Props.Validate()

	ticker := time.NewTicker(time.Duration(cmd.Props.StatusUpdateSec) * time.Second)
	go func() {
		for range ticker.C {
			cmd.status("Status update")
		}
	}()

	for groupName, userIDs := range cmd.Props.UserIDs {
		for _, userID := range userIDs {
			cmd.getUserTweetsFromUserID(groupName, userID)
		}
	}

	defer func() {
		ticker.Stop()
	}()

	summary := map[string]interface{}{
		"Type":  "RandomUsers",
		"Props": cmd.Props,
		"Status": map[string]interface{}{
			"TweetCount":    len(cmd.TweetIDs),
			"StartDateTime": startDateTime,
			"EndDateTime":   time.Now(),
			"TotalSec":      time.Since(startDateTime).Seconds(),
			"SecPerTweet":   time.Since(startDateTime).Seconds() / float64(len(cmd.TweetIDs)),
		},
	}
	tools.Log(cmd.CmdName, []string{"Summary"}, summary, true)
	cmd.status("End of the process")
}

func (cmd *GroupUsersCmd) status(msg string) {
	zap.L().Info(msg, zap.String("CmdName", cmd.CmdName), zap.Int("TweetCount", len(cmd.TweetIDs)))
}

func (cmd *GroupUsersCmd) getUserTweetsFromUserID(groupName string, userID string) {
	zap.L().Debug("Get user tweets", zap.String("Group", groupName), zap.String("UserID", userID))
	var tweets []model.Tweet
	tweets, _, err := cmd.GuestClient.UserTweets(userID)
	if err != nil {
		zap.L().Error(err.Error())
		if err.Error() == "instruction not found" && cmd.Props.RetryOnGuestFail {
			tweets, _, err = cmd.GuestClient.UserTweets(userID)
			if err != nil {
				zap.L().Error(err.Error())
				return
			}
		} else {
			return
		}
	}
	if len(tweets) == 0 {
		return
	}
	tweetDetails := map[string]interface{}{}
	for _, tweet := range tweets {
		if _, exists := cmd.TweetIDs[tweet.RestID]; !exists {
			cmd.TweetIDs[tweet.RestID] = groupName + "_" + userID
			bottomCursor := ""
			var tweet model.Tweet
			var conversation []model.Tweet
			for i := 0; i < cmd.Props.MaxConversationRequest; i++ {
				resTweet, resConversation, cursor, err := cmd.Client.TweetDetail(tweet.RestID, bottomCursor)
				if err != nil {
					zap.L().Error(err.Error())
					break
				}
				if tweet.RestID == "" {
					tweet = resTweet
				}
				conversation = append(conversation, resConversation...)
				if cursor.IsAfterLast {
					break
				}
				bottomCursor = cursor.BottomCursor
			}
			tweetDetails[tweet.RestID] = map[string]interface{}{
				"Tweet":        tweet,
				"Conversation": conversation,
			}
		}
	}
	tools.Log(cmd.CmdName, []string{"Tweets", groupName, userID}, tweetDetails, false)
}
