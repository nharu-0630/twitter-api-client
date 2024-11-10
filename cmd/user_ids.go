package cmd

import (
	"os"
	"strconv"
	"time"

	"github.com/nharu-0630/twitter-api-client/api"
	"github.com/nharu-0630/twitter-api-client/model"
	"github.com/nharu-0630/twitter-api-client/util"
	"go.uber.org/zap"
)

type UserIDsProps struct {
	UserIDs                []string `yaml:"user_ids"`
	MaxUserTweetsRequest   int      `yaml:"max_user_tweets_request"`
	MaxConversationRequest int      `yaml:"max_conversation_request"`
}

func (props UserIDsProps) Validate() {
	if len(props.UserIDs) == 0 {
		zap.L().Fatal("User ids are required")
	}
	if props.MaxUserTweetsRequest < 1 {
		zap.L().Fatal("Max user tweets request must be greater than 0")
	}
	if props.MaxConversationRequest < 0 {
		zap.L().Fatal("Max conversation request must be greater than or equal to 0")
	}
}

type UserIDsCmd struct {
	CmdName  string
	Props    UserIDsProps
	Client   *api.Client
	TweetIDs map[string]string
}

func (cmd *UserIDsCmd) Execute() {
	zap.L().Info("Start of the process")
	cmd.Props.Validate()
	cmd.Client = api.NewAuthorizedClient(
		os.Getenv("AUTH_TOKEN"),
		os.Getenv("CSRF_TOKEN"),
	)

	for _, userID := range cmd.Props.UserIDs {
		startDateTime := time.Now()
		cmd.CmdName = userID + "_" + time.Now().Format("20060102150405")
		cmd.TweetIDs = make(map[string]string)

		ticker := util.NewStatusTicker()
		go func() {
			for range ticker.C {
				zap.L().Info("Status update", zap.String("CmdName", cmd.CmdName), zap.Int("TweetCount", len(cmd.TweetIDs)))
			}
		}()

		bottomCursor := ""
		for i := 0; i < cmd.Props.MaxUserTweetsRequest; i++ {
			tweets, cursor, err := cmd.Client.UserTweets(userID, bottomCursor)
			if err != nil {
				zap.L().Fatal(err.Error())
				break
			}
			for _, tweet := range tweets {
				cmd.TweetIDs[tweet.RestID] = ""
			}
			if cmd.Props.MaxConversationRequest == 0 {
				util.Log(cmd.CmdName, []string{"Tweet", userID, strconv.Itoa(i)}, map[string]interface{}{"Tweets": tweets}, false)
			} else {
				conversation := cmd.conversationExecute(tweets)
				util.Log(cmd.CmdName, []string{"Tweet", userID, strconv.Itoa(i)}, map[string]interface{}{"Tweets": conversation}, false)
			}
			if cursor.IsAfterLast {
				break
			}
			bottomCursor = cursor.BottomCursor
		}

		defer func() {
			ticker.Stop()
		}()
		summary := map[string]interface{}{
			"Type":  "UserIDs",
			"Props": cmd.Props,
			"Status": map[string]interface{}{
				"TweetCount":    len(cmd.TweetIDs),
				"StartDateTime": startDateTime,
				"EndDateTime":   time.Now(),
				"TotalSec":      time.Since(startDateTime).Seconds(),
				"SecPerTweet":   time.Since(startDateTime).Seconds() / float64(len(cmd.TweetIDs)),
			},
		}
		util.Log(cmd.CmdName, []string{"Summary"}, summary, true)
	}
	zap.L().Info("End of the process")
}

func (cmd *UserIDsCmd) conversationExecute(tweets []model.Tweet) map[string]interface{} {
	tweetDetails := map[string]interface{}{}
	for _, tweet := range tweets {
		tweetDetails[tweet.RestID] = map[string]interface{}{
			"Tweet": tweet,
		}
		bottomCursor := ""
		var conversation []model.Tweet
		for i := 0; i < cmd.Props.MaxConversationRequest; i++ {
			_, resConversation, cursor, err := cmd.Client.TweetDetail(tweet.RestID, bottomCursor)
			if err != nil {
				zap.L().Error(err.Error())
				break
			}
			conversation = append(conversation, resConversation...)
			if cursor.IsAfterLast {
				break
			}
			bottomCursor = cursor.BottomCursor
		}
		tweetDetails[tweet.RestID].(map[string]interface{})["Conversation"] = conversation
	}
	return tweetDetails
}
