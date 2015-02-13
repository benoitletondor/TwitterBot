package main

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/robfig/cron"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

var api *anaconda.TwitterApi
var db *sql.DB

func main() {
	// Init Twitter API
	anaconda.SetConsumerKey(CONSUMER_KEY)
	anaconda.SetConsumerSecret(CONSUMER_SECRET)
	api = anaconda.NewTwitterApi(TOKEN, TOKEN_SECRET)

	// Init Mysql DB
	dbLink, err := sql.Open("mysql", MYSQL_USER+":"+MYSQL_PASSWORD+"@/"+MYSQL_SCHEMA)
	if err != nil {
	    panic(err.Error())
	}

	// Open doesn't open a connection. Validate DSN data:
	err = dbLink.Ping()
	if err != nil {
	    panic(err.Error()) 
	}

	// Set up global var
	db = dbLink

	c := cron.New()
	c.AddFunc("@every 5s", bot)
	c.Start()

	select {} // block forever
}

func bot() {
	fmt.Println("Hello world")

	performAction()
}