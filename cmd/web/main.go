package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"chilliweb.com/snippetbox/pkg/models/mysql"

	// If we try to import this normally the Go compiler will raise an error.
	// However, we need the driver's init() function to run so that it can register itself with the database/sql package.
	// The trick to getting around this is to alias the package name to the blank identifier
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type contextKey string

var contextKeyUser = contextKey("user")

// Define an application struct to hold the application-wide dependencies for the
// web application. User Model has now been added
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
	users         *mysql.UserModel
}

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls.
	// The value of he flag will be stored in the addr variable at runtime
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Define a new command-line flag for the MySQL DSN string
	dsn := flag.String("dsn", "web:jiwa@/snippetbox?parseTime=true", "MySQL data source name")

	// Define a new command-line flag for the session secret (a random key which
	// will be used to encrypt and authenticate session cookies).
	// It should be 32 bytes long
	secret := flag.String("secret", "j@883r_w0c|<%-@_pO3m4alL8|`|`rQD", "Secret key")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr variable.
	// This needs to be called *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000".
	// If any errors are encountered during parsing the application will be terminated
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
	// If we want to force UTC datetimes, we can use the log.LUTC flag
	errorLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	// To keep the main() fnction tidy, the code for creating a connection pool
	// has been put into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// We also defer a call to db.Close() so that the connection pool is closed
	// before the main() function exits.
	// This is a bit superfluous. Our application is only ever terminated by a signal interrupt
	// (i.e. Ctrl+c) or by errorLog.Fatal().
	// In both of those cases, the program exits immediately and deferred functions are never run
	defer db.Close()

	// Initialize a new template cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Use the session.New() function to initialize a new session manager,
	// passing in the scret key as the parameter. Then we configure it so
	// sessions always expire after 12 hours
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true // Set the Secure flag on session cookies

	// Initialize a new instance of application containing the dependencies
	app := &application{
		// Logging dependencies
		errorLog: errorLog,
		infoLog:  infoLog,
		// Session management
		session: session,
		// Add the mysql.SnippetModel instance to the dependencies
		snippets: &mysql.SnippetModel{DB: db},
		// Add the template cache to the dependencies
		templateCache: templateCache,
		// Initialize a mysql.UserModel instance and add to the dependencies
		users: &mysql.UserModel{DB: db},
	}

	// Initialize A tls.Config struct to hold the non-default TLS settings
	// we want the server to use
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Initialze a new http.Server struct. We will set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog Logger
	// in the event of any problems
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
		// Set the server's TLSConfig field to use the variable just created
		TLSConfig: tlsConfig,
		// Adding Idle, Read & Write timeouts to the server
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Call the ListenAndServe() method on our new http.Server struct
	infoLog.Printf("Starting server on %s", *addr)
	// Using the ListenAndServeTLS() method to start the HTTPS server.
	// We pass in the paths to the TLS certificate and corresponding
	// private key as the two parameters
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

// The OpenDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
