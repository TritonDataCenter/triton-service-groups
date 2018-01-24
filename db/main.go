package db

import (
	"log"

	"github.com/jackc/pgx"
)

var database = make(map[string]interface{})

func FindBy(key string) (interface{}, bool) {
	com, ok := database[key]

	return com, ok
}

func Save(key string, item interface{}) {
	database[key] = item
}

func Remove(key string) {
	delete(database, key)
}

func Test(key string) {
	config := pgx.ConnConfig{
		Host:     "",
		Database: "",
	}

	log.Printf("%v", config)
}
