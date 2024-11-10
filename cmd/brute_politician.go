package cmd

import (
	"time"

	"github.com/nharu-0630/twitter-api-client/api"
	"github.com/nharu-0630/twitter-api-client/model"
	"github.com/nharu-0630/twitter-api-client/util"
	"go.uber.org/zap"
)

type BrutePoliticianProps struct {
	model.Config `yaml:",inline"`
	Politicians  []struct {
		BlockName      string `yaml:"block_name"`      // 政治家の所属ブロック名
		PartyName      string `yaml:"party_name"`      // 政治家の所属政党名
		CandidateKanji string `yaml:"candidate_kanji"` // 政治家の漢字氏名
		CandidateKana  string `yaml:"candidate_kana"`  // 政治家のカナ氏名
		XID            string `yaml:"x_id"`            // 政治家のTwitter ID
	} `yaml:"politicians"`
	RequestIntervalDays int       `yaml:"request_interval_days"` // 検索リクエストの絞り込み期間の日数
	IntervalDays        int       `yaml:"interval_days"`         // 同一ユーザに連続して検索リクエストを送信する間隔の日数
	SinceTimestamp      time.Time `yaml:"since_timestamp"`       // 検索リクエストの開始日時
	UntilTimestamp      time.Time `yaml:"until_timestamp"`       // 検索リクエストの終了日時
}

type BrutePolitician struct {
	Props       BrutePoliticianProps
	ClientsPipe *api.ClientsPipe
}

func (bp *BrutePolitician) Execute() {
	zap.L().Info("政治家のツイートを取得します", zap.String("config", bp.Props.Config.Name))
	bp.ClientsPipe = api.NewClientsPipe(bp.Props.Config)

	partyMap := make(map[string][]Politician)
	partyMapIdx := make(map[string]int)
	for _, politician := range bp.Props.Politicians {
		if politician.XID == "" {
			continue
		}
		partyMap[politician.PartyName] = append(partyMap[politician.PartyName], politician)
		partyMapIdx[politician.PartyName] = 0
	}

	sinceTimestamp := bp.Props.SinceTimestamp
	for sinceTimestamp.Before(bp.Props.UntilTimestamp) {
		untilTimestamp := sinceTimestamp.AddDate(0, 0, bp.Props.IntervalDays)
		zap.L().Debug("対象の日時範囲", zap.Time("since_timestamp", sinceTimestamp), zap.Time("until_timestamp", untilTimestamp))
		for {
			flag := true
			for partyName := range partyMap {
				if partyMapIdx[partyName] >= len(partyMap[partyName]) {
					continue
				}
				flag = false
				politician := partyMap[partyName][partyMapIdx[partyName]]
				zap.L().Debug("検索リクエストの対象政治家", zap.String("party_name", politician.PartyName), zap.String("candidate_kanji", politician.CandidateKanji))

				childSinceTimestamp := sinceTimestamp
				for childSinceTimestamp.Before(untilTimestamp) {
					childUntilTimestamp := childSinceTimestamp.AddDate(0, 0, bp.Props.RequestIntervalDays)
					zap.L().Debug("検索リクエストの日時範囲", zap.Time("since_timestamp", childSinceTimestamp), zap.Time("until_timestamp", childUntilTimestamp))
					rawQuery := "(from:" + politician.XID + ") until:" + childUntilTimestamp.Format("2006-01-02") + " since:" + childSinceTimestamp.Format("2006-01-02")
					var resTweets []model.Tweet
					bottomCursor := ""
					for {
						tweets, nextCursor, err := bp.ClientsPipe.SearchTimeline(rawQuery, bottomCursor)
						if err != nil {
							break
						}
						resTweets = append(resTweets, tweets...)
						if nextCursor.IsAfterLast {
							break
						}
						bottomCursor = nextCursor.BottomCursor
					}
					util.Log(bp.Props.Config.Name, []string{"Tweets", politician.PartyName, politician.CandidateKanji}, map[string]interface{}{"Tweets": resTweets}, false)
					childSinceTimestamp = childUntilTimestamp
				}
				partyMapIdx[partyName]++
			}
			if flag {
				break
			}
		}
		sinceTimestamp = untilTimestamp
		partyMapIdx = make(map[string]int)
	}

	zap.L().Info("政治家のツイート取得が完了しました", zap.String("config", bp.Props.Config.Name))
}
