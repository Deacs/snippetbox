package main

import "chilliweb.com/snippetbox/pkg/models"

// Define a templateData type to act as the holdimg structure for
// any dynamic data that we want to pass to our HTML templates.
type templateData struct {
	Snippet *models.Snippet
}
