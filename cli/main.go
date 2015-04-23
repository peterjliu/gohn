package main

import (
	"fmt"
	"log"

	"github.com/peterjliu/gohn"
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}
func main() {
	items, err := gohn.TopStories()
	check(err)
	for _, i := range items {
		fmt.Println(i.PrettyString())
	}
}
