package errors

import "fmt"

type UniqueValue struct {
	Field string
}

func (e UniqueValue) Error() string {
	return fmt.Sprintf("%s already exist", e.Field)
}

func ErrUniqueValue(field string) UniqueValue {
	return UniqueValue{
		Field: field,
	}
}
