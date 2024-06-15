package cmd

import (
	"log"
	"os"

	"github.com/nharu-0630/twitter-api-client/api"
	"github.com/nharu-0630/twitter-api-client/handler"
)

type RandomUserFromSeedUsersProps struct {
	SeedScreenName      []string
	MaxFollowersRequest int
}

func RandomUsersFromSeedUsers(props RandomUserFromSeedUsersProps) {
	seedUserID := []string{}
	for _, screenName := range props.SeedScreenName {
		userID, err := handler.ScreenNameToUserID(screenName)
		if err != nil {
			log.Fatal(err)
		}
		seedUserID = append(seedUserID, userID)
	}
	log.Println(seedUserID)
	client := api.NewClient(
		api.ClientConfig{
			IsGuestTokenEnabled: false,
			AuthToken:           os.Getenv("AUTH_TOKEN"),
			CsrfToken:           os.Getenv("CSRF_TOKEN"),
		},
	)
	for _, userID := range seedUserID {
		for i := 0; i < props.MaxFollowersRequest; i++ {
			followers, cursor, err := client.Followers(userID)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(followers)
			log.Println(cursor)
		}
	}
}
