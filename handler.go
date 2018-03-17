package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func (a *App) GetPersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		id, _ := strconv.Atoi(params["id"])
		if item.ID == id {
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
	id, _ := strconv.Atoi(params["id"])
	person.ID = id
	people = append(people, person)

	json.NewEncoder(w).Encode(people)
}

func (a *App) DeletePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		id, _ := strconv.Atoi(params["id"])
		if item.ID == id {
			// everything before and everything after.
			people = append(people[:index], people[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(people)
}

func (a *App) TestDB(w http.ResponseWriter, r *http.Request) {
	var err error
	results := []Person{}

	a.DB, err = sql.Open("sqlite3", os.Getenv("DBFILE"))
	if err != nil {
		log.Fatal("sql.Open:", err)
	}
	defer a.DB.Close()

	rows, err := a.DB.Query(`SELECT * from people`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Person
		var a Address
		err = rows.Scan(&p.ID, &p.Firstname, &p.Lastname, &a.City, &a.State)
		if err != nil {
			log.Fatal(err)
		}
		p.Address = &a

		results = append(results, p)
	}

	json.NewEncoder(w).Encode(results)
}