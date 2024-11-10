package api

import (
	"errors"

	"github.com/nharu-0630/twitter-api-client/model"
	"github.com/nharu-0630/twitter-api-client/util"
	"go.uber.org/zap"
)

func (c *Client) TweetResultByRestId(tweetID string) (model.Tweet, error) {
	zap.L().Debug("TweetResultByRestId", zap.String("tweetID", tweetID))
	variables := map[string]interface{}{"tweetId": tweetID, "includePromotedContent": true, "withBirdwatchNotes": true, "withVoice": true, "withCommunity": true}
	features := map[string]interface{}{"creator_subscriptions_tweet_preview_api_enabled": true, "communities_web_enable_tweet_community_results_fetch": true, "c9s_tweet_anatomy_moderator_badge_enabled": true, "articles_preview_enabled": true, "tweetypie_unmention_optimization_enabled": true, "responsive_web_edit_tweet_api_enabled": true, "graphql_is_translatable_rweb_tweet_is_translatable_enabled": true, "view_counts_everywhere_api_enabled": true, "longform_notetweets_consumption_enabled": true, "responsive_web_twitter_article_tweet_consumption_enabled": true, "tweet_awards_web_tipping_enabled": false, "creator_subscriptions_quote_tweet_preview_enabled": false, "freedom_of_speech_not_reach_fetch_enabled": true, "standardized_nudges_misinfo": true, "tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true, "tweet_with_visibility_results_prefer_gql_media_interstitial_enabled": true, "rweb_video_timestamps_enabled": true, "longform_notetweets_rich_text_read_enabled": true, "longform_notetweets_inline_media_enabled": true, "rweb_tipjar_consumption_enabled": true, "responsive_web_graphql_exclude_directive_enabled": true, "verified_phone_label_enabled": false, "responsive_web_graphql_skip_user_profile_image_extensions_enabled": false, "responsive_web_graphql_timeline_navigation_enabled": true, "responsive_web_enhance_cards_enabled": false}
	params := map[string]interface{}{
		"variables": variables,
		"features":  features,
	}
	res, err := c.gql("GET", "7xflPyRiUxGVbJd4uWmbfg", "TweetResultByRestId", params)
	util.LogRaw([]string{"TweetResultByRestId", tweetID}, res, false)
	if err != nil {
		return model.Tweet{}, err
	}
	if res["tweetResult"] == nil {
		return model.Tweet{}, errors.New("tweetResult not found")
	}
	if res["tweetResult"].(map[string]interface{})["result"] == nil {
		return model.Tweet{}, errors.New("result not found")
	}
	tweet, err := ParseTweet(res["tweetResult"].(map[string]interface{})["result"].(map[string]interface{}))
	if err != nil {
		return model.Tweet{}, err
	}
	return tweet, nil
}

func (c *Client) UserByScreenName(screenName string) (model.User, error) {
	zap.L().Debug("UserByScreenName", zap.String("screenName", screenName))
	variables := map[string]interface{}{"screen_name": screenName, "withSafetyModeUserFields": true}
	features := map[string]interface{}{"hidden_profile_likes_enabled": true, "hidden_profile_subscriptions_enabled": true, "rweb_tipjar_consumption_enabled": true, "responsive_web_graphql_exclude_directive_enabled": true, "verified_phone_label_enabled": false, "subscriptions_verification_info_is_identity_verified_enabled": true, "subscriptions_verification_info_verified_since_enabled": true, "highlights_tweets_tab_ui_enabled": true, "responsive_web_twitter_article_notes_tab_enabled": true, "creator_subscriptions_tweet_preview_api_enabled": true, "responsive_web_graphql_skip_user_profile_image_extensions_enabled": false, "responsive_web_graphql_timeline_navigation_enabled": true}
	fieldToggles := map[string]interface{}{"withAuxiliaryUserLabels": false}
	params := map[string]interface{}{
		"variables":    variables,
		"features":     features,
		"fieldToggles": fieldToggles,
	}
	res, err := c.gql("GET", "qW5u-DAuXpMEG0zA1F7UGQ", "UserByScreenName", params)
	util.LogRaw([]string{"UserByScreenName", screenName}, res, false)
	if err != nil {
		return model.User{}, err
	}
	if res["user"] == nil {
		return model.User{}, errors.New("user not found")
	}
	if res["user"].(map[string]interface{})["result"] == nil {
		return model.User{}, errors.New("result not found")
	}
	user, err := ParseUser(res["user"].(map[string]interface{})["result"].(map[string]interface{}))
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (c *Client) UserTweets(userID string, cursor ...string) ([]model.Tweet, model.Cursor, error) {
	zap.L().Debug("UserTweets", zap.String("userID", userID))
	variables := map[string]interface{}{"userId": userID, "count": 20, "includePromotedContent": true, "withQuickPromoteEligibilityTweetFields": true, "withVoice": true, "withV2Timeline": true}
	if len(cursor) > 0 {
		if cursor[0] != "" {
			variables["cursor"] = cursor[0]
		}
	}
	features := map[string]interface{}{"rweb_tipjar_consumption_enabled": true, "responsive_web_graphql_exclude_directive_enabled": true, "verified_phone_label_enabled": false, "creator_subscriptions_tweet_preview_api_enabled": true, "responsive_web_graphql_timeline_navigation_enabled": true, "responsive_web_graphql_skip_user_profile_image_extensions_enabled": false, "communities_web_enable_tweet_community_results_fetch": true, "c9s_tweet_anatomy_moderator_badge_enabled": true, "articles_preview_enabled": true, "tweetypie_unmention_optimization_enabled": true, "responsive_web_edit_tweet_api_enabled": true, "graphql_is_translatable_rweb_tweet_is_translatable_enabled": true, "view_counts_everywhere_api_enabled": true, "longform_notetweets_consumption_enabled": true, "responsive_web_twitter_article_tweet_consumption_enabled": true, "tweet_awards_web_tipping_enabled": false, "creator_subscriptions_quote_tweet_preview_enabled": false, "freedom_of_speech_not_reach_fetch_enabled": true, "standardized_nudges_misinfo": true, "tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true, "tweet_with_visibility_results_prefer_gql_media_interstitial_enabled": true, "rweb_video_timestamps_enabled": true, "longform_notetweets_rich_text_read_enabled": true, "longform_notetweets_inline_media_enabled": true, "responsive_web_enhance_cards_enabled": false}
	fieldToggles := map[string]interface{}{"withArticlePlainText": false}
	params := map[string]interface{}{
		"variables":    variables,
		"features":     features,
		"fieldToggles": fieldToggles,
	}
	res, err := c.gql("GET", "9zyyd1hebl7oNWIPdA8HRw", "UserTweets", params)
	util.LogRaw([]string{"UserTweets", userID}, res, false)
	if err != nil {
		return nil, model.Cursor{}, err
	}
	return ParseTimelineEntriesTweets(res)
}

func (c *Client) TweetDetail(focalTweetID string, cursor ...string) (model.Tweet, []model.Tweet, model.Cursor, error) {
	zap.L().Debug("TweetDetail", zap.String("focalTweetID", focalTweetID))
	variables := map[string]interface{}{"focalTweetId": focalTweetID, "with_rux_injections": false, "includePromotedContent": true, "withCommunity": true, "withQuickPromoteEligibilityTweetFields": true, "withBirdwatchNotes": true, "withVoice": true, "withV2Timeline": true}
	if len(cursor) > 0 {
		if cursor[0] != "" {
			variables["cursor"] = cursor[0]
		}
	}
	features := map[string]interface{}{"rweb_tipjar_consumption_enabled": true, "responsive_web_graphql_exclude_directive_enabled": true, "verified_phone_label_enabled": false, "creator_subscriptions_tweet_preview_api_enabled": true, "responsive_web_graphql_timeline_navigation_enabled": true, "responsive_web_graphql_skip_user_profile_image_extensions_enabled": false, "communities_web_enable_tweet_community_results_fetch": true, "c9s_tweet_anatomy_moderator_badge_enabled": true, "articles_preview_enabled": true, "tweetypie_unmention_optimization_enabled": true, "responsive_web_edit_tweet_api_enabled": true, "graphql_is_translatable_rweb_tweet_is_translatable_enabled": true, "view_counts_everywhere_api_enabled": true, "longform_notetweets_consumption_enabled": true, "responsive_web_twitter_article_tweet_consumption_enabled": true, "tweet_awards_web_tipping_enabled": false, "creator_subscriptions_quote_tweet_preview_enabled": false, "freedom_of_speech_not_reach_fetch_enabled": true, "standardized_nudges_misinfo": true, "tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true, "rweb_video_timestamps_enabled": true, "longform_notetweets_rich_text_read_enabled": true, "longform_notetweets_inline_media_enabled": true, "responsive_web_enhance_cards_enabled": false}
	fieldToggles := map[string]interface{}{"withArticleRichContentState": true, "withArticlePlainText": false, "withGrokAnalyze": false}
	params := map[string]interface{}{
		"variables":    variables,
		"features":     features,
		"fieldToggles": fieldToggles,
	}
	res, err := c.gql("GET", "VwKJcAd7zqlBOitPLUrB8A", "TweetDetail", params)
	util.LogRaw([]string{"TweetDetail", focalTweetID}, res, false)
	if err != nil {
		return model.Tweet{}, nil, model.Cursor{}, err
	}
	return ParseTimelineEntriesTweetsWithInjections(res)
}

func (c *Client) Following(userID string, cursor ...string) ([]model.User, model.Cursor, error) {
	zap.L().Debug("Following", zap.String("userID", userID))
	variables := map[string]interface{}{"userId": userID, "count": 20, "includePromotedContent": false}
	if len(cursor) > 0 {
		if cursor[0] != "" {
			variables["cursor"] = cursor[0]
		}
	}
	features := map[string]interface{}{"rweb_tipjar_consumption_enabled": true, "responsive_web_graphql_exclude_directive_enabled": true, "verified_phone_label_enabled": false, "creator_subscriptions_tweet_preview_api_enabled": true, "responsive_web_graphql_timeline_navigation_enabled": true, "responsive_web_graphql_skip_user_profile_image_extensions_enabled": false, "communities_web_enable_tweet_community_results_fetch": true, "c9s_tweet_anatomy_moderator_badge_enabled": true, "articles_preview_enabled": true, "tweetypie_unmention_optimization_enabled": true, "responsive_web_edit_tweet_api_enabled": true, "graphql_is_translatable_rweb_tweet_is_translatable_enabled": true, "view_counts_everywhere_api_enabled": true, "longform_notetweets_consumption_enabled": true, "responsive_web_twitter_article_tweet_consumption_enabled": true, "tweet_awards_web_tipping_enabled": false, "creator_subscriptions_quote_tweet_preview_enabled": false, "freedom_of_speech_not_reach_fetch_enabled": true, "standardized_nudges_misinfo": true, "tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true, "rweb_video_timestamps_enabled": true, "longform_notetweets_rich_text_read_enabled": true, "longform_notetweets_inline_media_enabled": true, "responsive_web_enhance_cards_enabled": false}
	params := map[string]interface{}{
		"variables": variables,
		"features":  features,
	}
	res, err := c.gql("GET", "7FEKOPNAvxWASt6v9gfCXw", "Following", params)
	util.LogRaw([]string{"Following", userID}, res, false)
	if err != nil {
		return nil, model.Cursor{}, err
	}
	return ParseTimelineEntriesUsers(res)
}

func (c *Client) Followers(userID string, cursor ...string) ([]model.User, model.Cursor, error) {
	zap.L().Debug("Followers", zap.String("userID", userID))
	variables := map[string]interface{}{"userId": userID, "count": 20, "includePromotedContent": false}
	if len(cursor) > 0 {
		if cursor[0] != "" {
			variables["cursor"] = cursor[0]
		}
	}
	features := map[string]interface{}{"rweb_tipjar_consumption_enabled": true, "responsive_web_graphql_exclude_directive_enabled": true, "verified_phone_label_enabled": false, "creator_subscriptions_tweet_preview_api_enabled": true, "responsive_web_graphql_timeline_navigation_enabled": true, "responsive_web_graphql_skip_user_profile_image_extensions_enabled": false, "communities_web_enable_tweet_community_results_fetch": true, "c9s_tweet_anatomy_moderator_badge_enabled": true, "articles_preview_enabled": true, "tweetypie_unmention_optimization_enabled": true, "responsive_web_edit_tweet_api_enabled": true, "graphql_is_translatable_rweb_tweet_is_translatable_enabled": true, "view_counts_everywhere_api_enabled": true, "longform_notetweets_consumption_enabled": true, "responsive_web_twitter_article_tweet_consumption_enabled": true, "tweet_awards_web_tipping_enabled": false, "creator_subscriptions_quote_tweet_preview_enabled": false, "freedom_of_speech_not_reach_fetch_enabled": true, "standardized_nudges_misinfo": true, "tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true, "rweb_video_timestamps_enabled": true, "longform_notetweets_rich_text_read_enabled": true, "longform_notetweets_inline_media_enabled": true, "responsive_web_enhance_cards_enabled": false}
	params := map[string]interface{}{
		"variables": variables,
		"features":  features,
	}
	res, err := c.gql("GET", "DMcBoZkXf9axSfV2XND0Ig", "Followers", params)
	util.LogRaw([]string{"Followers", userID}, res, false)
	if err != nil {
		return nil, model.Cursor{}, err
	}
	return ParseTimelineEntriesUsers(res)
}

func (c *Client) Likes(userID string, cursor ...string) ([]model.Tweet, model.Cursor, error) {
	if c.isGuest {
		return nil, model.Cursor{}, errors.New("Likes API is not available with guest token")
	}
	zap.L().Debug("Likes", zap.String("userID", userID))
	variables := map[string]interface{}{"userId": userID, "count": 20, "includePromotedContent": false, "withClientEventToken": false, "withBirdwatchNotes": false, "withVoice": true, "withV2Timeline": true}
	if len(cursor) > 0 {
		if cursor[0] != "" {
			variables["cursor"] = cursor[0]
		}
	}
	features := map[string]interface{}{"rweb_tipjar_consumption_enabled": true, "responsive_web_graphql_exclude_directive_enabled": true, "verified_phone_label_enabled": false, "creator_subscriptions_tweet_preview_api_enabled": true, "responsive_web_graphql_timeline_navigation_enabled": true, "responsive_web_graphql_skip_user_profile_image_extensions_enabled": false, "communities_web_enable_tweet_community_results_fetch": true, "c9s_tweet_anatomy_moderator_badge_enabled": true, "articles_preview_enabled": true, "tweetypie_unmention_optimization_enabled": true, "responsive_web_edit_tweet_api_enabled": true, "graphql_is_translatable_rweb_tweet_is_translatable_enabled": true, "view_counts_everywhere_api_enabled": true, "longform_notetweets_consumption_enabled": true, "responsive_web_twitter_article_tweet_consumption_enabled": true, "tweet_awards_web_tipping_enabled": false, "creator_subscriptions_quote_tweet_preview_enabled": false, "freedom_of_speech_not_reach_fetch_enabled": true, "standardized_nudges_misinfo": true, "tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true, "tweet_with_visibility_results_prefer_gql_media_interstitial_enabled": true, "rweb_video_timestamps_enabled": true, "longform_notetweets_rich_text_read_enabled": true, "longform_notetweets_inline_media_enabled": true, "responsive_web_enhance_cards_enabled": false}
	fieldToggles := map[string]interface{}{"withArticlePlainText": false}
	params := map[string]interface{}{
		"variables":    variables,
		"features":     features,
		"fieldToggles": fieldToggles,
	}
	res, err := c.gql("GET", "RaAkBb4XXis-atDL3rV-xw", "Likes", params)
	util.LogRaw([]string{"Likes", userID}, res, false)
	if err != nil {
		return nil, model.Cursor{}, err
	}
	return ParseTimelineEntriesTweets(res)
}

func (c *Client) Bookmarks(cursor ...string) ([]model.Tweet, model.Cursor, error) {
	if c.isGuest {
		return nil, model.Cursor{}, errors.New("Bookmarks API is not available with guest token")
	}
	zap.L().Debug("Bookmarks")
	variables := map[string]interface{}{"count": 20, "includePromotedContent": true}
	if len(cursor) > 0 {
		if cursor[0] != "" {
			variables["cursor"] = cursor[0]
		}
	}
	features := map[string]interface{}{"graphql_timeline_v2_bookmark_timeline": true, "rweb_tipjar_consumption_enabled": true, "responsive_web_graphql_exclude_directive_enabled": true, "verified_phone_label_enabled": false, "creator_subscriptions_tweet_preview_api_enabled": true, "responsive_web_graphql_timeline_navigation_enabled": true, "responsive_web_graphql_skip_user_profile_image_extensions_enabled": false, "communities_web_enable_tweet_community_results_fetch": true, "c9s_tweet_anatomy_moderator_badge_enabled": true, "articles_preview_enabled": true, "tweetypie_unmention_optimization_enabled": true, "responsive_web_edit_tweet_api_enabled": true, "graphql_is_translatable_rweb_tweet_is_translatable_enabled": true, "view_counts_everywhere_api_enabled": true, "longform_notetweets_consumption_enabled": true, "responsive_web_twitter_article_tweet_consumption_enabled": true, "tweet_awards_web_tipping_enabled": false, "creator_subscriptions_quote_tweet_preview_enabled": false, "freedom_of_speech_not_reach_fetch_enabled": true, "standardized_nudges_misinfo": true, "tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true, "tweet_with_visibility_results_prefer_gql_media_interstitial_enabled": true, "rweb_video_timestamps_enabled": true, "longform_notetweets_rich_text_read_enabled": true, "longform_notetweets_inline_media_enabled": true, "responsive_web_enhance_cards_enabled": false}
	params := map[string]interface{}{
		"variables": variables,
		"features":  features,
	}
	res, err := c.gql("GET", "yzqS_xq0glDD7YZJ2YDaiA", "Bookmarks", params)
	util.LogRaw([]string{"Bookmarks"}, res, false)
	if err != nil {
		return nil, model.Cursor{}, err
	}
	return ParseTimelineEntriesBookmarksTweets(res)
}

func (c *Client) SearchTimeline(rawQuery string, cursor ...string) ([]model.Tweet, model.Cursor, error) {
	if c.isGuest {
		return nil, model.Cursor{}, errors.New("SearchTimeline API is not available with guest token")
	}
	zap.L().Debug("SearchTimeline", zap.String("rawQuery", rawQuery))
	variables := map[string]interface{}{"rawQuery": rawQuery, "count": 20, "querySource": "typeahead_click", "product": "Latest"}
	if len(cursor) > 0 {
		if cursor[0] != "" {
			variables["cursor"] = cursor[0]
		}
	}
	features := map[string]interface{}{"rweb_tipjar_consumption_enabled": true, "responsive_web_graphql_exclude_directive_enabled": true, "verified_phone_label_enabled": false, "creator_subscriptions_tweet_preview_api_enabled": true, "responsive_web_graphql_timeline_navigation_enabled": true, "responsive_web_graphql_skip_user_profile_image_extensions_enabled": false, "communities_web_enable_tweet_community_results_fetch": true, "c9s_tweet_anatomy_moderator_badge_enabled": true, "articles_preview_enabled": true, "responsive_web_edit_tweet_api_enabled": true, "graphql_is_translatable_rweb_tweet_is_translatable_enabled": true, "view_counts_everywhere_api_enabled": true, "longform_notetweets_consumption_enabled": true, "responsive_web_twitter_article_tweet_consumption_enabled": true, "tweet_awards_web_tipping_enabled": false, "creator_subscriptions_quote_tweet_preview_enabled": false, "freedom_of_speech_not_reach_fetch_enabled": true, "standardized_nudges_misinfo": true, "tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true, "rweb_video_timestamps_enabled": true, "longform_notetweets_rich_text_read_enabled": true, "longform_notetweets_inline_media_enabled": true, "responsive_web_enhance_cards_enabled": false}
	params := map[string]interface{}{
		"variables": variables,
		"features":  features,
	}
	res, err := c.gql("GET", "MJpyQGqgklrVl_0X9gNy3A", "SearchTimeline", params)
	util.LogRaw([]string{"SearchTimeline", rawQuery}, res, false)
	if err != nil {
		return nil, model.Cursor{}, err
	}
	return ParseSearchTimelineEntriesTweets(res)
}
