package mysql

import (
	"database/sql"

	"chilliweb.com/snippetbox/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

// We'll use the Insert method to add a new record to the users table
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// We'll use the Authenticate method to verify whether a user exisits with
// the provided email address and password.
// This will return a user ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// We'll use the Get method to fetch details for a specific user
// based on their ID
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
