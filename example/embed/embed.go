package main

import (
	"fmt"
	"time"

	"github.com/sonh/qs"
)

type User struct {
	Verified bool           `qs:"verified"`
	From     time.Time      `qs:"from,millis"`
	Roles    map[string]int `qs:"roles"`
}

type Query struct {
	// The "dot" option encodes the nested struct using dot-notation keys
	// (user.verified) instead of the default bracket scoping (user[verified]).
	User User `qs:"user,dot"`
}

func main() {
	query := Query{
		User: User{
			Verified: true,
			From:     time.Now(),
			Roles: map[string]int{
				"admin":  1,
				"member": 2,
			},
		},
	}

	encoder := qs.NewEncoder()
	values, err := encoder.Values(query)
	if err != nil {
		// Handle error
		fmt.Println("failed")
		return
	}
	// (unescaped) output:
	// user.from=1601623397728&user.roles[admin]=1&user.roles[member]=2&user.verified=true
	fmt.Println(values.Encode())
}
