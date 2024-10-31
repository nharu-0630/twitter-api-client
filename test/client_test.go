package test

import (
	"log"
	"testing"

	"github.com/nharu-0630/twitter-api-client/api"
)

func TestTweetResultByRestIdWithGuest(t *testing.T) {
	client := api.NewClient(api.ClientConfig{
		IsGuestTokenEnabled: true,
	})

	result, err := client.TweetResultByRestId("1762646776916635958")
	if err != nil {
		t.Fatal(err)
	}
	if result.RestID == "" {
		t.Fatal("result.RestID is empty")
	}
	log.Default().Println(result)
}

func TestUserByScreenNameWithGuest(t *testing.T) {
	client := api.NewClient(api.ClientConfig{
		IsGuestTokenEnabled: true,
	})

	user, err := client.UserByScreenName("nharu_0630")
	if err != nil {
		t.Fatal(err)
	}
	if user.ID == "" {
		t.Fatal("user.ID is empty")
	}
	log.Default().Println(user)
}

func TestUserTweetsWithGuest1(t *testing.T) {
	client := api.NewClient(api.ClientConfig{
		IsGuestTokenEnabled: true,
	})

	tweets, cursor, err := client.UserTweets("1003084799592972288")
	if err != nil {
		t.Fatal(err)
	}
	if len(tweets) == 0 {
		t.Fatal("tweets is empty")
	}
	log.Default().Println(tweets)
	log.Default().Println(cursor)
}

func TestUserTweetsWithGuest2(t *testing.T) {
	client := api.NewClient(api.ClientConfig{
		IsGuestTokenEnabled: true,
	})

	_, _, err := client.UserTweets("1679423214500585473")
	if err.Error() != "instruction not found" {
		t.Fatal(err)
	}
}

// func TestTweetDetail(t *testing.T) {
// 	tools.LoadEnv()

// 	client := api.NewClient(api.ClientConfig{
// 		IsGuestTokenEnabled: false,
// 		AuthToken:           os.Getenv("AUTH_TOKEN"),
// 		CsrfToken:           os.Getenv("CSRF_TOKEN"),
// 	})

// 	tweet, conversation, cursor, err := client.TweetDetail("1750775577437724978")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	log.Default().Println(tweet)
// 	log.Default().Println(conversation)
// 	log.Default().Println(cursor)
// }

// func TestFollowing(t *testing.T) {
// 	tools.LoadEnv()

// 	client := api.NewClient(api.ClientConfig{
// 		IsGuestTokenEnabled: false,
// 		AuthToken:           os.Getenv("AUTH_TOKEN"),
// 		CsrfToken:           os.Getenv("CSRF_TOKEN"),
// 	})

// 	users, cursor, err := client.Following("1003084799592972288")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(users) == 0 {
// 		t.Fatal("users is empty")
// 	}
// 	log.Default().Println(users)
// 	log.Default().Println(cursor)
// }

// func TestFollowers(t *testing.T) {
// 	tools.LoadEnv()

// 	client := api.NewClient(api.ClientConfig{
// 		IsGuestTokenEnabled: false,
// 		AuthToken:           os.Getenv("AUTH_TOKEN"),
// 		CsrfToken:           os.Getenv("CSRF_TOKEN"),
// 	})

// 	users, cursor, err := client.Followers("1003084799592972288")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(users) == 0 {
// 		t.Fatal("users is empty")
// 	}
// 	log.Default().Println(users)
// 	log.Default().Println(cursor)
// }
