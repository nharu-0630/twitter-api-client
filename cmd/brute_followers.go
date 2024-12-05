package cmd

import (
	"time"

	"github.com/nharu-0630/twitter-api-client/api"
	"github.com/nharu-0630/twitter-api-client/model"
	"github.com/nharu-0630/twitter-api-client/util"
	"go.uber.org/zap"
)

type BruteFollowersProps struct {
	model.Config `yaml:",inline"`
	MaxLimit     struct {
		PerUser  int `yaml:"per_user"`
		PerTweet int `yaml:"per_tweet"`
	} `yaml:"max_limit"`
	Timestamp struct {
		Since time.Time `yaml:"since"`
		Until time.Time `yaml:"until"`
	} `yaml:"timestamp"`
	Queries []struct {
		Key         string   `yaml:"key"`
		ScreenNames []string `yaml:"screen_names"`
	}
}

type BruteFollowers struct {
	Props       BruteFollowersProps
	ClientsPipe *api.ClientsPipe
}

func (bf *BruteFollowers) Execute() {
	zap.L().Info("フォロワーのツイート取得を開始します", zap.String("config", bf.Props.Config.Name))
	bf.ClientsPipe = api.NewClientsPipe(bf.Props.Config)

	queryMap := make(map[string][]string)
	queryMapIdx := make(map[string]int)
	for _, query := range bf.Props.Queries {
		queryMap[query.Key] = append(queryMap[query.Key], query.ScreenNames...)
		queryMapIdx[query.Key] = 0
	}

	queryMapFollowers := make(map[string][]model.User)
	queryMapFollowersIdx := make(map[string]int)
	for {
		flag := true
		for key := range queryMap {
			if queryMapIdx[key] >= len(queryMap[key]) {
				continue
			}
			flag = false
			screenName := queryMap[key][queryMapIdx[key]]
			user, err := bf.ClientsPipe.UserByScreenName(screenName)
			if err != nil {
				zap.L().Error("ユーザー情報の取得に失敗しました", zap.String("screen_name", screenName), zap.Error(err))
				continue
			}
			util.Log(bf.Props.Config.Name, []string{"User", screenName}, map[string]interface{}{"User": user}, false)

			followers := []model.User{}
			bottomCursor := ""
			for {
				res, cursor, err := bf.ClientsPipe.Followers(user.ID, bottomCursor)
				if err != nil {
					zap.L().Error("フォロワーの取得に失敗しました", zap.String("screen_name", screenName), zap.Error(err))
					break
				}
				followers = append(followers, res...)
				if cursor.IsAfterLast {
					break
				}
				if len(followers) >= bf.Props.MaxLimit.PerUser {
					break
				}
			}
			util.Log(bf.Props.Config.Name, []string{"Followers", key, screenName}, map[string]interface{}{"Followers": followers}, false)

			queryMapIdx[key]++

			queryMapFollowers[key] = append(queryMapFollowers[key], followers...)
			queryMapFollowersIdx[key] = 0
		}
		if flag {
			break
		}
	}

	for {
		flag := true
		for key := range queryMap {
			if queryMapFollowersIdx[key] >= len(queryMapFollowers[key]) {
				continue
			}
			flag = false
			user := queryMapFollowers[key][queryMapFollowersIdx[key]]

			tweets := []model.Tweet{}
			bottomCursor := ""
			for {
				rawQuery := "from:" + user.Legacy.ScreenName
				rawQuery = rawQuery + " until:" + bf.Props.Timestamp.Until.Format("2006-01-02") + " since:" + bf.Props.Timestamp.Since.Format("2006-01-02")
				res, cursor, err := bf.ClientsPipe.SearchTimeline(rawQuery, bottomCursor)
				if err != nil {
					zap.L().Error("ツイートの取得に失敗しました", zap.String("screen_name", user.Legacy.ScreenName), zap.Error(err))
					break
				}
				tweets = append(tweets, res...)
				if cursor.IsAfterLast {
					break
				}
				if len(tweets) >= bf.Props.MaxLimit.PerTweet {
					break
				}
			}
			util.Log(bf.Props.Config.Name, []string{"Tweets", key, user.Legacy.ScreenName}, map[string]interface{}{"Tweets": tweets}, false)

			queryMapFollowersIdx[key]++
		}
		if flag {
			break
		}
	}

	zap.L().Info("フォロワーのツイート取得が完了しました", zap.String("config", bf.Props.Config.Name))
}
