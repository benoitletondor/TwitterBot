package main

import (
	"./db"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Content struct {
	text     string
	url      string
	hashtags string
}

func generateTweetContent() (Content, error) {
	contents, err := callAPI()
	if err != nil {
		return Content{}, err
	}

	for _, content := range contents {
		if strings.Contains(content.text, "\\") {
			continue
		}

		tweetExists, err := db.HasTweetWithContent(content.text)

		if err == nil && !tweetExists {
			return addHashTags(content), nil
		}
	}

	return Content{}, errors.New("No tweet content found")
}

type ByRandom []string

func (a ByRandom) Len() int           { return len(a) }
func (a ByRandom) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRandom) Less(i, j int) bool { return rand.Intn(2) > 0 }

func addHashTags(content Content) Content {
	numberOfTags := rand.Intn(4) // Max 3 hashtags
	tags := make([]string, len(HASHTAGS))
	copy(tags, HASHTAGS)

	sort.Sort(ByRandom(tags))

	margin := 140 - T_CO_URL_LENGTH - len(content.text) - 1 // -1 for the space between link and text

	for i, hashtag := range tags {
		if i >= numberOfTags {
			return content
		}

		if margin-len(hashtag)-2 < 0 { // -2 = the space before the new hashtag and the #
			return content
		}

		margin -= len(hashtag) + 2
		content.hashtags = content.hashtags + " #" + hashtag
	}

	return content
}

func hasReachDailyTweetLimit() (bool, error) {
	var from time.Time
	var to time.Time

	now := time.Now()

	if now.Hour() >= WAKE_UP_HOUR {
		from = time.Date(now.Year(), now.Month(), now.Day(), WAKE_UP_HOUR, 0, 0, 0, now.Location())
	} else {
		yesterday := now.Add(-time.Duration(24) * time.Hour)
		from = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), WAKE_UP_HOUR, 0, 0, 0, yesterday.Location())
	}

	if now.Hour() < GO_TO_BED_HOUR {
		to = time.Date(now.Year(), now.Month(), now.Day(), GO_TO_BED_HOUR, 0, 0, 0, now.Location())
	} else {
		tomorrow := now.Add(time.Duration(24) * time.Hour)
		to = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), GO_TO_BED_HOUR, 0, 0, 0, tomorrow.Location())
	}

	count, err := db.GetNumberOfTweetsBetweenDates(from, to)
	if err != nil {
		return true, err
	}

	return count >= MAX_TWEET_IN_A_DAY, nil
}

func callAPI() ([]Content, error) {
	content, err := getWebserviceContent(DATA_SOURCES[rand.Intn(len(DATA_SOURCES))])
	if err != nil {
		fmt.Println("Error while calling API")
		return nil, err
	} else {
		// Fill the record with the data from the JSON
		var jsonObj map[string]interface{}
		err = json.Unmarshal(content, &jsonObj)
		if err != nil {
			fmt.Println("An error occurred while converting our JSON to an object")
			return nil, err
		}

		contents := make([]Content, 0)

		if val, ok := jsonObj["results"]; ok {

			if results, ok := val.(map[string]interface{}); ok {

				if val, ok := results["collection1"]; ok {

					if collection, ok := val.([]interface{}); ok {

						for _, item := range collection {

							if property, ok := item.(map[string]interface{}); ok {

								if property1, ok := property["property1"]; ok {

									if contentData, ok := property1.(map[string]interface{}); ok {

										if text, ok := contentData["text"].(string); ok {

											if href, ok := contentData["href"].(string); ok {

												contents = append(contents, Content{text: text, url: href})

											} else {
												fmt.Println("Error mapping href as string")
												return nil, errors.New("json mapping error")
											}

										} else {
											fmt.Println("Error mapping text as string")
											return nil, errors.New("json mapping error")
										}

									} else {
										fmt.Println("Error mapping property1 as json obj")
										return nil, errors.New("json mapping error")
									}

								} else {
									fmt.Println("Missing mappings property1")
									return nil, errors.New("json mapping error")
								}

							} else {
								fmt.Println("Error mappings property1")
								return nil, errors.New("json mapping error")
							}

						}

					} else {
						fmt.Println("Error mappings collection1 as array")
						return nil, errors.New("json mapping error")
					}

				} else {
					fmt.Println("Error mappings collection1")
					return nil, errors.New("json mapping error")
				}

			} else {
				fmt.Println("Error mappings results as json obj")
				return nil, errors.New("json mapping error")
			}

			return contents, nil

		} else {
			fmt.Println("No field results in json")
			return nil, errors.New("json mapping error")
		}
	}
}

func getWebserviceContent(url string) ([]byte, error) {
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// Defer the closing of the body
	defer resp.Body.Close()
	// Read the content into a byte array
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// At this point we're done - simply return the bytes
	return body, nil
}
