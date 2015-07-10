package gohn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

const rootUrl = "https://hacker-news.firebaseio.com/v0"
const maxStories = 10
const timeBetweenReqMs = 10

const maxRetries = 100
const waitTimeMeanMs = 10

func jsonBytes(url string) ([]byte, error) {
	var err error
	var resp *http.Response
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			log.Printf("Retrying again %d\n", i)
		}
		resp, err = http.Get(url)
		if err == nil {
			defer resp.Body.Close()
			break
		}
		waitTime := time.Millisecond * time.Duration(i+1) * time.Duration(math.Max(waitTimeMeanMs*rand.ExpFloat64(), 100))
		time.Sleep(waitTime)
	}
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func TopStories() ([]*Item, error) {
	body, err := jsonBytes(rootUrl + "/topstories.json")
	list := make([]*Item, 0, 500)
	if err != nil {
		return list, err
	}

	itemids := make([]int, 0, 500)

	if err := json.Unmarshal(body, &itemids); err != nil {
		return list, err
	}
	num_stories := 0
	for _, id := range itemids {
		it, err := GetItem(id)
		if err != nil {
			log.Printf("Failed to get story %d, %s", id, err)
		}
		list = append(list, it)
		time.Sleep(timeBetweenReqMs * time.Millisecond)
		num_stories += 1
		if num_stories == maxStories {
			break
		}
	}
	return list, nil
}

// TODO retry this
func GetItem(id int) (*Item, error) {
	i := Item{}
	body, err := jsonBytes(fmt.Sprintf("%s/item/%d.json", rootUrl, id))
	if err != nil {
		return &i, err
	}
	err = json.Unmarshal(body, &i)
	if err != nil {
		return &i, err
	}
	return &i, nil
}
