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

package content

import (
	"./../db"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type ContentAPI interface {
	callAPI() ([]Content, error)
}

var apis []ContentAPI
var hashtags []string
var urlLength int

type Content struct {
	Text     string
	Url      string
	Hashtags string
}

func Init(tags []string, urlL int) {
	apis = make([]ContentAPI, 0)
	hashtags = tags
	urlLength = urlL
}

func RegisterAPI(contentAPI ContentAPI) {
	apis = append(apis, contentAPI)
}

func GenerateTweetContent() (Content, error) {
	contents, err := apis[rand.Intn(len(apis))].callAPI()
	if err != nil {
		return Content{}, err
	}

	for _, content := range contents {
		if strings.Contains(strconv.QuoteToASCII(content.Text), "\\") {
			continue
		}

		tweetExists, err := db.HasTweetWithContent(content.Text)

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
	tags := make([]string, len(hashtags))
	copy(tags, hashtags)

	sort.Sort(ByRandom(tags))

	margin := 140 - urlLength - len(content.Text) - 1 // -1 for the space between link and text

	for i, hashtag := range tags {
		if i >= numberOfTags {
			return content
		}

		if margin-len(hashtag)-2 < 0 { // -2 = the space before the new hashtag and the #
			return content
		}

		margin -= len(hashtag) + 2
		content.Hashtags = content.Hashtags + " #" + hashtag
	}

	return content
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
