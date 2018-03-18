package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Initialize sets up the database connection and routes for the app
func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run starts the app and serves on the specified addr
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8000", a.Router))
}

// App represents the application
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")

	a.Router.HandleFunc("/events", a.getEvents).Methods("GET")
	a.Router.HandleFunc("/event", a.createEvent).Methods("POST")
	a.Router.HandleFunc("/event/{id:[0-9]+}", a.getEvent).Methods("GET")
	a.Router.HandleFunc("/event/{id:[0-9]+}", a.updateEvent).Methods("PUT")
	a.Router.HandleFunc("/event/{id:[0-9]+}", a.deleteEvent).Methods("DELETE")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
