package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: mo matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
	Byy     string
}
type User struct {
	ID              int
	Name            string
	Email           string
	HashdedPassword []byte
	Created         time.Time
	Active          bool
}
type Stranger struct {
	Name   string
	Year   string
	Place2 string
	Branch string
	Dob    string
}
