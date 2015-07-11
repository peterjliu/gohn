// Concurrently downloads hacker news items.
package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"sync"
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

const waitTimeMeanMs = 10

var start = flag.Int("start", 2, "Start item to download")
var end = flag.Int("end", 100, "End item to download")
var numWorkers = flag.Int("numworkers", 50, "Number of concurrent downloads.")
var verbose = flag.Bool("verbose", false, "print a lot of info to log")

func getItems(in <-chan int, out chan<- *gohn.Item, wg *sync.WaitGroup) {
	for {
		waitTime := time.Millisecond * gohn.ExpWithMax(waitTimeMeanMs, 100.0)
		if *verbose {
			log.Printf("getitems, wait for %s", waitTime)
		}
		time.Sleep(waitTime)
		i, more := <-in
		if more {
			it, err := gohn.GetItem(i)
			if err != nil {
				log.Fatal(fmt.Sprintf("Failed to get item %d", i))
			}
			out <- it
		} else {
			// channel is depleted
			wg.Done()
			return
		}
	}
}

func saveItem(db *leveldb.DB, it *gohn.Item) {
	pbmsg, err := proto.Marshal(it)
	check(err)
	key := []byte(strconv.Itoa(int(it.GetId())))
	err = db.Put(key, pbmsg, nil)
	check(err)
}

func logItem(it *gohn.Item) {
	log.Printf("Item %d\n", it.GetId())
	if it.Title != nil {
		log.Printf("Title: %s\n", it.GetTitle())
	}
	if it.Text != nil {
		log.Printf("Text: %s\n", it.GetText())
	}
}

func main() {
	flag.Parse()
	db, err := leveldb.OpenFile("hackernewsdb", nil)
	check(err)
	defer db.Close()
	items := make(chan *gohn.Item)
	itemqueue := make(chan int)
	startTime := time.Now()
	go func() {
		for i := *start; i <= *end; i++ {
			key := []byte(strconv.Itoa(i))
			exists, err := db.Has(key, nil)
			check(err)
			if exists {
				log.Printf("skipping item %d, already in db\n", i)
				continue
			}
			log.Printf("enqueue item %d\n", i)
			itemqueue <- i
		}
		close(itemqueue)
	}()

	added := 0
	var wg sync.WaitGroup
	wg.Add(*numWorkers)
	for i := 0; i < *numWorkers; i++ {
		go getItems(itemqueue, items, &wg)
	}
	go func() {
		wg.Wait()
		close(items)
	}()

	for it := range items {
		key := []byte(strconv.Itoa(int(it.GetId())))
		pbmsg, err := proto.Marshal(it)
		err = db.Put(key, pbmsg, nil)
		check(err)
		if *verbose {
			logItem(it)
		}
		added += 1

	}
	log.Printf("%s\n", time.Since(startTime))
	log.Printf("%g items per second\n", float64(added)/float64(time.Since(startTime).Seconds()))
}
