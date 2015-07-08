// Downloads hacker news items
package main

import (
	"log"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/peterjliu/gohn"
	"github.com/syndtr/goleveldb/leveldb"
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

const waitTimeMs = 5

func main() {
	db, err := leveldb.OpenFile("hackernewsdb", nil)
	check(err)
	check(err)
	defer db.Close()
	startTime := time.Now()
	added := 0
	for i := 500; i < 1000; i++ {
		key := []byte(strconv.Itoa(i))
		exists, err := db.Has(key, nil)
		check(err)
		if exists {
			log.Printf("skipping item %d, already in db\n", i)
			continue
		}
		it, err := gohn.GetItem(i)
		if err != nil {
			log.Fatal("Failed to get item %d", i)
		}
		if it.Title != nil {
			log.Printf("Title: %s\n", *it.Title)
		}
		if it.Text != nil {
			log.Printf("Text: %s\n", *it.Text)

		}
		pbmsg, err := proto.Marshal(it)
		check(err)
		err = db.Put(key, pbmsg, nil)
		added += 1
		check(err)
		time.Sleep(time.Millisecond * waitTimeMs)
		log.Printf("Added item %d\n", i)
	}
	log.Printf("%s\n", time.Since(startTime))
	log.Printf("%g items per second\n", float64(added)/float64(time.Since(startTime).Seconds()))
}
