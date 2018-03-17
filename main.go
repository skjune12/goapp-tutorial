package main

import (
	"log"
	"os"

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

var people []Person

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
