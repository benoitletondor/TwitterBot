package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/jsgoecke/go-wit"
	"math/rand"
	"strings"
)

func buildReply(tweet anaconda.Tweet) (string, error) {
	message := cleanTweetMessage(tweet.Text)
	if message == "" {
		return "", nil
	}

	// Process a text message
	request := &wit.MessageRequest{}
	request.Query = message

	result, err := witclient.Message(request)
	if err != nil {
		return "", err
	}

	intent := result.Outcome.Intent
	if result.Outcome.Confidence < 0.5 {
		fmt.Println("Not enough confidence for intent : " + intent)
		return "", nil
	}

	if intent == "hi" {
		return buildHiIntentResponse(tweet), nil
	} else if intent == "nice_article" {
		return buildNiceArticleIntentResponse(tweet), nil
	} else if intent == "thank_follow" {
		return buildThanksFollowIntentResponse(tweet), nil
	}

	return "", nil
}

func buildHiIntentResponse(tweet anaconda.Tweet) string {
	greetings := []string{"hello!", "hey", "yo"}

	return buildMention(tweet.User, greetings[rand.Intn(len(greetings))])
}

func buildNiceArticleIntentResponse(tweet anaconda.Tweet) string {
	greetings := []string{"hello!", "hey", "hi", "well,", ""}
	thanks := []string{"thanks", "thank you", "many thanks", "thx"}
	messages := []string{"reading", "your tweet", "your message"}

	greet := greetings[rand.Intn(len(greetings))]
	thank := thanks[rand.Intn(len(thanks))]
	message := messages[rand.Intn(len(messages))]

	return buildMention(tweet.User, greet+" "+thank+" for "+message)
}

func buildThanksFollowIntentResponse(tweet anaconda.Tweet) string {
	greetings := []string{"hello!", "hey", "hi", "well,", ""}
	thanks := []string{"thanks", "thank you", "many thanks", "thx"}
	follows := []string{"following me", "the follow"}
	messages := []string{"your message", "your tweet", "your mention"}
	reciprocals := []string{"too", "as well", ""}

	following, err := isUserFollowing(tweet.User.ScreenName)
	if following && err == nil {
		greet := greetings[rand.Intn(len(greetings))]
		thank := thanks[rand.Intn(len(thanks))]
		follow := follows[rand.Intn(len(follows))]
		reciprocal := reciprocals[rand.Intn(len(reciprocals))]

		return buildMention(tweet.User, greet+" "+thank+" for "+follow+" "+reciprocal)
	} else {
		greet := greetings[rand.Intn(len(greetings))]
		thank := thanks[rand.Intn(len(thanks))]
		message := messages[rand.Intn(len(messages))]

		return buildMention(tweet.User, greet+" "+thank+" for "+message)
	}
}

func buildMention(user anaconda.User, text string) string {
	return "@" + user.ScreenName + " " + text
}

func cleanTweetMessage(message string) string {
	cleaned := ""

	words := strings.Split(message, " ")
	for _, word := range words {
		if strings.HasPrefix(word, "@") {
			continue
		} else if strings.HasPrefix(word, "#") {
			cleaned += strings.TrimPrefix(word, "#") + " "
		} else if strings.HasPrefix(word, "http") {
			continue
		}

		cleaned += word + " "
	}

	return cleaned
}
