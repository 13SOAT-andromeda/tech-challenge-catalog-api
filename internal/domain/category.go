package domain

import (
	"errors"
	"time"
)

type Category struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCategory(name string) (*Category, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	return &Category{
		Name: name,
	}, nil
}
