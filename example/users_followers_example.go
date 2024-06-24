package example

import (
	"math"
	"time"

	"github.com/nharu-0630/twitter-api-client/cmd"
)

func NewUserFollowersCmdPoliticalParty(screenName string) cmd.UserFollowersCmd {
	props := cmd.UserFollowersProps{
		SeedScreenName:      screenName,
		MaxFollowersRequest: math.MaxInt,
		MaxChildRequest:     1,
		MaxUserLimit:        3000,
		StatusUpdateSec:     int(time.Minute.Seconds() * 10),
	}
	cmd := cmd.UserFollowersCmd{
		Props: props,
	}
	return cmd
}
