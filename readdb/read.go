// Reads hacker news leveldb created by ../download
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/peterjliu/gohn"
	"github.com/syndtr/goleveldb/leveldb"
)

var dbpath = flag.String("dbpath", "", "path to hacker news leveldb database")

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	flag.Parse()
	db, err := leveldb.OpenFile(*dbpath, nil)
	check(err)
	defer db.Close()
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		var it gohn.Item
		err := proto.Unmarshal(iter.Value(), &it)
		check(err)
		if it.GetTitle() != "" {
			fmt.Printf("%d: %s\n", it.GetId(), it.GetTitle())
		}
	}
}
