package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf" // CSRF management
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		// Any code here will execute on the way DOWN the chain
		next.ServeHTTP(w, r)
		// Any code here will execute on the way back UP the chain
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function
		// This will always be run in the event of a panic as Go unwinds the stack
		defer func() {
			// Use the builtin recover() function to check if there has been a panic
			// If there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header to the response.
				w.Header().Set("Connection", "close")
				// Call the app.ServerError helper method to return a
				// 500 Internal Server Error response
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authenticated, redirect them to the login page and
		// return from the middleware chain so no subsequent handlers in
		// the chain are executed
		if app.authenticatedUser(r) == 0 {
			http.Redirect(w, r, "/user/login", 302)
			return
		}

		// Otherwise call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// Create a NoSurf middleware function which uses a customized CSRF cookie with
// the Secure, OPath and HttpOnly flags set.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}
