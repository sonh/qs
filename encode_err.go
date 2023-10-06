package qs

import "fmt"

type ErrInvalidInput struct {
	inputType string
}

func (err ErrInvalidInput) Error() string {
	return fmt.Sprintf("expects struct input, got %v", err.inputType)
}
