package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	// initialize a new httprouter router instance
	router := httprouter.New()

	// here we convert default notFound response of router to our custom method
	// in order to send JSON response
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// also convert methodNotAllowed response to our custom method
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// now register relevant methods and handlers for our endpoints
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/books", app.createBookHandler)
	router.HandlerFunc(http.MethodGet, "/v1/books/:id", app.showBookHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/books/:id", app.updateBookHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/books/:id", app.deleteBookHandler)
	router.HandlerFunc(http.MethodGet, "/v1/books", app.listBooksHandler)

	router.HandlerFunc(http.MethodPost, "/v1/cart", app.addToCartHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/cart", app.deleteBookFromCartHandler)
	router.HandlerFunc(http.MethodGet, "/v1/cart", app.listBooksInCartHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// return router instance
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
