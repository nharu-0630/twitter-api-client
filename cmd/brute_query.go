package cmd

import (
	"time"

	"github.com/nharu-0630/twitter-api-client/api"
	"github.com/nharu-0630/twitter-api-client/model"
	"github.com/nharu-0630/twitter-api-client/util"
	"go.uber.org/zap"
)

type BruteQueryProps struct {
	model.Config `yaml:",inline"`
	IntervalDays struct {
		PerQuery   int `yaml:"per_query"`
		PerRequest int `yaml:"per_request"`
	} `yaml:"interval_days"`
	Timestamp struct {
		Since time.Time `yaml:"since"`
		Until time.Time `yaml:"until"`
	} `yaml:"timestamp"`
	Queries []struct {
		Key        string   `yaml:"key"`
		RawQueries []string `yaml:"raw_queries"`
	} `yaml:"queries"`
}

type BruteQuery struct {
	Props       BruteQueryProps
	ClientsPipe *api.ClientsPipe
}

func (bq *BruteQuery) Execute() {
	zap.L().Info("検索クエリによるツイート取得を開始します", zap.String("config", bq.Props.Config.Name))
	bq.ClientsPipe = api.NewClientsPipe(bq.Props.Config)

	queryMap := make(map[string][]string)
	queryMapIdx := make(map[string]int)
	for _, query := range bq.Props.Queries {
		queryMap[query.Key] = append(queryMap[query.Key], query.RawQueries...)
		queryMapIdx[query.Key] = 0
	}

	for since := bq.Props.Timestamp.Since; since.Before(bq.Props.Timestamp.Until); since = since.AddDate(0, 0, bq.Props.IntervalDays.PerQuery) {
		for {
			flag := true
			for key := range queryMap {
				if queryMapIdx[key] >= len(queryMap[key]) {
					continue
				}
				flag = false
				rawQuery := queryMap[key][queryMapIdx[key]]

				for childSince := since; childSince.Before(bq.Props.Timestamp.Until); childSince = childSince.AddDate(0, 0, bq.Props.IntervalDays.PerRequest) {
					childUntil := childSince.AddDate(0, 0, bq.Props.IntervalDays.PerRequest)
					rawQuery = rawQuery + " until:" + childUntil.Format("2006-01-02") + " since:" + childSince.Format("2006-01-02")
					tweets := []model.Tweet{}
					bottomCursor := ""
					for {
						res, cursor, err := bq.ClientsPipe.SearchTimeline(rawQuery, bottomCursor)
						if err != nil {
							break
						}
						tweets = append(tweets, res...)
						if cursor.IsAfterLast {
							break
						}
						bottomCursor = cursor.BottomCursor
					}
					util.Log(bq.Props.Config.Name, []string{"Tweets", key}, map[string]interface{}{"Tweets": tweets}, false)
				}
				queryMapIdx[key]++
			}
			if flag {
				break
			}
		}
		queryMapIdx = make(map[string]int)
	}
	zap.L().Info("検索クエリによるツイート取得が完了しました", zap.String("config", bq.Props.Config.Name))
}
