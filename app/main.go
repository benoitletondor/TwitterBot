package main

import "github.com/ChimeraCoder/anaconda"
import "fmt"

func main() {
	anaconda.SetConsumerKey(CONSUMER_KEY)
	anaconda.SetConsumerSecret(CONSUMER_SECRET)
	api := anaconda.NewTwitterApi(TOKEN, TOKEN_SECRET)

	searchResult, _ := api.GetSearch("golang", nil)
	for _ , tweet := range searchResult.Statuses {
    	fmt.Println(tweet.Text)
	}
}