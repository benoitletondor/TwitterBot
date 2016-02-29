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
	"./content"
	"./db"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	_FOLLOW     = iota
	_UNFOLLOW   = iota
	_FAVORITE   = iota
	_UNFAVORITE = iota
	_TWEET      = iota
	_REPLY      = iota
)

type Action struct {
	name   int
	weight int
}

func performDailyAction() {
	actions := make([]Action, 0, 6)

	actions = append(actions, Action{name: _FOLLOW, weight: ACTION_FOLLOW_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name: _UNFOLLOW, weight: ACTION_UNFOLLOW_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name: _FAVORITE, weight: ACTION_FAVORITE_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name: _UNFAVORITE, weight: ACTION_UNFAVORITE_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name: _TWEET, weight: ACTION_TWEET_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name: _REPLY, weight: ACTION_REPLY_WEIGHT * rand.Intn(100)})

	selectAndPerformAction(actions)
}

func performNightlyAction() {
	actions := make([]Action, 0, 6)

	actions = append(actions, Action{name: _FOLLOW, weight: ACTION_NIGHTLY_FOLLOW_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name: _UNFOLLOW, weight: ACTION_NIGHTLY_UNFOLLOW_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name: _FAVORITE, weight: ACTION_NIGHTLY_FAVORITE_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name: _UNFAVORITE, weight: ACTION_NIGHTLY_UNFAVORITE_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name: _TWEET, weight: ACTION_NIGHTLY_TWEET_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name: _REPLY, weight: ACTION_NIGHTLY_REPLY_WEIGHT * rand.Intn(100)})

	selectAndPerformAction(actions)
}

func selectAndPerformAction(actions []Action) {
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
	case _UNFOLLOW:
		actionUnfollow()
		break
	case _FAVORITE:
		actionFavorite()
		break
	case _UNFAVORITE:
		actionUnfavorite()
		break
	case _TWEET:
		actionTweet()
		break
	case _REPLY:
		actionReply()
		break
	}
}

func actionFollow() {
	log.Println("Action follow")

	search, v := generateAPISearchValues(KEYWORDS[rand.Intn(len(KEYWORDS))])

	searchResult, err := api.GetSearch(search, v)
	if err != nil {
		log.Println("Error while querying twitter API", err)
		return
	}

	for _, tweet := range searchResult.Statuses {
		if !isUserAcceptable(tweet) {
			log.Println("Ignoring user for follow : @" + tweet.User.ScreenName)
			continue
		}

		if isMentionOrRT(tweet) {
			log.Println("Ignoring tweet for follow, mention or RT")
			continue
		}

		if isMe(tweet) {
			log.Println("Ignoring my own tweet for follow")
			continue
		}

		follow, err := db.AlreadyFollow(tweet.User.Id)
		if err != nil {
			log.Println("Error while checking if already follow")
			return
		}

		if follow {
			log.Println("Ignoring user for follow, already follow @" + tweet.User.ScreenName)
			continue
		}

		alreadyFollowMe, err := isUserFollowing(tweet.User.ScreenName)
		if err != nil {
			log.Println("Error while checking user already follow", err)
			return
		}

		if alreadyFollowMe {
			log.Println("Ignoring user @" + tweet.User.ScreenName + " for follow cause he already follow us")
			continue
		}

		err = db.Follow{UserId: tweet.User.Id, UserName: tweet.User.ScreenName, Status: tweet.Text, FollowDate: time.Now()}.Persist()
		if err != nil {
			log.Println("Error while persisting follow", err)
			return
		}

		_, err = api.FollowUser(tweet.User.ScreenName)
		if err != nil {
			log.Println("Error while following user "+tweet.User.ScreenName+" : ", err)
		}

		log.Println("Now follow ", tweet.User.ScreenName)
		return
	}
}

func actionUnfollow() {
	log.Println("Action unfollow")

	date := time.Now()
	duration, err := time.ParseDuration("-72h") // -3 days
	date = date.Add(duration)

	follows, err := db.GetNotUnfollowed(date, UNFOLLOW_LIMIT_IN_A_ROW)
	if err != nil {
		log.Println("Error while querying db to find people to unfollow", err)
		return
	}

	for _, follow := range follows {
		follow.LastAction = time.Now()

		isFollowing, err := isUserFollowing(follow.UserName)
		if err != nil {
			log.Println("Error while querying API for friendships", err)
			return
		}

		if isFollowing {
			err = follow.Persist()
			if err != nil {
				log.Println("Error while persisting follow", err)
				return
			}

			continue
		}

		follow.UnfollowDate = time.Now()

		err = follow.Persist()
		if err != nil {
			log.Println("Error while persisting follow", err)
			return
		}

		_, err = api.UnfollowUser(follow.UserName)
		if err != nil {
			log.Println("Error while querying API to unfollow @"+follow.UserName, err)
			continue
		}

		log.Println("Unfollowed @" + follow.UserName)
	}
}

