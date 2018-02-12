package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type Person struct {
	ID        int      `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

var people []Person

func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.InitializeDB()
	a.initializeRoutes()
}

func (a *App) InitializeDB() {
	var err error

	os.Remove(os.Getenv("DBFILE"))
	a.DB, err = sql.Open("sqlite3", os.Getenv("DBFILE"))
	if err != nil {
		log.Fatal("sql.Open:", err)
	}

	_, err = a.DB.Exec(`CREATE TABLE IF NOT EXISTS "people" (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"firstname" VARCHAR(255),
			"lastname" VARCHAR(255),
			"city" VARCHAR(255),
			"state" VARCHAR(255))`)

	if err != nil {
		log.Fatal("CREATE TABLE: ", err)
	}

	_, err = a.DB.Exec(
		`INSERT INTO "people" ("firstname", "lastname", "city", "state") VALUES(?, ?, ?, ?)`,
		"Taro",
		"Hoge",
		"Fujisawa",
		"Kanagawa",
	)

	if err != nil {
		log.Fatal("INSERT INTO: ", err)
	}
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/people", a.GetPeopleEndpoint).Methods("GET")
	a.Router.HandleFunc("/people/{id}", a.GetPersonEndpoint).Methods("GET")
	a.Router.HandleFunc("/people/{id}", a.CreatePersonEndpoint).Methods("POST")
	a.Router.HandleFunc("/people/{id}", a.DeletePersonEndpoint).Methods("DELETE")
	a.Router.HandleFunc("/testdb", a.TestDB).Methods("GET")
}

func (a *App) Run(addr string) {
	loggedRouter := handlers.LoggingHandler(os.Stdout, a.Router)
	log.Fatal(http.ListenAndServe(addr, loggedRouter))
}

func main() {
	var app App
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// dummy data
	people = append(people, Person{ID: 1, Firstname: "foo", Lastname: "hoge", Address: &Address{City: "Fujisawa", State: "Kanagawa"}})
	people = append(people, Person{ID: 2, Firstname: "bar", Lastname: "fuga"})

	app.Initialize()
	app.Run(os.Getenv("ADDR"))
}

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
	results := []Person{}

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
