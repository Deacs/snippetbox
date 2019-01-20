package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Change the signature of the handler so it is defined as a method against *application
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// Now using the notFound() helper
		app.notFound(w)
		return
	}

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
}

// Change the signature of the handler so it is defined as a method against *application
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		// Now using the notFound() helper
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Display a snippet with ID %d...", id)
}

// Change the signature of the handler so it is defined as a method against *application
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		// Use the clientError() helper
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a new snippet..."))
}
