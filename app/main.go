package main

import "github.com/ChimeraCoder/anaconda"
import "github.com/robfig/cron"
import "fmt"

var api *anaconda.TwitterApi

func main() {
	anaconda.SetConsumerKey(CONSUMER_KEY)
	anaconda.SetConsumerSecret(CONSUMER_SECRET)
	api = anaconda.NewTwitterApi(TOKEN, TOKEN_SECRET)

	c := cron.New()
	c.AddFunc("@every 5s", bot)
	c.Start()

	select {} // block forever
}

func bot() {
	fmt.Println("Hello world")

	performAction(api)
}