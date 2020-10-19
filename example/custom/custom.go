package main

import (
	"fmt"
	"github.com/sonh/qs"
)

type NullableName struct {
	First string
	Last  string
}

func (n NullableName) EncodeParam() (string, error) {
	return n.First + n.Last, nil
}

func (n NullableName) IsZero() bool {
	return n.First == "" && n.Last == ""
}

func main() {

	type Struct struct {
		User  NullableName `qs:"user"`
		Admin NullableName `qs:"admin,omitempty"`
	}

	s := Struct{
		User: NullableName{
			First: "son",
			Last:  "huynh",
		},
	}
	encoder := qs.NewEncoder()

	values, err := encoder.Values(&s)
	if err != nil {
		// Handle error
		fmt.Println("failed")
		return
	}
	fmt.Println(values.Encode())
}
