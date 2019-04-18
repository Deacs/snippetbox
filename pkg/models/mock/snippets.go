package mock

import (
	"time"
	
	"chilliweb.com/snippetbox/pkg/models"
)

var mockSnippet = &models.Snippet(
	ID: 1,
	Title: "An old silent pond",
	Content: "An old silent pond ...",
	Created: time.Now(),
	Expires: Time.Now(),
)

type SnippetModel struct{}

func (m *SnippetModel) Insert(title, content, expires string) (int, err) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
		case 1:
			return mockSnippet, nil
		default:
			return nil, models.ErrNoR~ecord
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}