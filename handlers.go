package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Id    int
	Name  string
	Email string
	Age   int
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println("server failed to handle", err)
		return
	}
	defer rows.Close()

	data := make([]User, 0)
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Age)
		if err != nil {
			fmt.Println("server failed to handle", err)
		}
		data = append(data, user)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("server failed to handle", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	user := User{}
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Age)

	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("server failed to handle", err)
		return
	}

	_, err = db.Exec("INSERT INTO users (name, email, age) VALUES ($1, $2, $3)", user.Name, user.Email, user.Age)
	if err != nil {
		fmt.Println("failed to insert", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	up := User{}
	err := json.NewDecoder(r.Body).Decode(&up)
	if err != nil {
		fmt.Println("server failed to handle", err)
		return
	}

	row := db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	u := User{}
	err = row.Scan(&u.Id, &u.Name, &u.Email, &u.Age)

	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return

	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	if up.Name != "" {
		u.Name = up.Name
	}

	if up.Email != "" {
		u.Email = up.Email
	}

	if up.Age != 0 {
		u.Age = up.Age
	}

	_, err = db.Exec("UPDATE users SET name = $1, email = $2, age = $3 WHERE id = $4", u.Name, u.Email, u.Age, u.Id)
	if err != nil {
		fmt.Println("failed to update", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")

	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		fmt.Println("failed to delete", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}