package main

import (
	"fmt"
	"github.com/sonh/qs"
	"time"
)

type Query struct {
	Tags   []string  `qs:"tags"`
	Limit  int       `qs:"limit"`
	From   time.Time `qs:"from"`
	Open   bool      `qs:"open,int"`
	Active bool      `qs:"active,omitempty"` //omit empty value
	Ignore float64   `qs:"-"`                //ignore
}

func main() {
	query := &Query{
		Tags:   []string{"docker", "golang", "reactjs"},
		Limit:  24,
		From:   time.Unix(1580601600, 0).UTC(),
		Ignore: 0,
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
