package gohn

import (
	"fmt"
)

func (i *Item) PrettyString() string {
	var p string
	if i.GetType() == "story" {
		p = fmt.Sprintf("%d: %s", i.GetId(), i.GetTitle())
	} else {
		p = fmt.Sprintf("unsupported type: %s", i.GetType())
	}
	return p

}
