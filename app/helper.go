package main

import (
	"./db"
	"fmt"
	"github.com/benoitletondor/anaconda"
	"net/url"
	"strings"
	"time"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func isUserFollowing(userName string) (bool, error) {
	friendships, err := api.GetFriendshipsLookup(url.Values{"screen_name": []string{userName}})
	if err != nil {
		fmt.Println("Error while querying twitter api for friendships", err)
		return false, err
	}

	following := false
	for _, friendship := range friendships {
		if stringInSlice("followed_by", friendship.Connections) {
			following = true
		}
	}

	return following, nil
}

func isUserAcceptable(tweet anaconda.Tweet) bool {
	words := strings.Split(tweet.Text, " ")
	for _, word := range words {
		if stringInSlice(strings.ToLower(word), BANNED_KEYWORDS) {
			return false
		}
	}

	if tweet.User.Description == "" {
		return false
	}

	words = strings.Split(tweet.User.Description, " ")
	for _, word := range words {
		if stringInSlice(strings.ToLower(word), BANNED_KEYWORDS) {
			return false
		}
	}

	return true
}

func generateAPISearchValues(word string) (string, url.Values) {
	searchString := word

	for _, word := range BANNED_KEYWORDS {
		searchString += " -" + word
	}

	v := url.Values{}
	v.Add("lang", ACCEPTED_LANGUAGE)
	v.Add("count", "100")
	v.Add("result_type", "recent")

	return url.QueryEscape(searchString), v
}

func isMentionOrRT(tweet anaconda.Tweet) bool {
	return strings.HasPrefix(tweet.Text, "RT") || strings.HasPrefix(tweet.Text, "@")
}

func isMe(tweet anaconda.Tweet) bool {
	return strings.ToLower(tweet.User.ScreenName) == strings.ToLower(USER_NAME)
}

func hasReachDailyTweetLimit() (bool, error) {
	var from time.Time
	var to time.Time

	now := time.Now()

	if now.Hour() >= WAKE_UP_HOUR {
		from = time.Date(now.Year(), now.Month(), now.Day(), WAKE_UP_HOUR, 0, 0, 0, now.Location())
	} else {
		yesterday := now.Add(-time.Duration(24) * time.Hour)
		from = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), WAKE_UP_HOUR, 0, 0, 0, yesterday.Location())
	}

	if now.Hour() < GO_TO_BED_HOUR {
		to = time.Date(now.Year(), now.Month(), now.Day(), GO_TO_BED_HOUR, 0, 0, 0, now.Location())
	} else {
		tomorrow := now.Add(time.Duration(24) * time.Hour)
		to = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), GO_TO_BED_HOUR, 0, 0, 0, tomorrow.Location())
	}

	count, err := db.GetNumberOfTweetsBetweenDates(from, to)
	if err != nil {
		return true, err
	}

	return count >= MAX_TWEET_IN_A_DAY, nil
}
