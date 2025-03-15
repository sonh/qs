package qs

import (
	"fmt"
	"reflect"
)

type InvalidInputErr struct {
	InputKind reflect.Kind
}

func (e InvalidInputErr) Error() string {
	return fmt.Sprintf(`input should be struct type, got "%v"`, e.InputKind)
}
