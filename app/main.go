package main

import (
	"./db"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/robfig/cron"
)

var api *anaconda.TwitterApi

func main() {
	// Init Twitter API
	anaconda.SetConsumerKey(CONSUMER_KEY)
	anaconda.SetConsumerSecret(CONSUMER_SECRET)
	api = anaconda.NewTwitterApi(TOKEN, TOKEN_SECRET)

	database, err := db.Init(MYSQL_USER, MYSQL_PASSWORD, MYSQL_SCHEMA)
	if err != nil {
		panic(err.Error())
	}

	defer database.Close()

	c := cron.New()
	c.AddFunc(ACTIONS_INTERVAL, bot)
	c.Start()

	fmt.Println("Hello world")

	bot()

	select {} // block forever
}

func bot() {
	fmt.Println("----------- Waking up!")

	performAction()

	fmt.Println("----------- Goes to sleep")
}
