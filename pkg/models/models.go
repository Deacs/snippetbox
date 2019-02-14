package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// Add a new ErrInvalidCredentials error. We'll use this if a user
	// tries to login with an invalid email address or password.
	ErrInvalidCredentials = errors.New("")
	// Add a new ErrDuplicateEmail error. We'll use this if a user
	//. tries to signup with an em,ail address that's already in use.
	ErrDuplicateEmail = errors.New("")
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}
