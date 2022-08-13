package errors

import "fmt"

type IdNotFound struct {
	Field string
	Id    int
}

func (e IdNotFound) Error() string {
	return fmt.Sprintf("%s[%d] not found", e.Field, e.Id)
}

func ErrIdNotFound(field string, id int) IdNotFound {
	return IdNotFound{
		Field: field,
		Id:    id,
	}
}
