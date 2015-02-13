package main

import "math/rand"
import "fmt"

const (
	_FOLLOW   = iota
	_RETWEET  = iota
	_FAVORITE = iota
)

type Action struct {
	name int
	weight int
}

func performAction() {
	actions := make([]Action, 0, 3)

	actions = append(actions, Action{name:_FOLLOW, weight: ACTION_FOLLOW_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name:_RETWEET, weight: ACTION_RETWEET_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name:_FAVORITE, weight: ACTION_FAVORITE_WEIGHT * rand.Intn(100)})

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