func actionFavorite() {
	log.Println("Action fav")

	search, v := generateAPISearchValues(KEYWORDS[rand.Intn(len(KEYWORDS))])

	searchResult, err := api.GetSearch(search, v)
	if err != nil {
		log.Println("Error while querying twitter API", err)
		return
	}

	i := 0
	for _, tweet := range searchResult.Statuses {
		if i >= FAV_LIMIT_IN_A_ROW {
			return
		}

		if !isUserAcceptable(tweet) {
			log.Println("Ignoring user for favorite : @" + tweet.User.ScreenName)
			continue
		}

		if isMentionOrRT(tweet) {
			log.Println("Ignoring tweet for favorite, mention or RT")
			continue
		}

		if isMe(tweet) {
			log.Println("Ignoring my own tweet for favorite")
			continue
		}

		follow, err := db.AlreadyFollow(tweet.User.Id)
		if err == nil && follow {
			log.Println("Ignoring tweet for favorite, already follow @" + tweet.User.ScreenName)
			continue
		}

		alreadyFav, err := db.HasAlreadyFav(tweet.Id)
		if err != nil {
			log.Println("Error while checking already fav", err)
			return
		}

		if alreadyFav {
			log.Println("Ignoring tweet for favorite, already fav tweet from @" + tweet.User.ScreenName)
			continue
		}

		err = db.Favorite{UserId: tweet.User.Id, UserName: tweet.User.ScreenName, TweetId: tweet.Id, Status: tweet.Text, FavDate: time.Now()}.Persist()
		if err != nil {
			log.Println("Error while persisting fav", err)
			return
		}

		_, err = api.Favorite(tweet.Id)
		if err != nil {
			if strings.Contains(err.Error(), "139") { // Case of an already favorited tweet
				continue
			}

			log.Println("Error while favoriting tweet", err)
		} else {
			log.Println("Just favorited tweet : ", tweet.Text)
		}

		i++
	}
}

func actionUnfavorite() {
	log.Println("Action unfavorite")

	date := time.Now()
	duration, err := time.ParseDuration("-72h") // -3 days
	date = date.Add(duration)

	favs, err := db.GetNotUnfavorite(date, UNFAVORITE_LIMIT_IN_A_ROW)
	if err != nil {
		log.Println("Error while querying db to find tweet to unfav", err)
		return
	}

	for _, fav := range favs {
		fav.LastAction = time.Now()

		isFollowing, err := isUserFollowing(fav.UserName)
		if err != nil {
			log.Println("Error while querying API for friendships", err)
			return
		}

		if isFollowing {
			err = fav.Persist()
			if err != nil {
				log.Println("Error while persisting fav", err)
				return
			}

			continue
		}

		fav.UnfavDate = time.Now()

		err = fav.Persist()
		if err != nil {
			log.Println("Error while persisting fav", err)
			return
		}

		_, err = api.Unfavorite(fav.TweetId)
		if err != nil {
			log.Println("Error while querying API to unfav : "+fav.Status, err)
			continue
		}

		log.Println("Unfaved @" + fav.Status)
	}
}

func actionTweet() {
	log.Println("Action tweet")

	hasReachLimit, err := hasReachDailyTweetLimit()
	if err != nil {
		log.Println("Error while getting daily limit reached", err)
		return
	}

	if hasReachLimit {
		log.Println("Day tweet limit reached, abording")
		return
	}

	content, err := content.GenerateTweetContent()
	if err != nil {
		log.Println("Error while getting tweet content : ", err)
		return
	}

	tweetText := content.Text + " " + content.Url + content.Hashtags

	err = db.Tweet{Content: tweetText, Date: time.Now()}.Persist()
	if err != nil {
		log.Println("Error while persisting tweet", err)
		return
	}

	tweet, err := api.PostTweet(tweetText, nil)
	if err != nil {
		log.Println("Error while posting tweet", err)
		return
	}

	log.Println("Tweet posted : ", tweet.Text)
}

func actionReply() {
	log.Println("Action reply")

	tweets, err := api.GetMentionsTimeline(nil)
	if err != nil {
		log.Println("Error while querying twitter mention API", err)
		return
	}

	for _, tweet := range tweets {

		replied, err := db.HasAlreadyReplied(tweet.Id)
		if err == nil && !replied {
			log.Println("Building reply for tweet : " + tweet.Text)

			response, err := buildReply(tweet)
			if err != nil {
				log.Println("Error while building reply", err)
				return
			}

			err = db.Reply{UserId: tweet.User.Id, UserName: tweet.User.ScreenName, TweetId: tweet.Id, Status: tweet.Text, Answer: response, ReplyDate: time.Now()}.Persist()
			if err != nil {
				log.Println("Error while persisting reply", err)
				return
			}

			if response != "" {
				v := url.Values{}
				v.Add("in_reply_to_status_id", strconv.FormatInt(tweet.Id, 10))

				respTweet, err := api.PostTweet(response, v)
				if err != nil {
					log.Println("Error while posting reply", err)
					return
				}

				log.Println("Reply posted : ", respTweet.Text)
			} else {
				log.Println("No response found for tweet : " + tweet.Text)
			}

			return
		}
	}

	log.Println("Nothing to reply found :(")
}
