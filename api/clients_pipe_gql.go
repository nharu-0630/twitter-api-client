package api

import (
	"github.com/nharu-0630/twitter-api-client/model"
)

func (cp *ClientsPipe) TweetResultByRestId(tweetID string) (model.Tweet, error) {
	client := cp.operation("TweetResultByRestId")
	return client.TweetResultByRestId(tweetID)
}

func (cp *ClientsPipe) UserByScreenName(screenName string) (model.User, error) {
	client := cp.operation("UserByScreenName")
	return client.UserByScreenName(screenName)
}

func (cp *ClientsPipe) UserTweets(userID string, cursor ...string) ([]model.Tweet, model.Cursor, error) {
	client := cp.operation("UserTweets")
	return client.UserTweets(userID, cursor...)
}

func (cp *ClientsPipe) TweetDetail(focalTweetID string, cursor ...string) (model.Tweet, []model.Tweet, model.Cursor, error) {
	client := cp.operation("TweetDetail")
	return client.TweetDetail(focalTweetID, cursor...)
}

func (cp *ClientsPipe) Following(userID string, cursor ...string) ([]model.User, model.Cursor, error) {
	client := cp.operation("Following")
	return client.Following(userID, cursor...)
}

func (cp *ClientsPipe) Followers(userID string, cursor ...string) ([]model.User, model.Cursor, error) {
	client := cp.operation("Followers")
	return client.Followers(userID, cursor...)
}

func (cp *ClientsPipe) SearchTimeline(rawQuery string, cursor ...string) ([]model.Tweet, model.Cursor, error) {
	client := cp.operation("SearchTimeline")
	return client.SearchTimeline(rawQuery, cursor...)
}
