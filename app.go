package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

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
	defer a.DB.Close()

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
