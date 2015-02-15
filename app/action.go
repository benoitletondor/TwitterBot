package main

import (
	"./db"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	_FOLLOW   = iota
	_REPLY    = iota
	_FAVORITE = iota
	_TWEET    = iota
)

type Action struct {
	name   int
	weight int
}

func performAction() {
	actions := make([]Action, 0, 4)

	actions = append(actions, Action{name: _FOLLOW, weight: ACTION_FOLLOW_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name: _REPLY, weight: ACTION_REPLY_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name: _FAVORITE, weight: ACTION_FAVORITE_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name: _TWEET, weight: ACTION_TWEET_WEIGHT * rand.Intn(100)})

	selectedAction := Action{name: -1, weight: -1}

	for _, action := range actions {
		if action.weight > selectedAction.weight {
			selectedAction = action
		}
	}

	switch selectedAction.name {
	case _FOLLOW:
		actionFollow()
		break
	case _REPLY:
		actionReply()
		break
	case _FAVORITE:
		actionFavorite()
		break
	case _TWEET:
		actionTweet()
		break
	}
}

func actionFollow() {
	fmt.Println("Action follow")

	searchResult, err := api.GetSearch(KEYWORDS[rand.Intn(len(KEYWORDS))], nil)
	if err != nil {
		fmt.Println("Error while querying twitter API", err)
		return
	}

	for _, tweet := range searchResult.Statuses {

		follow, err := db.AlreadyFollow(tweet.User.Id)
		if err == nil && !follow {

			err := db.Follow{UserId: tweet.User.Id, UserName: tweet.User.ScreenName, Status: tweet.Text, FollowDate: time.Now()}.Persist()
			if err != nil {
				fmt.Println("Error while persisting follow", err)
				return
			}

			_, err = api.FollowUser(tweet.User.ScreenName)
			if err != nil {
				fmt.Println("Error while following user "+tweet.User.ScreenName+" : ", err)
			}

			fmt.Println("Now follow ", tweet.User.ScreenName)
			return
		}

	}
}

func actionReply() {
	fmt.Println("Action reply")

	tweets, err := api.GetMentionsTimeline(nil)
	if err != nil {
		fmt.Println("Error while querying twitter mention API", err)
		return
	}

	for _, tweet := range tweets {

		replied, err := db.HasAlreadyReplied(tweet.Id)
		if err == nil && !replied && shouldReply(tweet) {

			response := buildReply(tweet)

			err = db.Reply{UserId: tweet.User.Id, UserName: tweet.User.ScreenName, TweetId: tweet.Id, Status: tweet.Text, Answer: response, ReplyDate: time.Now()}.Persist()
			if err != nil {
				fmt.Println("Error while persisting reply", err)
				return
			}

			respTweet, err := api.PostTweet(response, nil)
			if err != nil {
				fmt.Println("Error while posting reply", err)
				return
			}

			fmt.Println("Reply posted : ", respTweet.Text)
			return
		}
	}

	fmt.Println("Nothing to reply found :(")
}

func actionFavorite() {
	fmt.Println("Action fav")

	searchResult, err := api.GetSearch(KEYWORDS[rand.Intn(len(KEYWORDS))], nil)
	if err != nil {
		fmt.Println("Error while querying twitter API", err)
		return
	}

	i := 0
	for _, tweet := range searchResult.Statuses {
		if i >= FAV_LIMIT_IN_A_ROW {
			return
		}

		_, err = api.Favorite(tweet.Id)
		if err != nil {
			if strings.Contains(err.Error(), "139") { // Case of an already favorited tweet
				continue
			}

			fmt.Println("Error while favoriting tweet", err)
		} else {
			fmt.Println("Just favorited tweet : ", tweet.Text)
		}

		i++
	}
}

func actionTweet() {
	fmt.Println("Action tweet")

	content, err := generateTweetContent()
	if err != nil {
		fmt.Println("Error while getting tweet content : ", err)
		return
	}

	tweetText := content.text + " " + content.url

	err = db.Tweet{Content: tweetText, Date: time.Now()}.Persist()
	if err != nil {
		fmt.Println("Error while persisting tweet", err)
		return
	}

	tweet, err := api.PostTweet(tweetText, nil)
	if err != nil {
		fmt.Println("Error while posting tweet", err)
		return
	}

	fmt.Println("Tweet posted : ", tweet.Text)
}
