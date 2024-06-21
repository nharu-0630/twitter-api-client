package api

import (
	"encoding/json"
	"errors"
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

func ParseTimelineEntriesTweets(res map[string]interface{}) ([]model.Tweet, model.Cursor, error) {
	if res["user"] == nil {
		return nil, model.Cursor{}, errors.New("user not found")
	}
	if res["user"].(map[string]interface{})["result"] == nil {
		return nil, model.Cursor{}, errors.New("result not found")
	}
	if res["user"].(map[string]interface{})["result"].(map[string]interface{})["timeline_v2"] == nil {
		return nil, model.Cursor{}, errors.New("timeline_v2 not found")
	}
	if res["user"].(map[string]interface{})["result"].(map[string]interface{})["timeline_v2"].(map[string]interface{})["timeline"] == nil {
		return nil, model.Cursor{}, errors.New("timeline not found")
	}
	if res["user"].(map[string]interface{})["result"].(map[string]interface{})["timeline_v2"].(map[string]interface{})["timeline"].(map[string]interface{})["instructions"] == nil {
		return nil, model.Cursor{}, errors.New("instructions not found")
	}
	instructions := res["user"].(map[string]interface{})["result"].(map[string]interface{})["timeline_v2"].(map[string]interface{})["timeline"].(map[string]interface{})["instructions"].([]interface{})
	var data map[string]interface{}
	for _, instruction := range instructions {
		instructionType := instruction.(map[string]interface{})["type"]
		if instructionType == "TimelineAddEntries" {
			data = instruction.(map[string]interface{})
			break
		}
	}
	if data == nil {
		return nil, model.Cursor{}, errors.New("instruction not found")
	}

	tweets := make([]model.Tweet, 0)
	resCursor := model.Cursor{}
	entries := data["entries"].([]interface{})
	for _, entry := range entries {
		content := entry.(map[string]interface{})["content"]
		entryID := entry.(map[string]interface{})["entryId"].(string)
		entryType := content.(map[string]interface{})["entryType"]
		if entryType == "TimelineTimelineItem" {
			if strings.HasPrefix(entryID, "tweet-") {
				tweetResults := content.(map[string]interface{})["itemContent"].(map[string]interface{})["tweet_results"].(map[string]interface{})
				if tweetResults["result"] == nil {
					continue
				}
				tweet, err := ParseTweet(tweetResults["result"].(map[string]interface{}))
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

func ParseTimelineEntriesBookmarksTweets(res map[string]interface{}) ([]model.Tweet, model.Cursor, error) {
	if res["bookmark_timeline_v2"] == nil {
		return nil, model.Cursor{}, errors.New("bookmark_timeline_v2 not found")
	}
	if res["bookmark_timeline_v2"].(map[string]interface{})["timeline"] == nil {
		return nil, model.Cursor{}, errors.New("timeline not found")
	}
	if res["bookmark_timeline_v2"].(map[string]interface{})["timeline"].(map[string]interface{})["instructions"] == nil {
		return nil, model.Cursor{}, errors.New("instructions not found")
	}
	instructions := res["bookmark_timeline_v2"].(map[string]interface{})["timeline"].(map[string]interface{})["instructions"].([]interface{})
	var data map[string]interface{}
	for _, instruction := range instructions {
		instructionType := instruction.(map[string]interface{})["type"]
		if instructionType == "TimelineAddEntries" {
			data = instruction.(map[string]interface{})
			break
		}
	}
	if data == nil {
		return nil, model.Cursor{}, errors.New("instruction not found")
	}

	tweets := make([]model.Tweet, 0)
	resCursor := model.Cursor{}
	entries := data["entries"].([]interface{})
	for _, entry := range entries {
		content := entry.(map[string]interface{})["content"]
		entryID := entry.(map[string]interface{})["entryId"].(string)
		entryType := content.(map[string]interface{})["entryType"]
		if entryType == "TimelineTimelineItem" {
			if strings.HasPrefix(entryID, "tweet-") {
				tweetResults := content.(map[string]interface{})["itemContent"].(map[string]interface{})["tweet_results"].(map[string]interface{})
				if tweetResults["result"] == nil {
					continue
				}
				tweet, err := ParseTweet(tweetResults["result"].(map[string]interface{}))
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

func ParseTimelineEntriesTweetsWithInjections(res map[string]interface{}) (model.Tweet, []model.Tweet, model.Cursor, error) {
	if res["threaded_conversation_with_injections_v2"] == nil {
		return model.Tweet{}, nil, model.Cursor{}, errors.New("threaded_conversation_with_injections_v2 not found")
	}
	if res["threaded_conversation_with_injections_v2"].(map[string]interface{})["instructions"] == nil {
		return model.Tweet{}, nil, model.Cursor{}, errors.New("instructions not found")
	}
	instructions := res["threaded_conversation_with_injections_v2"].(map[string]interface{})["instructions"].([]interface{})
	var data map[string]interface{}
	for _, instruction := range instructions {
		instructionType := instruction.(map[string]interface{})["type"]
		if instructionType == "TimelineAddEntries" {
			data = instruction.(map[string]interface{})
			break
		}
	}
	if data == nil {
		return model.Tweet{}, nil, model.Cursor{}, errors.New("instruction not found")
	}

	resTweet := model.Tweet{}
	conversation := make([]model.Tweet, 0)
	resCursor := model.Cursor{}
	entries := data["entries"].([]interface{})
	for _, entry := range entries {
		content := entry.(map[string]interface{})["content"]
		entryID := entry.(map[string]interface{})["entryId"].(string)
		entryType := content.(map[string]interface{})["entryType"]
		if entryType == "TimelineTimelineItem" || entryType == "TimelineTimelineModule" {
			if strings.HasPrefix(entryID, "tweet-") {
				tweetResults := content.(map[string]interface{})["itemContent"].(map[string]interface{})["tweet_results"].(map[string]interface{})
				if tweetResults["result"] == nil {
					continue
				}
				tweet, err := ParseTweet(tweetResults["result"].(map[string]interface{}))
				if err != nil {
					continue
				}
				resTweet = tweet
			}
			if strings.HasPrefix(entryID, "conversationthread-") {
				items := content.(map[string]interface{})["items"].([]interface{})
				for _, item := range items {
					itemEntryID := item.(map[string]interface{})["entryId"].(string)
					if !strings.Contains(itemEntryID, "tweet-") {
						continue
					}
					tweetResults := item.(map[string]interface{})["item"].(map[string]interface{})["itemContent"].(map[string]interface{})["tweet_results"].(map[string]interface{})
					if tweetResults["result"] == nil {
						continue
					}
					tweet, err := ParseTweet(tweetResults["result"].(map[string]interface{}))
					if err != nil {
						continue
					}
					conversation = append(conversation, tweet)
				}
			}
			if strings.HasPrefix(entryID, "cursor-top") {
				resCursor.TopCursor = content.(map[string]interface{})["itemContent"].(map[string]interface{})["value"].(string)
			} else if strings.HasPrefix(entryID, "cursor-bottom") {
				resCursor.BottomCursor = content.(map[string]interface{})["itemContent"].(map[string]interface{})["value"].(string)
			}
		}
	}
	return resTweet, conversation, resCursor, nil
}

func ParseTimelineEntriesUsers(res map[string]interface{}) ([]model.User, model.Cursor, error) {
	if res["user"] == nil {
		return nil, model.Cursor{}, errors.New("user not found")
	}
	if res["user"].(map[string]interface{})["result"] == nil {
		return nil, model.Cursor{}, errors.New("result not found")
	}
	if res["user"].(map[string]interface{})["result"].(map[string]interface{})["timeline"] == nil {
		return nil, model.Cursor{}, errors.New("timeline not found")
	}
	if res["user"].(map[string]interface{})["result"].(map[string]interface{})["timeline"].(map[string]interface{})["timeline"] == nil {
		return nil, model.Cursor{}, errors.New("timeline not found")
	}
	if res["user"].(map[string]interface{})["result"].(map[string]interface{})["timeline"].(map[string]interface{})["timeline"].(map[string]interface{})["instructions"] == nil {
		return nil, model.Cursor{}, errors.New("instructions not found")
	}
	instructions := res["user"].(map[string]interface{})["result"].(map[string]interface{})["timeline"].(map[string]interface{})["timeline"].(map[string]interface{})["instructions"].([]interface{})
	data := instructions[len(instructions)-1].(map[string]interface{})

	users := make([]model.User, 0)
	resCursor := model.Cursor{}
	entries := data["entries"].([]interface{})
	for _, entry := range entries {
		content := entry.(map[string]interface{})["content"]
		entryID := entry.(map[string]interface{})["entryId"].(string)
		if strings.HasPrefix(entryID, "user-") {
			userResults := content.(map[string]interface{})["itemContent"].(map[string]interface{})["user_results"].(map[string]interface{})
			if userResults["result"] == nil {
				continue
			}
			user, err := ParseUser(userResults["result"].(map[string]interface{}))
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
