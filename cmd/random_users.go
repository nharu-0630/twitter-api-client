package cmd

import (
	"log"
	"os"
	"strconv"

	"github.com/nharu-0630/twitter-api-client/api"
	"github.com/nharu-0630/twitter-api-client/tools"
)

type RandomUsersProps struct {
	CmdName             string
	SeedScreenName      string
	MaxFollowersRequest int
	MaxChildRequest     int
	MaxUserLimit        int
}

type RandomUsersCmd struct {
	Props            RandomUsersProps
	GuestClient      *api.Client
	Client           *api.Client
	UserIDs          map[string]string
	LeftChildRequest int
}

func (cmd *RandomUsersCmd) Execute() {
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
	cmd.UserIDs = make(map[string]string)
	seedUserID := []string{}
	user, err := cmd.GuestClient.UserByScreenName(cmd.Props.SeedScreenName)
	if err != nil {
		log.Fatal(err)
	}
	cmd.Props.CmdName = user.Legacy.Name
	tools.Log(cmd.Props.CmdName, []string{"User", user.RestID}, map[string]interface{}{"User": user})
	seedUserID = append(seedUserID, user.RestID)
	cmd.UserIDs[user.RestID] = "ROOT"
	cmd.getUserTweetsFromUserID(user.RestID)

	tools.LogOverwrite(cmd.Props.CmdName, []string{"UserIDs"}, map[string]interface{}{"UserIDs": cmd.UserIDs})
	cmd.getUsersFromUserIDs(seedUserID)
}

func (cmd *RandomUsersCmd) getUsersFromUserIDs(userIDs []string) {
	cmd.LeftChildRequest--
	childUserIDs := []string{}
	for _, userID := range userIDs {
		bottomCursor := ""
		for i := 0; i < cmd.Props.MaxFollowersRequest; i++ {
			followers, cursor, err := cmd.Client.Followers(userID, bottomCursor)
			if err != nil {
				log.Fatal(err)
			}
			tools.Log(cmd.Props.CmdName, []string{"Followers", userID, strconv.Itoa(i)}, map[string]interface{}{"Followers": followers})
			for _, follower := range followers {
				if _, exists := cmd.UserIDs[follower.RestID]; !exists {
					cmd.UserIDs[follower.RestID] = userID
					if !follower.Legacy.Protected {
						cmd.getUserTweetsFromUserID(follower.RestID)
						if tools.IsJapaneseUser(follower) {
							childUserIDs = append(childUserIDs, follower.RestID)
						}
					}
				}
			}
			tools.LogOverwrite(cmd.Props.CmdName, []string{"UserIDs"}, map[string]interface{}{"UserIDs": cmd.UserIDs})
			if len(cmd.UserIDs) > cmd.Props.MaxUserLimit {
				return
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

func (cmd *RandomUsersCmd) getUserTweetsFromUserID(userID string) {
	tweets, _, err := cmd.GuestClient.UserTweets(userID)
	if err != nil {
		log.Default().Println(err)
		return
	}
	tools.Log(cmd.Props.CmdName, []string{"Tweets", userID}, map[string]interface{}{"Tweets": tweets})
}
