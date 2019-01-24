package main

import (
	"fmt"
	"html/template"
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

	// Temporarily dump the latest snippets directly to the page
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
	}

	for _, snippet := range s {
		fmt.Fprintf(w, "%v\n", snippet)
	}

	/*

		// Initialize a slice containing the paths to the two files.
		// Note tha the home.page.tmpl file must be the *first* file in the slice
		files := []string{
			"./ui/html/home.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/footer.partial.tmpl",
		}

		// Use the template.ParseFiles() function to read the template file into a template set.
		// If there's an error, we log the detailed error message and use the
		// http.Error() function to send a generic 500 Internal Server Error response to the user
		ts, err := template.ParseFiles(files...)
		if err != nil {
			// Now using the serverError() helper
			app.serverError(w, err)
			return
		}

		// We then we use the Execute() method on the template set to write the template
		// content as the response body. The last parameter to Execute() represents any
		// dynamic data that we want to pass in, which for now, we'll leave as nil
		err = ts.Execute(w, nil)
		if err != nil {
			// Now using the serverError() helper
			app.serverError(w, err)
		}

	*/
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

	// Initialize a slice containing the paths to the show.page.tmpl file,
	// plus the base layout and footer partials
	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Parse the template files ...
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Execute the templates
	// The snippet data (a models.Snippet.struct) is passed in as the final parameter
	err = ts.Execute(w, s)
	if err != nil {
		app.serverError(w, err)
	}
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
	title := "O snail"
	content := "O snail\nClimb Mount Fuji\nBut slowly, slowly!\n\n- Kobayashi Issa"
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
