package cmd

import (
	"flag"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/nharu-0630/twitter-api-client/api"
	"github.com/nharu-0630/twitter-api-client/tools"
	"go.uber.org/zap"
)

type UserIDProps struct {
	UserID               string
	MaxUserTweetsRequest int
	StatusUpdateSec      int
}

type UserIDCmd struct {
	CmdName  string
	Props    UserIDProps
	Client   *api.Client
	TweetIDs map[string]string
}

func UserIDCmdFromFlag() UserIDCmd {
	userID := flag.String("from", "", "対象となるユーザのuser_id (必須)")
	maxUserTweetsRequest := flag.Int("req", math.MaxInt, "最大ユーザツイートリクエスト数 指定しない場合は全てのツイートを取得")
	statusUpdateSec := flag.Int("watch", 600, "ステータスを更新する間隔(秒) 指定しない場合は10分ごとに更新")
	flag.Parse()
	props := UserIDProps{
		UserID:               *userID,
		MaxUserTweetsRequest: *maxUserTweetsRequest,
		StatusUpdateSec:      *statusUpdateSec,
	}
	cmd := UserIDCmd{
		Props: props,
	}
	return cmd
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
	zap.L().Info("Start of the process", zap.String("CmdName", cmd.CmdName), zap.String("UserID", cmd.Props.UserID), zap.Int("StatusUpdateSec", cmd.Props.StatusUpdateSec))

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
		if cursor.BottomCursor == "" || (cursor.BottomCursor != "" && len(tweets) == 0) {
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
