package util

import "encoding/json"

func ToPointer[T any](val T) *T {
	return &val
}

func FromPointer[T any](val *T) T {
	var newVal T
	if val != nil {
		newVal = *val
	}
	return newVal
}

func WrapEvent[T any](handler func(*T) error) func([]byte) error {
	return func(a []byte) error {
		var ev *T
		if err := json.Unmarshal(a, &ev); err != nil {
			return err
		}
		if err := handler(ev); err != nil {
			return err
		}
		return nil
	}
}
