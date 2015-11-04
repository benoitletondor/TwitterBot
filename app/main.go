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
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/jsgoecke/go-wit"
	"github.com/robfig/cron"
	"math/rand"
	"time"
)

var api *anaconda.TwitterApi
var witclient *wit.Client

func main() {
	// Init Twitter API
	anaconda.SetConsumerKey(CONSUMER_KEY)
	anaconda.SetConsumerSecret(CONSUMER_SECRET)
	api = anaconda.NewTwitterApi(TOKEN, TOKEN_SECRET)

	// Init DB
	database, err := db.Init(MYSQL_USER, MYSQL_PASSWORD, MYSQL_SCHEMA)
	if err != nil {
		panic(err.Error())
	}

	defer database.Close()

	// Init Content
	content.Init(HASHTAGS, T_CO_URL_LENGTH)

	for _, kimonoDataSourcesUrl := range KIMONO_DATA_SOURCES {
		content.RegisterAPI(content.KimonoContent{Url: kimonoDataSourcesUrl})
	}

	// Init WIT api
	witclient = wit.NewClient(WIT_ACCESS_TOKEN)

	// Init cron
	c := cron.New()
	c.AddFunc(ACTIONS_INTERVAL, bot)
	c.Start()

	// Init random
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Hello world")

	bot()

	select {} // block forever
}

func bot() {
	fmt.Println("----------- Waking up!")

	hour := time.Now().Hour()

	if GO_TO_BED_HOUR < WAKE_UP_HOUR && (hour >= WAKE_UP_HOUR || hour < GO_TO_BED_HOUR) {
		performDailyAction()
	} else if GO_TO_BED_HOUR > WAKE_UP_HOUR && (hour >= WAKE_UP_HOUR && hour < GO_TO_BED_HOUR) {
		performDailyAction()
	} else {
		performNightlyAction()
	}

	fmt.Println("----------- Goes to sleep")
}
