package tools

import "github.com/nharu-0630/twitter-api-client/model"

func ContainsJapanese(str string) bool {
	for _, r := range str {
		if (r >= '\u3040' && r <= '\u309F') || (r >= '\u30A0' && r <= '\u30FF') || (r >= '\u4E00' && r <= '\u9FFF') {
			return true
		}
	}
	return false
}

func IsJapaneseUser(user model.User) bool {
	if ContainsJapanese(user.Legacy.Description) || ContainsJapanese(user.Legacy.Name) || ContainsJapanese(user.Legacy.Location) {
		return true
	}
	return false
}

func IsJapaneseTweet(tweet model.Tweet) bool {
	return ContainsJapanese(tweet.Legacy.FullText)
}
