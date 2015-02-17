package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/kapsteur/franco"
	"math/rand"
	"net/url"
	"strings"
)

func shouldReply(tweet anaconda.Tweet) bool {
	friendships, err := api.GetFriendshipsLookup(url.Values{"screen_name": []string{tweet.User.ScreenName}})
	if err != nil {
		fmt.Println("Error while querying twitter api for friendships", err)
		return false
	}

	following := false
	for _, friendship := range friendships {
		if stringInSlice("followed_by", friendship.Connections) {
			following = true
		}
	}

	if !following {
		fmt.Println("Avoid reply to " + tweet.User.ScreenName + " cause he's not following us :C")
		return false
	}

	words := strings.Split(tweet.Text, " ")

	keywords := []string{"thank", "thanks", "follow", "following", "following!"}

	occurrences := 0
	max_occurrences := len(keywords)

	for _, word := range words {
		if stringInSlice(strings.ToLower(word), keywords) {
			occurrences++
		}
	}

	confidence := (float64(occurrences) / float64(max_occurrences))

	return confidence > 0.2
}

func buildReply(tweet anaconda.Tweet) string {
	greetings := []string{"hello!", "hey", "hi", "well,", ""}

	thanks := []string{"thanks", "thank you", "many thanks", "thx"}

	follows := []string{"following me", "the follow"}

	reciprocals := []string{"too", "as well", ""}

	greet := greetings[rand.Intn(len(greetings))]
	thank := thanks[rand.Intn(len(thanks))]
	follow := follows[rand.Intn(len(follows))]
	reciprocal := reciprocals[rand.Intn(len(reciprocals))]

	return "@" + tweet.User.ScreenName + " " + greet + " " + thank + " for " + follow + " " + reciprocal
}

func isRightLanguage(tweet anaconda.Tweet) bool {
	lang := franco.DetectOne(removeUserNames(tweet.Text))
	userLang := franco.DetectOne(tweet.User.Description)

	if !stringInSlice(lang.Code, ACCEPTED_LANGUAGES) {
		fmt.Println("Ignoring tweet in " + lang.Code + ", not english : " + tweet.Text)
		return false
	}

	if !stringInSlice(lang.Code, ACCEPTED_LANGUAGES) {
		fmt.Println("Ignoring user desc in " + userLang.Code + ", not english : " + tweet.User.Description)
		return false
	}

	return true
}

func removeUserNames(tweetText string) string {
	value := ""

	words := strings.Split(tweetText, " ")
	for _, word := range words {
		if !strings.HasPrefix(word, "@") {
			value += word + " "
		}
	}

	return value
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
