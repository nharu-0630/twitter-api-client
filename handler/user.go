package handler

import (
	"github.com/nharu-0630/twitter-api-client/api"
)

func ScreenNameToUserID(screenName string) (string, error) {
	client := api.NewClient(
		api.ClientConfig{
			IsGuestTokenEnabled: true,
		},
	)
	user, err := client.UserByScreenName(screenName)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}
