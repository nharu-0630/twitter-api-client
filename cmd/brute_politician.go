package cmd

import (
	"time"

	"github.com/nharu-0630/twitter-api-client/api"
	"github.com/nharu-0630/twitter-api-client/model"
	"github.com/nharu-0630/twitter-api-client/util"
)

type BrutePoliticianProps struct {
	model.Config   `yaml:",inline"`
	Politicians    []Politician `yaml:"politicians"`
	IntervalDays   int          `yaml:"interval_days"`
	SinceTimestamp time.Time    `yaml:"since_timestamp"`
	UntilTimestamp time.Time    `yaml:"until_timestamp"`
}

type BrutePolitician struct {
	Props       BrutePoliticianProps
	ClientsPipe *api.ClientsPipe
}

func (bp *BrutePolitician) Execute() {
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

	for partyName := range partyMap {
		if partyMapIdx[partyName] >= len(partyMap[partyName]) {
			continue
		}
		politician := partyMap[partyName][partyMapIdx[partyName]]

		user, err := bp.ClientsPipe.UserByScreenName(politician.XID)
		if err != nil {
			continue
		}
		util.Log(bp.Props.Config.Name, []string{"User", politician.PartyName, politician.CandidateKanji}, map[string]interface{}{"User": user}, false)

		currentTimestamp := bp.Props.SinceTimestamp
		for currentTimestamp.Before(bp.Props.UntilTimestamp) {

			currentTimestamp = currentTimestamp.AddDate(0, 0, bp.Props.IntervalDays)
		}

		partyMapIdx[partyName]++
	}
}
