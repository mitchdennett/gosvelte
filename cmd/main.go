package main

import (
	"log"
	"net/http"
	"os"

	handler "github.com/mitchdennett/gosvelte/handler"
)

func main() {

	router := handler.NewRouter()

	router.ServeFiles("/static/*filepath", http.Dir("static/public"))
	log.Println(http.Dir("static/public"))

	env := &handler.Env{
		Port: os.Getenv("PORT"),
		Host: os.Getenv("HOST"),
	}

	router.Get("/", handler.Handler{Env: env, H: handler.Index})
	router.Get("/blog", handler.Handler{Env: env, H: handler.Blog})

	log.Fatal(http.ListenAndServe(":8000", router))

}
