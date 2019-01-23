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
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	// Create the SQL statement to execute
	// Split over 2 lines for readibility
	stmt := `SELECT id, title, content, created, expires FROM snippets 
	WHERE expires > UTC_TIMESTAMP() AND id =?`

	// Use the QueryRow() method on the comnnection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result set from the database.
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct.
	s := &models.Snippet{}

	// Use row.Scan() to copy the values from each field in the sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place we want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by the statement. If the query returns no rows, then
	// row.Scan() will return a sql.ErrNoRows error. We check for that and return
	// our models.ErrNoRecord error instead of a Snippet object
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	// If everything went OK then return the Snippet object.
	return s, nil

	// Version above is long hand.
	// As errors from DB.QueryRow() are deferred until Scan() is called, it can be shortened to:
	// --------------------------------------
	// s := &models.Snippet{}
	// err := m.DB.QueryRow("SELECT ...", id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	// if err == sql.ErrNoRows {
	// return nil, models.ErrNoRecord }
	// else if err != nil {
	// return nil, err }
	// return s, nil
}

// This will return the 10 latest snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	// The SQL query that we want to execute.
	stmt := `SELECT id, title, content, created, expires FROM snippets 
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	// Use the Query() method on the connection pool to execute our SQL statement.
	// This returns a sql.Rows resultset containing the result of the query
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We defer rows.CLose() to ensure the sql.Rows resultset is always properly
	// closed before the Latest() method returns. This defer statement
	// should come *after* we check for an error from the Query() method.
	// Otherwise, if Query() returns an error, we'll get a panic
	// trying to close the nil resultset
	defer rows.Close()

	// Initialize an empty slice to hold the models.Snippets objects
	snippets := []*models.Snippet{}

	// Use rows.Next to iterate through the rows of the resultset.
	// This prepares the first (and then each subsequent) row to be acted upon by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection
	for rows.Next() {
		// Create a pointer to a new Zeroed Snippet struct.
		s := &models.Snippet{}
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet() object that we created. Again, the arguments to row.Scan()
		// must be pointers to the place we want to copy the data into, and the
		// number of arguments must be exactly the same number of
		// columns returned by the statement.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK, return the snippers slice
	return snippets, nil
}
