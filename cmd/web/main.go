package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include fields for the new custom loggers,
// but we'll add more as the build progresses
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls.
	// The value of he flag will be stored in the addr variable at runtime
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr variable.
	// This needs to be called *before* you use the addr variable
	// otherwise it will always contain the default valiue of "":4000".
	// If any errors are encountered during parsing teh application will be terminated
	flag.Parse()

	// Use log.New() to create a logger for writing information messages.
	// This takes three parameters: the destination to write the logs to (os.Stdout),
	// a string prefix for message (INFO followed by a tab), and flags to indicate what
	// additional information to include (local date and time).
	// Note that the flags are joined using the bitwise OR operator |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error messages in the same way. but use stderr as
	// the destination amd also ther log.Lshortfile flag to include the relevant
	// file name and line number
	// If we wanted the full file path we can use log.Llongfile instead
	// If we wanr to force UTC datetimes, we can use the log.LUTC flag
	errorLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize a new instance of application containing the dependencies
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Initialze a new http.Server struct. We will set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom erroLog Logger
	// in the event of any problems
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		// We are now using the encapsulated routes rather than defining them directly here
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	// Call the ListenAndServe() method on our new http.Server struct
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
