package cmd

import (
	"flag"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/nharu-0630/twitter-api-client/api"
	"github.com/nharu-0630/twitter-api-client/model"
	"github.com/nharu-0630/twitter-api-client/tools"
	"go.uber.org/zap"
)

type UsersFollowersProps struct {
	CmdName             string
	SeedScreenName      string
	MaxFollowersRequest int
	MaxChildRequest     int
	MaxUserLimit        int
	RetryOnGuestFail    bool
	StatusUpdateSec     int
}

func UsersFollowersCmdFromFlag() UsersFollowersCmd {
	seedScreenName := flag.String("from", "", "シードとなるユーザーのscreen_name (必須)")
	maxFollowersRequest := flag.Int("req", math.MaxInt, "1ユーザあたりの最大フォロワーリクエスト数 指定しない場合は全てのフォロワーを取得")
	maxChildRequest := flag.Int("depth", 1, "シードとなるユーザからの最大深度 指定しない場合は1(シードとなるユーザのフォロワーのみ取得)")
	maxUserLimit := flag.Int("limit", math.MaxInt, "取得するユーザー数の上限 指定しない場合は全てのユーザを取得")
	retryOnGuestFail := flag.Bool("retry", false, "ゲストトークンでのリクエスト失敗時に認証済みトークンでリトライする")
	statusUpdateSec := flag.Int("watch", 600, "ステータスを更新する間隔(秒) 指定しない場合は10分ごとに更新")
	flag.Parse()
	props := UsersFollowersProps{
		SeedScreenName:      *seedScreenName,
		MaxFollowersRequest: *maxFollowersRequest,
		MaxChildRequest:     *maxChildRequest,
		MaxUserLimit:        *maxUserLimit,
		RetryOnGuestFail:    *retryOnGuestFail,
		StatusUpdateSec:     *statusUpdateSec,
	}
	cmd := UsersFollowersCmd{
		Props: props,
	}
	return cmd
}

func (props UsersFollowersProps) Validate() {
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
	if props.StatusUpdateSec < 1 {
		zap.L().Fatal("Status update sec must be greater than 0")
	}
}

type UsersFollowersCmd struct {
	Props            UsersFollowersProps
	GuestClient      *api.Client
	Client           *api.Client
	UserIDs          map[string]string
	TweetIDs         map[string]string
	LeftChildRequest int
}

func (cmd *UsersFollowersCmd) Execute() {
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
	cmd.Props.CmdName = cmd.Props.SeedScreenName + "_" + time.Now().Format("20060102150405")

	startDateTime := time.Now()
	zap.L().Info("Start of the process", zap.String("CmdName", cmd.Props.CmdName), zap.String("SeedScreenName", cmd.Props.SeedScreenName), zap.Int("MaxFollowersRequest", cmd.Props.MaxFollowersRequest), zap.Int("MaxChildRequest", cmd.Props.MaxChildRequest), zap.Int("MaxUserLimit", cmd.Props.MaxUserLimit), zap.Bool("RetryOnGuestFail", cmd.Props.RetryOnGuestFail), zap.Int("StatusUpdateSec", cmd.Props.StatusUpdateSec))

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
		log.Fatal(err)
	}
	tools.Log(cmd.Props.CmdName, []string{"User", user.RestID}, map[string]interface{}{"User": user}, false)
	seedUserID = append(seedUserID, user.RestID)
	cmd.UserIDs[user.RestID] = "ROOT"
	cmd.getUserTweetsFromUserID(user.RestID)
	tools.LogOverwrite(cmd.Props.CmdName, []string{"UserIDs"}, map[string]interface{}{"UserIDs": cmd.UserIDs}, false)

	cmd.getUsersFromUserIDs(seedUserID)

	defer func() {
		ticker.Stop()
	}()

	summary := map[string]interface{}{
		"Type":  "UsersFollowers",
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
	tools.Log(cmd.Props.CmdName, []string{"Summary"}, summary, true)
	cmd.status("End of the process")
}

func (cmd *UsersFollowersCmd) status(msg string) {
	zap.L().Info(msg, zap.String("CmdName", cmd.Props.CmdName), zap.Int("UserCount", len(cmd.UserIDs)), zap.Int("TweetCount", len(cmd.TweetIDs)))
}

func (cmd *UsersFollowersCmd) getUsersFromUserIDs(userIDs []string) {
	zap.L().Debug("Get users", zap.Int("UserCount", len(userIDs)))
	cmd.LeftChildRequest--
	childUserIDs := []string{}
	for _, userID := range userIDs {
		bottomCursor := ""
		for i := 0; i < cmd.Props.MaxFollowersRequest; i++ {
			followers, cursor, err := cmd.Client.Followers(userID, bottomCursor)
			if err != nil {
				log.Default().Println(err)
				break
			}
			tools.Log(cmd.Props.CmdName, []string{"Followers", userID, strconv.Itoa(i)}, map[string]interface{}{"Followers": followers}, false)
			for _, follower := range followers {
				if _, exists := cmd.UserIDs[follower.RestID]; !exists {
					cmd.UserIDs[follower.RestID] = userID
					if !follower.Legacy.Protected {
						cmd.getUserTweetsFromUserID(follower.RestID)
						// if tools.IsJapaneseUser(follower) {
						childUserIDs = append(childUserIDs, follower.RestID)
						// }
					}
					tools.LogOverwrite(cmd.Props.CmdName, []string{"UserIDs"}, map[string]interface{}{"UserIDs": cmd.UserIDs}, false)
				}
			}
			if len(cmd.UserIDs) > cmd.Props.MaxUserLimit {
				return
			}
			if cursor.BottomCursor == "" {
				break
			}
			bottomCursor = cursor.BottomCursor
		}
	}
	if cmd.LeftChildRequest > 0 {
		cmd.getUsersFromUserIDs(childUserIDs)
	}
}

func (cmd *UsersFollowersCmd) getUserTweetsFromUserID(userID string) {
	zap.L().Debug("Get user tweets", zap.String("UserID", userID))
	var tweets []model.Tweet
	tweets, _, err := cmd.GuestClient.UserTweets(userID)
	if err != nil {
		log.Default().Println(err)
		if err.Error() == "instruction not found" && cmd.Props.RetryOnGuestFail {
			tweets, _, err = cmd.GuestClient.UserTweets(userID)
			if err != nil {
				log.Default().Println(err)
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
	tools.Log(cmd.Props.CmdName, []string{"Tweets", userID}, map[string]interface{}{"Tweets": tweets}, false)
}
