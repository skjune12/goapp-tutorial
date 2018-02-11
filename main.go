package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type App struct {
	Router *mux.Router
	// DB     *sql.DB
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var people []Person

func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/people", a.GetPeopleEndpoint).Methods("GET")
	a.Router.HandleFunc("/people/{id}", a.GetPersonEndpoint).Methods("GET")
	a.Router.HandleFunc("/people/{id}", a.CreatePersonEndpoint).Methods("POST")
	a.Router.HandleFunc("/people/{id}", a.DeletePersonEndpoint).Methods("DELETE")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func main() {
	var app App

	// dummy data
	people = append(people, Person{ID: "1", Firstname: "foo", Lastname: "hoge", Address: &Address{City: "Fujisawa", State: "Kanagawa"}})
	people = append(people, Person{ID: "2", Firstname: "bar", Lastname: "fuga"})

	app.Initialize()
	app.Run(":8080")
}

func (a *App) GetPersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	// returns empty object
	json.NewEncoder(w).Encode(&Person{})
}

func (a *App) GetPeopleEndpoint(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

func (a *App) CreatePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	var person Person
	params := mux.Vars(r)

	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)

	json.NewEncoder(w).Encode(people)
}

func (a *App) DeletePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			// everything before and everything after.
			people = append(people[:index], people[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(people)
}
