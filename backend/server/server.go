package server

import (
	"api/database"
	"api/response"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type user struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

//CreateUser create user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		response.Erro(w, http.StatusUnprocessableEntity, errors.New("an error occurred while reading the request body"))
		return
	}

	var user user
	if err = json.Unmarshal(body, &user); err != nil {
		response.Erro(w, http.StatusUnprocessableEntity, errors.New("an error occurred when converting user to struct"))
		return
	}

	fmt.Println(user)

	db, err := database.Connection()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, errors.New("an error occurred while connecting to database"))
		return
	}

	defer db.Close()

	var sql string = "INSERT into users (name, email) VALUES (?, ?)"
	statement, err := db.Prepare(sql)

	if err != nil {
		response.Erro(w, http.StatusInternalServerError, errors.New("an error occurred while creating statement"))
		return
	}

	defer statement.Close()

	insert, err := statement.Exec(user.Name, user.Email)
	if err != nil {
		w.Write([]byte("an error occurred while running statement"))
		return
	}

	lastId, err := insert.LastInsertId()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, errors.New("an error occurred while getting id"))
		return
	}
	user.ID = uint32(lastId)
	response.JSON(w, http.StatusCreated, user)
}
