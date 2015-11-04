/*
 *   Copyright 2015 Benoit LETONDOR
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/jsgoecke/go-wit"
	"math/rand"
	"strings"
)

const (
	INTENT_HI           = "hi"
	INTENT_NICE_ARTICLE = "nice_article"
	INTENT_THANK_FOLLOW = "thank_follow"
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

	if intent == INTENT_HI {
		return buildHiIntentResponse(tweet), nil
	} else if intent == INTENT_NICE_ARTICLE {
		return buildNiceArticleIntentResponse(tweet), nil
	} else if intent == INTENT_THANK_FOLLOW {
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
		if strings.HasPrefix(word, "@") || strings.HasPrefix(word, "http") {
			continue
		} else if strings.HasPrefix(word, "#") {
			cleaned += strings.TrimPrefix(word, "#") + " "
		}

		cleaned += word + " "
	}

	return cleaned
}
