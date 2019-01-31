package main

import "net/http"

// Updated signature for the routes() method as it now returns
// a http.Handler instead of the original *http.ServeMux
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Wrap the exisiting chain in the recoverPanic() middleware.
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
