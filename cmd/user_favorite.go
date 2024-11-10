package cmd

import (
	"os"
	"strconv"
	"time"

	"github.com/nharu-0630/twitter-api-client/api"
	"github.com/nharu-0630/twitter-api-client/util"
	"go.uber.org/zap"
)

type UserFavoriteProps struct {
	UserID  string `yaml:"user_id"`
	UntilID string `yaml:"until_id"`
}

func (props UserFavoriteProps) Validate() {
	if props.UserID == "" {
		zap.L().Fatal("User id is required")
	}
}

type UserFavoriteCmd struct {
	CmdName  string
	Props    UserFavoriteProps
	Client   *api.Client
	TweetIDs map[string]string
}

func (cmd *UserFavoriteCmd) Execute() {
	zap.L().Info("Start of the process")
	cmd.Props.Validate()
	cmd.Client = api.NewAuthorizedClient(
		os.Getenv("AUTH_TOKEN"),
		os.Getenv("CSRF_TOKEN"),
	)

	startDateTime := time.Now()
	cmd.CmdName = cmd.Props.UserID + "_" + "Likes"
	cmd.TweetIDs = make(map[string]string)

	ticker := util.NewStatusTicker()
	go func() {
		for range ticker.C {
			zap.L().Info("Status update", zap.String("CmdName", cmd.CmdName), zap.Int("TweetCount", len(cmd.TweetIDs)))
		}
	}()

	bottomCursor := ""
	for i := 0; true; i++ {
		tweets, cursor, err := cmd.Client.Likes(cmd.Props.UserID, bottomCursor)
		if err != nil {
			zap.L().Fatal(err.Error())
			break
		}
		util.Log(cmd.CmdName, []string{"Tweets", cmd.Props.UserID, strconv.Itoa(i)}, map[string]interface{}{"Tweets": tweets}, false)
		for _, tweet := range tweets {
			if tweet.RestID == cmd.Props.UntilID {
				zap.L().Info("Reached until_id", zap.String("UntilID", cmd.Props.UntilID))
				return
			}
			cmd.TweetIDs[tweet.RestID] = ""
			cmd.Client.DownloadAllMedia(cmd.CmdName, tweet)
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
		"Type":  "UserFavorite",
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
	zap.L().Info("End of the process")
}
