package gohn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const rootUrl = "https://hacker-news.firebaseio.com/v0"
const maxStories = 10
const timeBetweenReqMs = 1

func jsonBytes(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func TopStories() ([]Item, error) {
	body, err := jsonBytes(rootUrl + "/topstories.json")
	list := make([]Item, 0, 500)
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
func GetItem(id int) (Item, error) {
	i := Item{}
	body, err := jsonBytes(fmt.Sprintf("%s/item/%d.json", rootUrl, id))
	if err != nil {
		return i, err
	}
	err = json.Unmarshal(body, &i)
	if err != nil {
		return i, err
	}
	return i, nil
}
