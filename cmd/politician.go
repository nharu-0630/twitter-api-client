package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/nharu-0630/twitter-api-client/api"
	"github.com/nharu-0630/twitter-api-client/model"
	"github.com/nharu-0630/twitter-api-client/util"
	"go.uber.org/zap"
)

type PoliticianProps struct {
	Politicians []Politician `yaml:"politicians"`
}

type Politician struct {
	BlockName      string `yaml:"block_name"`
	PartyName      string `yaml:"party_name"`
	CandidateKanji string `yaml:"candidate_kanji"`
	CandidateKana  string `yaml:"candidate_kana"`
	XID            string `yaml:"x_id"`
}

type PoliticianCmd struct {
	CmdName string
	Props   PoliticianProps
	Client  *api.Client
}

func (cmd *PoliticianCmd) Execute() {
	zap.L().Info("Start of the process")
	cmd.Client = api.NewClient(
		api.ClientConfig{
			IsGuestTokenEnabled: false,
			AuthToken:           os.Getenv("AUTH_TOKEN"),
			CsrfToken:           os.Getenv("CSRF_TOKEN"),
		},
	)
	cmd.CmdName = "Politician_" + time.Now().Format("20060102150405")

	partyMap := make(map[string][]Politician)
	for _, politician := range cmd.Props.Politicians {
		if politician.XID == "" {
			continue
		}
		partyMap[politician.PartyName] = append(partyMap[politician.PartyName], politician)
	}

	for partyName, politicians := range partyMap {
		zap.L().Info("Processing party", zap.String("PartyName", partyName))
		for _, politician := range politicians {
			candidateName := strings.ReplaceAll(politician.CandidateKanji, "　", "")
			zap.L().Info("Processing politician", zap.String("CandidateName", candidateName))

			user, err := cmd.Client.UserByScreenName(politician.XID)
			if err != nil {
				zap.L().Error("Failed to get user", zap.Error(err))
				continue
			}
			util.Log(cmd.CmdName, []string{"User", politician.PartyName, candidateName}, map[string]interface{}{"User": user}, false)

			var tweets []model.Tweet
			bottomCursor := ""
			for i := 0; i < 5; i++ {
				resTweets, cursor, err := cmd.Client.UserTweets(user.RestID, bottomCursor)
				if err != nil {
					zap.L().Error("Failed to get tweets", zap.Error(err))
					break
				}
				tweets = append(tweets, resTweets...)
				if cursor.IsAfterLast {
					break
				}
				bottomCursor = cursor.BottomCursor
			}
			util.Log(cmd.CmdName, []string{"Tweets", politician.PartyName, candidateName}, map[string]interface{}{"Tweets": tweets}, false)
		}
	}

	zap.L().Info("End of the process")
}
