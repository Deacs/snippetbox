package mysql

import (
	"database/sql"

	"chilliweb.com/snippetbox/pkg/models"
)

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	// Create the SQL statement we want to execute. It's split over several lines
	// for readability - so it's surrounded by backquotes instead of normal double quotes
	stmt := `INSERT INTO snippets (title, content, created, expires)
	values (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP, INTERVAL ? DAY))`

	// Use the Exec() method on the embedded connection pool to execute the statement.
	// The first parameter is the SQL statement followed by the table fields.
	// The method returns a sql.Result object which contains some basic information
	// about what happened when the statement was executed
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId() method on the result object to get the ID of our
	// newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has the type int64 so we convert it to an int type before returning
	return int(id), nil

	return 0, nil
}

// This will return a specific snippert based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

// This will return the 10 latest snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
