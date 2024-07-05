package cmd

import (
	"os"
	"strconv"
	"time"

	"github.com/nharu-0630/twitter-api-client/api"
	"github.com/nharu-0630/twitter-api-client/tools"
	"go.uber.org/zap"
)

type UserIDProps struct {
	UserID               string `yaml:"UserID"`
	MaxUserTweetsRequest int    `yaml:"MaxUserTweetsRequest"`
	StatusUpdateSec      int    `yaml:"StatusUpdateSec"`
}

func (props UserIDProps) Validate() {
	if props.UserID == "" {
		zap.L().Fatal("User id is required")
	}
	if props.MaxUserTweetsRequest < 1 {
		zap.L().Fatal("Max user tweets request must be greater than 0")
	}
	if props.StatusUpdateSec < 1 {
		zap.L().Fatal("Status update sec must be greater than 0")
	}
}

type UserIDCmd struct {
	CmdName  string
	Props    UserIDProps
	Client   *api.Client
	TweetIDs map[string]string
}

func (cmd *UserIDCmd) Execute() {
	cmd.Client = api.NewClient(
		api.ClientConfig{
			IsGuestTokenEnabled: false,
			AuthToken:           os.Getenv("AUTH_TOKEN"),
			CsrfToken:           os.Getenv("CSRF_TOKEN"),
		},
	)
	cmd.TweetIDs = make(map[string]string)
	cmd.CmdName = cmd.Props.UserID + "_" + time.Now().Format("20060102150405")

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
	for i := 0; i < cmd.Props.MaxUserTweetsRequest; i++ {
		tweets, cursor, err := cmd.Client.UserTweets(cmd.Props.UserID, bottomCursor)
		if err != nil {
			zap.L().Fatal(err.Error())
		}
		tools.Log(cmd.CmdName, []string{"Tweet", cmd.Props.UserID, strconv.Itoa(i)}, map[string]interface{}{"Tweets": tweets}, false)
		if cursor.IsAfterLast {
			break
		}
		bottomCursor = cursor.BottomCursor
	}

	defer func() {
		ticker.Stop()
	}()

	summary := map[string]interface{}{
		"Type":  "UserID",
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

func (cmd *UserIDCmd) status(msg string) {
	zap.L().Info(msg, zap.String("CmdName", cmd.CmdName), zap.Int("TweetCount", len(cmd.TweetIDs)))
}
