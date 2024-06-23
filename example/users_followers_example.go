package example

import (
	"math"
	"time"

	"github.com/nharu-0630/twitter-api-client/cmd"
)

func NewUsersFollowersCmdPoliticalParty(screenName string) cmd.UsersFollowersCmd {
	props := cmd.UsersFollowersProps{
		SeedScreenName:      screenName,
		MaxFollowersRequest: math.MaxInt,
		MaxChildRequest:     1,
		MaxUserLimit:        3000,
		StatusUpdateSec:     int(time.Minute.Seconds() * 10),
	}
	cmd := cmd.UsersFollowersCmd{
		Props: props,
	}
	return cmd
}
