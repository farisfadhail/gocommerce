package repository

import "fmt"

type ErrRecordNotFound struct {
	Entity string
	ID     string
}

func (e *ErrRecordNotFound) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.Entity, e.ID)
}

func NewErrRecordNotFound(entity, id string) *ErrRecordNotFound {
	return &ErrRecordNotFound{
		Entity: entity,
		ID:     id,
	}
}
