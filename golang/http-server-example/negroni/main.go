package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type trivial struct {
}

func (t *trivial) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("Executing trivial middleware")
	next(w, r)
}

type badAuth struct {
	Username string
	Password string
}

func (b *badAuth) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	if username != b.Username || password != b.Password {
		http.Error(w, "Unaothorized", http.StatusUnauthorized)
		return
	}
	ctx := context.WithValue(r.Context(), "username", username)
	r = r.WithContext(ctx)
	next(w, r)
}

func hello(w http.ResponseWriter, r *http.Request) {
	usarname := r.Context().Value("username").(string)
	fmt.Fprintf(w, "Hello %s\n", usarname)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", hello).Methods("GET")
	
	n := negroni.Classic()
	n.Use(&trivial{})
	n.Use(&badAuth{Username: "admin", Password: "password"})
	n.UseHandler(r)

	http.ListenAndServe(":8000", n)
}
