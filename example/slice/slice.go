package main

import (
	"fmt"
	"github.com/sonh/qs"
)

type Query struct {
	Default []string `qs:"default_fmt"`
	Comma   []string `qs:"comma_fmt,comma"`
	Bracket []string `qs:"bracket_fmt,bracket"`
	Index   []string `qs:"index_fmt,index"`
}

func main() {
	tags := []string{"go", "docker"}
	query := &Query{
		Default: tags,
		Comma:   tags,
		Bracket: tags,
		Index:   tags,
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
