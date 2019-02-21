package main

import (
	"fmt"
	"net/http"
	"strconv"

	"chilliweb.com/snippetbox/pkg/forms"
	"chilliweb.com/snippetbox/pkg/models"
)

// Change the signature of the handler so it is defined as a method against *application
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// The Pat router explicitly matches the "/" path exactly, we can now remove the manual check
	// of r.URL.Path != "/" from this handler.

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
	// The colon isn't automatically stripped from the named capture key, so we need
	// to ge the value of ":id" from the query string instead of "id"
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
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

	// Flash message is now added to the addDefaultData helper
	// to make it available to all views.

	// Use the render helper with the newly created created flash message
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template
		Form: forms.New(nil),
	})
}

// Changed the signature of the handler so it is defined as a method against *application
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// First we call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way for PUT and PATCH
	// requests. If there are any errors we use our app.ClientError helper to send
	// a 400 Bad REquest response to the user
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	// If the form isn't valid, redisplay the template passing in the
	// form.Form object as the data.
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	// Because the form data (with type url.Values) has been anonymously embedded
	// in the form.Form struct, we use the Get() method to retrieve
	// the validated value fro a particular form field.
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Uset the Put() method to add a string value ("Snippet successfully created!")
	// and the corresponding key ("flash") to the session data.
	// Note that if there's nbo existing session for the current user
	// (or their session has expired) then a new, empty session will be automatically
	// created for them by the session middleware.
	app.session.Put(r, "flash", "Snippet successfully created!")

	// Redirect to view the newly created snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using the form helper we made earlier
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	// If there are any errors, redisplay the signup form
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	// Try to create a new user in the database. If the email address already exists
	// add an error message to the form and redisplay it
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked and asking them to login.
	app.session.Put(r, "flash", "Your signup was successful. Please login")

	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display the user login form...")
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login user...")
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
