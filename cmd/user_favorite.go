package cmd

import (
	"os"
	"strconv"
	"time"

	"github.com/nharu-0630/twitter-api-client/api"
	"github.com/nharu-0630/twitter-api-client/tools"
	"go.uber.org/zap"
)

type UserFavoriteProps struct {
	UserID          string `yaml:"UserID"`
	UntilID         string `yaml:"UntilID"`
	StatusUpdateSec int    `yaml:"StatusUpdateSec"`
}

func (props UserFavoriteProps) Validate() {
	if props.UserID == "" {
		zap.L().Fatal("User id is required")
	}
	if props.StatusUpdateSec < 1 {
		zap.L().Fatal("Status update sec must be greater than 0")
	}
}

type UserFavoriteCmd struct {
	CmdName  string
	Props    UserFavoriteProps
	Client   *api.Client
	TweetIDs map[string]string
}

func (cmd *UserFavoriteCmd) Execute() {
	cmd.Client = api.NewClient(
		api.ClientConfig{
			IsGuestTokenEnabled: false,
			AuthToken:           os.Getenv("AUTH_TOKEN"),
			CsrfToken:           os.Getenv("CSRF_TOKEN"),
		},
	)
	cmd.TweetIDs = make(map[string]string)
	cmd.CmdName = cmd.Props.UserID + "_" + "Likes"

	startDateTime := time.Now()
	zap.L().Info("Start of the process", zap.String("CmdName", cmd.CmdName))

	cmd.Props.Validate()

	ticker := time.NewTicker(time.Duration(cmd.Props.StatusUpdateSec) * time.Second)
	go func() {
		for range ticker.C {
			cmd.status("Status update")
		}
	}()

	bottomCursor := ""
	for i := 0; true; i++ {
		tweets, cursor, err := cmd.Client.Likes(cmd.Props.UserID, bottomCursor)
		if err != nil {
			zap.L().Fatal(err.Error())
			break
		}
		tools.Log(cmd.CmdName, []string{"Tweets", cmd.Props.UserID, strconv.Itoa(i)}, map[string]interface{}{"Tweets": tweets}, false)
		for _, tweet := range tweets {
			if tweet.RestID == cmd.Props.UntilID {
				cmd.status("UntilID found")
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
	tools.Log(cmd.CmdName, []string{"Summary"}, summary, true)
	cmd.status("End of the process")
}

func (cmd *UserFavoriteCmd) status(msg string) {
	zap.L().Info(msg, zap.String("CmdName", cmd.CmdName), zap.Int("TweetCount", len(cmd.TweetIDs)))
}
