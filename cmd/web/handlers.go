package main

import (
	"fmt"
	"net/http"
	"strconv"

	"chilliweb.com/snippetbox/pkg/models"
)

// Change the signature of the handler so it is defined as a method against *application
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// Now using the notFound() helper
		app.notFound(w)
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the render helper
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

// Change the signature of the handler so it is defined as a method against *application
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		// Now using the notFound() helper
		app.notFound(w)
		return
	}

	// Use the SnippetModel object's Get method to retrieve data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found resource
	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the render helper
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

// Change the signature of the handler so it is defined as a method against *application
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		// Use the clientError() helper
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Create some variables holding dummy data
	// We'll remove these later
	title := "O Jessie"
	content := "O Jessie\nSo many ghosts\nBut forget, forget!\n\n- Walter White"
	expires := "7"

	// Pass the data to the SnippetModel.Insert() method,
	// receiving the ID of the new record back
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
