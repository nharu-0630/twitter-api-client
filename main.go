package main

import (
	"github.com/nharu-0630/twitter-api-client/cmd"
	"github.com/nharu-0630/twitter-api-client/tools"
)

func main() {
	tools.LoadEnv()

	seedScreenName := []string{"nharu_0630"}
	maxFollowersRequest := 1
	maxChildRequest := 1
	props := cmd.RandomUserFromSeedUsersProps{
		SeedScreenName:      seedScreenName,
		MaxFollowersRequest: maxFollowersRequest,
		MaxChildRequest:     maxChildRequest,
	}
	cmd := cmd.RandomUserFromSeedUsersCmd{
		Props: props,
	}
	cmd.Execute()
}
