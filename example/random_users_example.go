package example

import (
	"math"

	"github.com/nharu-0630/twitter-api-client/cmd"
)

func NewRandomUsersCmdPoliticalParty(screenName string) cmd.RandomUsersCmd {
	props := cmd.RandomUsersProps{
		SeedScreenName:      screenName,
		MaxFollowersRequest: math.MaxInt,
		MaxChildRequest:     1,
		MaxUserLimit:        10000,
	}
	cmd := cmd.RandomUsersCmd{
		Props: props,
	}
	return cmd
}
