package api

import (
	"encoding/json"
	"strings"

	"github.com/nharu-0630/twitter-api-client/model"
)

func ParseTweet(data map[string]interface{}) (model.Tweet, error) {
	encodedResult, err := json.Marshal(data)
	if err != nil {
		return model.Tweet{}, err
	}
	var tweet model.Tweet
	err = json.Unmarshal(encodedResult, &tweet)
	if err != nil {
		return model.Tweet{}, err
	}
	return tweet, nil
}

func ParseUser(data map[string]interface{}) (model.User, error) {
	encodedResult, err := json.Marshal(data)
	if err != nil {
		return model.User{}, err
	}
	var user model.User
	err = json.Unmarshal(encodedResult, &user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func ParseTimelineEntriesTweets(data map[string]interface{}) ([]model.Tweet, model.Cursor, error) {
	tweets := make([]model.Tweet, 0)
	resCursor := model.Cursor{}
	entries := data["entries"].([]interface{})
	for _, entry := range entries {
		content := entry.(map[string]interface{})["content"]
		entryID := entry.(map[string]interface{})["entryId"].(string)
		entryType := content.(map[string]interface{})["entryType"]
		if entryType == "TimelineTimelineItem" {
			if strings.HasPrefix(entryID, "tweet-") {
				tweet, err := ParseTweet(content.(map[string]interface{})["itemContent"].(map[string]interface{})["tweet_results"].(map[string]interface{})["result"].(map[string]interface{}))
				if err != nil {
					continue
				}
				tweets = append(tweets, tweet)
			}
			if strings.HasPrefix(entryID, "cursor-top") {
				resCursor.TopCursor = content.(map[string]interface{})["value"].(string)
			} else if strings.HasPrefix(entryID, "cursor-bottom") {
				resCursor.BottomCursor = content.(map[string]interface{})["value"].(string)
			}
		} else if entryType == "TimelineTimelineCursor" {
			if strings.HasPrefix(entryID, "cursor-top") {
				resCursor.TopCursor = content.(map[string]interface{})["value"].(string)
			} else if strings.HasPrefix(entryID, "cursor-bottom") {
				resCursor.BottomCursor = content.(map[string]interface{})["value"].(string)
			}
		}
	}
	return tweets, resCursor, nil
}

func ParseTimelineEntriesUsers(data map[string]interface{}) ([]model.User, model.Cursor, error) {
	users := make([]model.User, 0)
	resCursor := model.Cursor{}
	entries := data["entries"].([]interface{})
	for _, entry := range entries {
		content := entry.(map[string]interface{})["content"]
		entryID := entry.(map[string]interface{})["entryId"].(string)
		if strings.HasPrefix(entryID, "user-") {
			user, err := ParseUser(content.(map[string]interface{})["itemContent"].(map[string]interface{})["user_results"].(map[string]interface{})["result"].(map[string]interface{}))
			if err != nil {
				continue
			}
			users = append(users, user)
		}
		if strings.HasPrefix(entryID, "cursor-top") {
			resCursor.TopCursor = content.(map[string]interface{})["value"].(string)
		} else if strings.HasPrefix(entryID, "cursor-bottom") {
			resCursor.BottomCursor = content.(map[string]interface{})["value"].(string)
		}
	}
	return users, resCursor, nil
}
