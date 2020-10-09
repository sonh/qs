package main

import (
	"fmt"
	"github.com/sonh/qs"
	"time"
)

type User struct {
	Verified bool      `qs:"verified"`
	From     time.Time `qs:"from,millis"`
}

type Query struct {
	User User `qs:"user"`
}

func main() {
	query := Query{
		User: User{
			Verified: true,
			From:     time.Now(),
		},
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
