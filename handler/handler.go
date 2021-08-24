package handler

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

//UserTokenKey is used to access the user from context
type key int

const psKey key = 0

type Env struct {
	Port string
	Host string
}

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Status Returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

// Router We could also put *httprouter.Router in a field to not get access to the original methods (GET, POST, etc. in uppercase)
type Router struct {
	*httprouter.Router
}

//Get allows us to wrap all func calls
func (r *Router) Get(path string, handler http.Handler) {
	r.GET(path, wrapHandler(handler))
}

//Post allows us to wrap all func calls
func (r *Router) Post(path string, handler http.Handler) {
	r.POST(path, wrapHandler(handler))
}

//Put allows us to wrap all func calls
func (r *Router) Put(path string, handler http.Handler) {
	r.PUT(path, wrapHandler(handler))
}

//Delete allows us to wrap all func calls
func (r *Router) Delete(path string, handler http.Handler) {
	r.DELETE(path, wrapHandler(handler))
}

//NewRouter creates a new wrapped Router
func NewRouter() *Router {
	return &Router{httprouter.New()}
}

//WrapHandler is a little but of Glue to use with HTTPROUTER
func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctxWithParams := context.WithValue(r.Context(), psKey, ps)
		rWithPS := r.WithContext(ctxWithParams)
		h.ServeHTTP(w, rWithPS)
	}
}

type Handle func(e *Env, w http.ResponseWriter, r *http.Request, ps httprouter.Params) error

// The Handler struct that takes a configured Env and a function matching
// our useful signature.
type Handler struct {
	*Env
	H func(e *Env, w http.ResponseWriter, r *http.Request, ps httprouter.Params) error
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ps := r.Context().Value(psKey).(httprouter.Params)
	err := h.H(h.Env, w, r, ps)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

func Index(e *Env, w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	lp := filepath.Join("templates", "base.html")
	fp := filepath.Join("templates", "index.html")

	tmpl, _ := template.ParseFiles(lp, fp)
	tmpl.ExecuteTemplate(w, "base", nil)
	return nil
}

func Blog(e *Env, w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	lp := filepath.Join("templates", "base.html")
	fp := filepath.Join("templates", "blog.html")

	tmpl, _ := template.ParseFiles(lp, fp)
	tmpl.ExecuteTemplate(w, "base", nil)
	return nil
}
