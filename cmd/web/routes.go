package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	dynamicMiddleware := alice.New(app.session.Enable, noSurf)

	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/", dynamicMiddleware.ThenFunc(app.home))

	router.Handler(http.MethodGet, "/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	router.Handler(http.MethodPost, "/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	router.Handler(http.MethodGet, "/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	router.Handler(http.MethodPost, "/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	router.Handler(http.MethodPost, "/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	// Middleware for every request - 'standard'.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standardMiddleware.Then(router)
}
