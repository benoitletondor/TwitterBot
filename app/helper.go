package main

import (
	"fmt"
	"net/url"
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
