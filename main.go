package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Define a home handle function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exacrly matches "/" If it doesn't,
	// use the http.NotFound() function to send a 404 to the client
	// Importantly, we then return from the handler. If we don't return, the handler
	// would keep executing and also write the "Hello from Snippetbox" message.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Snippetbox"))
}

// Add a showSnippet handler
func showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id from the query string and try to
	// convert it to an integer using the strconv.Atoi() function.
	// If it can't be converted to an integer, or the value is less than 1,
	// We return a 404 page not found response
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Use the fmt.Fprintf() function to interpolate teh id value with our response
	// and write it to the http.ResponseWriter
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// Add a createSnippet handler
func createSnippet(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not.
	// If it's not, use the w.WriteHeader() method to send a 405 status code and
	// the w.Write() to write a "Method Not Allowed" response body.
	// We then return from the function so the subsequent code is not executed
	// Use curl -i -X PUT http://localhost:4000/snippet/create to test
	if r.Method != "POST" {
		// Use the Header().Set() method to add an 'Allow: Post' header to the
		// response header map. The first parameter is the header name,
		// the second parameter is the header value
		// All changes to the header map must be made BEFORE calling WriteHeader or Write
		// If not, they are ignored
		w.Header().Set("Allow", "POST")
		// Use the http.Error() function to send a "Method Not Allowed" string and
		// 405 status code as the response body
		// The ResponseWriter is also sent to allow the helper to send the response
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	w.Write([]byte("Create a new snippet..."))
}

func main() {
	// Use the http.NewServeMux() function to initialise a new servemux, then
	// Register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Use the http.ListenAndServe() function to start a new web server.
	// We pass on two params: teh TCP network address to listen on (in our case ":4000")
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and exit
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
