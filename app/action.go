package main

import (
	"math/rand"
	"fmt"
	"net/url"
)

const (
	_FOLLOW   = iota
	_RETWEET  = iota
	_FAVORITE = iota
	_TWEET    = iota
)

type Action struct {
	name int
	weight int
}

func performAction() {
	actions := make([]Action, 0, 4)

	actions = append(actions, Action{name:_FOLLOW, weight: ACTION_FOLLOW_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name:_RETWEET, weight: ACTION_RETWEET_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name:_FAVORITE, weight: ACTION_FAVORITE_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name:_TWEET, weight: ACTION_TWEET_WEIGHT * rand.Intn(100)})

	selectedAction := Action{name:-1, weight:-1}

	for _,action := range actions {
        if( action.weight > selectedAction.weight ){
        	selectedAction = action
        }
    }

	switch selectedAction.name {
		case _FOLLOW:
			actionFollow()
			break
		case _RETWEET:
			actionRetweet()
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
}

func actionRetweet() {
	fmt.Println("Action retweet")
}

func actionFavorite() {
	fmt.Println("Action fav")
}

func actionTweet() {
	fmt.Println("Action tweet")

	content, err := generateTweetContent()
	if( err != nil ) {
		fmt.Println("Error while getting tweet content : ", err)
		return
	}
	
	tweet, err := api.PostTweet(content.text+" "+content.url, url.Values{})
	if( err != nil ) {
		fmt.Println("Error while posting tweet", err)
		return
	}

	// TODO save in DB

	fmt.Println("Tweet posted : ", tweet)
}