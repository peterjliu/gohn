package main

import (
	"fmt"
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
func main() {
	fmt.Println("------comments------")
	db, err := leveldb.OpenFile("hackernewsdb", nil)
	check(err)
	defer db.Close()
	for i := 2; i < 100; i++ {
		it, err := gohn.GetItem(i)
		check(err)
		pbmsg, err := proto.Marshal(it)
		err = db.Put([]byte(strconv.Itoa(i)), pbmsg, nil)
		check(err)
		time.Sleep(time.Millisecond * 100)
	}
}
