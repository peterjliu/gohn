package gohn

import (
	"fmt"
)

func (i *Item) PrettyString() string {
	var p string
	if *i.Type == "story" {
		p = fmt.Sprintf("%d: %s", *i.Id, *i.Title)
	} else {
		p = "error"
	}
	return p

}
