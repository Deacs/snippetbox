package main

import (
	"log"
	"net/http"
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
	w.Write([]byte("Display a specific snippet..."))
}

// Add a createSnippet handler
func createSnippet(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not.
	// If it's not, use the w.WriteHeader() method to send a 405 status code and
	// the w.Write() to write a "Method Not Allowed" response body.
	// We then return from the function so the subsequent code is not executed
	// Use curl -i -X PUT http://localhost:4000/snippet/create
	if r.Method != "POST" {
		// Use the Header().Set() method to add an 'Allow: Post' header to the
		// response header map. The first parameter is the header name,
		// the second parameter is the header value
		// All changes to the header map must be made BEFORE calling WriteHeader or Write
		// If not, they are ignored
		w.Header().Set("Allow", "POST")
		// Can obviously be any response header map
		w.Header().Set("Foo", "BAR")
		// WriteHeader can only be called once per response
		// If it is not present, a 200 OK will automatically be sent
		// This also the case if Write is called first
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
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
	log.Println("Staring server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
