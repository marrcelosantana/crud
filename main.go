package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	// Iniciando a conexão com o banco de dados.
	db, err = sql.Open("postgres", "postgres://postgres:postgres@postgres_database:5432/crud?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database")
}

func main() {
	http.HandleFunc("/users/read", GetAll)
	http.HandleFunc("/users/readById", GetUserByID)
	http.HandleFunc("/users/create", Create)
	http.HandleFunc("/users/update", Update)
	http.HandleFunc("/users/delete", Delete)
	http.ListenAndServe(":8080", nil)
}