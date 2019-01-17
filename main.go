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