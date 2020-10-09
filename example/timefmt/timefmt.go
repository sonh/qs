package main

import (
	"fmt"
	"github.com/sonh/qs"
	"time"
)

type Query struct {
	Default time.Time `qs:"default_fmt"`
	Second  time.Time `qs:"second_fmt,second"` //use `second` option
	Millis  time.Time `qs:"millis_fmt,millis"` //use `millis` option
}

func main() {
	t := time.Unix(1580601600, 0).UTC()
	query := &Query{
		Default: t,
		Second:  t,
		Millis:  t,
	}

	encoder := qs.NewEncoder()
	values, err := encoder.Values(query)
	if err != nil {
		// Handle error
		fmt.Println("failed")
		return
	}
	fmt.Println(values.Encode())
}
