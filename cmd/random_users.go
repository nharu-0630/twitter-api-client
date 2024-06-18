package cmd

import (
	"log"
	"os"
	"strconv"

	"github.com/nharu-0630/twitter-api-client/api"
	"github.com/nharu-0630/twitter-api-client/model"
	"github.com/nharu-0630/twitter-api-client/tools"
)

type RandomUserFromSeedUsersProps struct {
	SeedScreenName      []string
	MaxFollowersRequest int
	MaxChildRequest     int
}

type RandomUserFromSeedUsersCmd struct {
	Props            RandomUserFromSeedUsersProps
	GuestClient      *api.Client
	Client           *api.Client
	Users            []model.User
	LeftChildRequest int
}

func (cmd *RandomUserFromSeedUsersCmd) Execute() {
	cmd.GuestClient = api.NewClient(
		api.ClientConfig{
			IsGuestTokenEnabled: true,
		},
	)
	cmd.Client = api.NewClient(
		api.ClientConfig{
			IsGuestTokenEnabled: false,
			AuthToken:           os.Getenv("AUTH_TOKEN"),
			CsrfToken:           os.Getenv("CSRF_TOKEN"),
		},
	)
	cmd.LeftChildRequest = cmd.Props.MaxChildRequest

	seedUserID := []string{}
	for _, screenName := range cmd.Props.SeedScreenName {
		user, err := cmd.GuestClient.UserByScreenName(screenName)
		if err != nil {
			log.Fatal(err)
		}
		tools.Log("RandomUsersFromSeedUsers", []string{"User", user.ID}, map[string]interface{}{"User": user})
		seedUserID = append(seedUserID, user.ID)
		cmd.Users = append(cmd.Users, user)
	}
	cmd.getUsersFromUserIDs(seedUserID)
}

func (cmd *RandomUserFromSeedUsersCmd) getUsersFromUserIDs(userIDs []string) {
	cmd.LeftChildRequest--
	childUserIDs := []string{}
	for _, userID := range userIDs {
		bottomCursor := ""
		for i := 0; i < cmd.Props.MaxFollowersRequest; i++ {
			followers, cursor, err := cmd.Client.Followers(userID, bottomCursor)
			if err != nil {
				log.Fatal(err)
			}
			tools.Log("RandomUsersFromSeedUsers", []string{"Followers", userID, strconv.Itoa(i)}, map[string]interface{}{"Followers": followers})
			cmd.Users = append(cmd.Users, followers...)
			for _, follower := range followers {
				childUserIDs = append(childUserIDs, follower.ID)
			}

			if cursor.BottomCursor == "" {
				break
			}
			bottomCursor = cursor.BottomCursor
		}
	}
	if cmd.LeftChildRequest > 0 {
		cmd.getUsersFromUserIDs(childUserIDs)
	}
}
