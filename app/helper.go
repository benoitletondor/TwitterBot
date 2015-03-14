package main

import (
	"./anaconda"
	"fmt"
	"net/url"
	"strings"
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
	v.Add("count", "50")

	return url.QueryEscape(searchString), v
}

func isMentionOrRT(tweet anaconda.Tweet) bool {
	return strings.HasPrefix(tweet.Text, "RT") || strings.HasPrefix(tweet.Text, "@")
}

func isMe(tweet anaconda.Tweet) bool {
	return strings.ToLower(tweet.User.ScreenName) == strings.ToLower(USER_NAME)
}
