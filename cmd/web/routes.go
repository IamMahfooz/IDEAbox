package main

import (
	"fmt"
	"github.com/bmizerany/pat"
	"github.com/fsnotify/fsnotify"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/about", dynamicMiddleware.ThenFunc(app.about))
	mux.Get("/faceism", dynamicMiddleware.ThenFunc(app.faceism))
	mux.Post("/recognition", dynamicMiddleware.ThenFunc(app.recognition))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))
	//mux.Post("/user/forgot", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.forget))

	mux.Get("/ping", http.HandlerFunc(ping))
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	NewfileServer := http.FileServer(http.Dir("./test/"))
	mux.Get("/test/", http.StripPrefix("/test", NewfileServer))

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// Reload the served files when a file in the static directory changes
				if event.Op&fsnotify.Write == fsnotify.Write {
					NewfileServer := http.FileServer(http.Dir("./test/"))
					mux.Get("/test/", http.StripPrefix("/test", NewfileServer))
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./test")
	if err != nil {
		panic(err)
	}

	return standardMiddleware.Then(mux)
}
