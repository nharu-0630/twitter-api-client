package example

import (
	"math"
	"time"

	"github.com/nharu-0630/twitter-api-client/cmd"
)

func NewRandomUsersCmdPoliticalParty(screenName string) cmd.RandomUsersCmd {
	props := cmd.RandomUsersProps{
		SeedScreenName:      screenName,
		MaxFollowersRequest: math.MaxInt,
		MaxChildRequest:     1,
		MaxUserLimit:        3000,
		StatusUpdateSec:     int(time.Minute.Seconds() * 10),
	}
	cmd := cmd.RandomUsersCmd{
		Props: props,
	}
	return cmd
}

func NewRandomUsersCmdHead(screenName string) cmd.RandomUsersCmd {
	props := cmd.RandomUsersProps{
		SeedScreenName:      screenName,
		MaxFollowersRequest: math.MaxInt,
		MaxChildRequest:     1,
		MaxUserLimit:        3,
		StatusUpdateSec:     int(time.Minute.Seconds() * 1),
	}
	cmd := cmd.RandomUsersCmd{
		Props: props,
	}
	return cmd
}
