package main

import (
	"log"
	"net/http"

	"github.com/gogo-framework/router"
)

func usersListHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List of users"))
}

func usersGetHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("id")
	w.Write([]byte("User ID: " + userId))
}

func usersCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create user"))
}

func usersStoreHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Store user"))
}

func usersEditHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Edit user"))
}

func usersUpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update user"))
}

func usersDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete user confirm"))
}

func usersDeletePerformHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete user perform"))
}

func main() {
	r := router.NewRouter()

	r.Group("/api/users", func(r *router.Router) {
		r.GET("/", usersListHandler)
		r.GET("/{id}", usersGetHandler)
		r.GET("/create", usersCreateHandler)
		r.POST("/create", usersStoreHandler)
		r.GET("{id}/edit", usersEditHandler)
		r.POST("{id}/edit", usersUpdateHandler)
		r.GET("{id}/delete", usersDeleteHandler)
		r.POST("{id}/delete", usersDeletePerformHandler)
	})

	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
