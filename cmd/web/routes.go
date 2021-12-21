package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("./ui/static"))
	router.HandlerFunc(http.MethodGet, "/", app.home)

	// Middleware for every request - 'standard'.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standardMiddleware.Then(router)
}
