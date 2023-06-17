package util

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
