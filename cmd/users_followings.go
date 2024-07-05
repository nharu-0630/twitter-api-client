package cmd

import (
	"os"
	"strconv"
	"time"

	"github.com/nharu-0630/twitter-api-client/api"
	"github.com/nharu-0630/twitter-api-client/model"
	"github.com/nharu-0630/twitter-api-client/tools"
	"go.uber.org/zap"
)

type UserFollowingsProps struct {
	SeedScreenName       string `yaml:"SeedScreenName"`
	MaxFollowingsRequest int    `yaml:"MaxFollowingsRequest"`
	MaxChildRequest      int    `yaml:"MaxChildRequest"`
	MaxUserLimit         int    `yaml:"MaxUserLimit"`
	RetryOnGuestFail     bool   `yaml:"RetryOnGuestFail"`
	StatusUpdateSec      int    `yaml:"StatusUpdateSec"`
}

func (props UserFollowingsProps) Validate() {
	if props.SeedScreenName == "" {
		zap.L().Fatal("Seed screen name is required")
	}
	if props.MaxFollowingsRequest < 1 {
		zap.L().Fatal("Max followings request must be greater than 0")
	}
	if props.MaxChildRequest < 1 {
		zap.L().Fatal("Max child request must be greater than 0")
	}
	if props.MaxUserLimit < 1 {
		zap.L().Fatal("Max user limit must be greater than 0")
	}
	if props.StatusUpdateSec < 1 {
		zap.L().Fatal("Status update sec must be greater than 0")
	}
}

type UserFollowingsCmd struct {
	CmdName          string
	Props            UserFollowingsProps
	GuestClient      *api.Client
	Client           *api.Client
	UserIDs          map[string]string
	TweetIDs         map[string]string
	LeftChildRequest int
}

func (cmd *UserFollowingsCmd) Execute() {
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
	cmd.UserIDs = make(map[string]string)
	cmd.TweetIDs = make(map[string]string)
	cmd.LeftChildRequest = cmd.Props.MaxChildRequest
	cmd.CmdName = cmd.Props.SeedScreenName + "_" + time.Now().Format("20060102150405")

	startDateTime := time.Now()
	zap.L().Info("Start of the process", zap.String("CmdName", cmd.CmdName))

	cmd.Props.Validate()

	ticker := time.NewTicker(time.Duration(cmd.Props.StatusUpdateSec) * time.Second)
	go func() {
		for range ticker.C {
			cmd.status("Status update")
		}
	}()

	seedUserID := []string{}
	user, err := cmd.GuestClient.UserByScreenName(cmd.Props.SeedScreenName)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	tools.Log(cmd.CmdName, []string{"User", user.RestID}, map[string]interface{}{"User": user}, false)
	seedUserID = append(seedUserID, user.RestID)
	cmd.UserIDs[user.RestID] = "ROOT"
	cmd.getUserTweetsFromUserID(user.RestID)
	tools.LogOverwrite(cmd.CmdName, []string{"UserIDs"}, map[string]interface{}{"UserIDs": cmd.UserIDs}, false)

	cmd.getUsersFromUserIDs(seedUserID)

	defer func() {
		ticker.Stop()
	}()

	summary := map[string]interface{}{
		"Type":  "UserFollowings",
		"Props": cmd.Props,
		"Status": map[string]interface{}{
			"UserCount":     len(cmd.UserIDs),
			"TweetCount":    len(cmd.TweetIDs),
			"StartDateTime": startDateTime,
			"EndDateTime":   time.Now(),
			"TotalSec":      time.Since(startDateTime).Seconds(),
			"SecPerUser":    time.Since(startDateTime).Seconds() / float64(len(cmd.UserIDs)),
			"SecPerTweet":   time.Since(startDateTime).Seconds() / float64(len(cmd.TweetIDs)),
		},
	}
	tools.Log(cmd.CmdName, []string{"Summary"}, summary, true)
	cmd.status("End of the process")
}

func (cmd *UserFollowingsCmd) status(msg string) {
	zap.L().Info(msg, zap.String("CmdName", cmd.CmdName), zap.Int("UserCount", len(cmd.UserIDs)), zap.Int("TweetCount", len(cmd.TweetIDs)))
}

func (cmd *UserFollowingsCmd) getUsersFromUserIDs(userIDs []string) {
	cmd.LeftChildRequest--
	childUserIDs := []string{}
	for _, userID := range userIDs {
		bottomCursor := ""
		for i := 0; i < cmd.Props.MaxFollowingsRequest; i++ {
			followings, cursor, err := cmd.Client.Following(userID, bottomCursor)
			if err != nil {
				zap.L().Error(err.Error())
				break
			}
			tools.Log(cmd.CmdName, []string{"Followings", userID, strconv.Itoa(i)}, map[string]interface{}{"Followings": followings}, false)
			for _, follower := range followings {
				if _, exists := cmd.UserIDs[follower.RestID]; !exists {
					cmd.UserIDs[follower.RestID] = userID
					if !follower.Legacy.Protected {
						cmd.getUserTweetsFromUserID(follower.RestID)
						childUserIDs = append(childUserIDs, follower.RestID)
					}
					tools.LogOverwrite(cmd.CmdName, []string{"UserIDs"}, map[string]interface{}{"UserIDs": cmd.UserIDs}, false)
				}
			}
			if len(cmd.UserIDs) > cmd.Props.MaxUserLimit {
				return
			}
			if cursor.IsAfterLast {
				break
			}
			bottomCursor = cursor.BottomCursor
		}
	}
	if cmd.LeftChildRequest > 0 {
		cmd.getUsersFromUserIDs(childUserIDs)
	}
}

func (cmd *UserFollowingsCmd) getUserTweetsFromUserID(userID string) {
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
	for _, tweet := range tweets {
		if _, exists := cmd.TweetIDs[tweet.RestID]; !exists {
			cmd.TweetIDs[tweet.RestID] = userID
		}
	}
	tools.Log(cmd.CmdName, []string{"Tweets", userID}, map[string]interface{}{"Tweets": tweets}, false)
}
