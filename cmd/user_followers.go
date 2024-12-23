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

type UserFollowersProps struct {
	SeedScreenName      string `yaml:"seed_screen_name"`
	MaxFollowersRequest int    `yaml:"max_followers_request"`
	MaxChildRequest     int    ` yaml:"max_child_request"`
	MaxUserLimit        int    `yaml:"max_user_limit"`
	RetryOnGuestFail    bool   `yaml:"retry_on_guest_fail"`
}

func (props UserFollowersProps) Validate() {
	if props.SeedScreenName == "" {
		zap.L().Fatal("Seed screen name is required")
	}
	if props.MaxFollowersRequest < 1 {
		zap.L().Fatal("Max followers request must be greater than 0")
	}
	if props.MaxChildRequest < 1 {
		zap.L().Fatal("Max child request must be greater than 0")
	}
	if props.MaxUserLimit < 1 {
		zap.L().Fatal("Max user limit must be greater than 0")
	}
}

type UserFollowersCmd struct {
	CmdName          string
	Props            UserFollowersProps
	GuestClient      *api.Client
	Client           *api.Client
	UserIDs          map[string]string
	TweetIDs         map[string]string
	LeftChildRequest int
}

func (cmd *UserFollowersCmd) Execute() {
	zap.L().Info("Start of the process")
	cmd.Props.Validate()
	cmd.GuestClient = api.NewGuestClient()
	cmd.Client = api.NewAuthorizedClient(
		os.Getenv("AUTH_TOKEN"),
		os.Getenv("CSRF_TOKEN"),
	)

	startDateTime := time.Now()
	cmd.CmdName = cmd.Props.SeedScreenName + "_" + time.Now().Format("20060102150405")
	cmd.UserIDs = make(map[string]string)
	cmd.TweetIDs = make(map[string]string)
	cmd.LeftChildRequest = cmd.Props.MaxChildRequest

	ticker := util.NewStatusTicker()
	go func() {
		for range ticker.C {
			zap.L().Info("Status update", zap.String("CmdName", cmd.CmdName), zap.Int("UserCount", len(cmd.UserIDs)), zap.Int("TweetCount", len(cmd.TweetIDs)))
		}
	}()

	seedUserID := []string{}
	user, err := cmd.GuestClient.UserByScreenName(cmd.Props.SeedScreenName)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	util.Log(cmd.CmdName, []string{"User", user.RestID}, map[string]interface{}{"User": user}, false)
	seedUserID = append(seedUserID, user.RestID)
	cmd.UserIDs[user.RestID] = "ROOT"
	cmd.userTweetsExecute(user.RestID)
	util.LogOverwrite(cmd.CmdName, []string{"UserIDs"}, map[string]interface{}{"UserIDs": cmd.UserIDs}, false)

	cmd.usersExecute(seedUserID)

	defer func() {
		ticker.Stop()
	}()
	summary := map[string]interface{}{
		"Type":  "UserFollowers",
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
	util.Log(cmd.CmdName, []string{"Summary"}, summary, true)
	zap.L().Info("End of the process")
}

func (cmd *UserFollowersCmd) usersExecute(userIDs []string) {
	cmd.LeftChildRequest--
	childUserIDs := []string{}
	for _, userID := range userIDs {
		bottomCursor := ""
		for i := 0; i < cmd.Props.MaxFollowersRequest; i++ {
			followers, cursor, err := cmd.Client.Followers(userID, bottomCursor)
			if err != nil {
				zap.L().Error(err.Error())
				break
			}
			util.Log(cmd.CmdName, []string{"Followers", userID, strconv.Itoa(i)}, map[string]interface{}{"Followers": followers}, false)
			for _, follower := range followers {
				if _, exists := cmd.UserIDs[follower.RestID]; !exists {
					cmd.UserIDs[follower.RestID] = userID
					if !follower.Legacy.Protected {
						cmd.userTweetsExecute(follower.RestID)
						childUserIDs = append(childUserIDs, follower.RestID)
					}
					util.LogOverwrite(cmd.CmdName, []string{"UserIDs"}, map[string]interface{}{"UserIDs": cmd.UserIDs}, false)
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
		cmd.usersExecute(childUserIDs)
	}
}

func (cmd *UserFollowersCmd) userTweetsExecute(userID string) {
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
	util.Log(cmd.CmdName, []string{"Tweets", userID}, map[string]interface{}{"Tweets": tweets}, false)
}
