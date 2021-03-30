package server

import (
	"api/database"
	"api/response"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

//GetUsers get all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connection()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, errors.New("an error occurred while connecting to database"))
		return
	}

	defer db.Close()

	var sql string = "SELECT * FROM users ORDER BY id DESC"

	lines, err := db.Query(sql)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, errors.New("an error occurred while fetching users"))
		return
	}

	defer lines.Close()

	var users []user

	for lines.Next() {
		var user user

		if err := lines.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			response.Erro(w, http.StatusInternalServerError, errors.New("an error occurred while scaning user"))
			return
		}
		users = append(users, user)
	}

	response.JSON(w, http.StatusOK, users)
}

//GetUser get user
func GetUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	id, err := strconv.ParseUint(parameters["id"], 10, 32)
	if err != nil {
		response.Erro(w, http.StatusUnprocessableEntity, errors.New("an error occurred while converting ID"))
		return
	}

	db, err := database.Connection()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, errors.New("an error occurred while connecting to database"))
		return
	}

	defer db.Close()

	var sql string = "SELECT * FROM users WHERE id = ?"

	line, err := db.Query(sql, id)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, errors.New("an error occurred while fetching user"))
		return
	}

	var user user
	if line.Next() {
		if err := line.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			response.Erro(w, http.StatusInternalServerError, errors.New("an error occurred while scaning user"))
			return
		}
	}

	if user.ID == 0 {
		response.JSON(w, http.StatusOK, "there is no user with this id")
		return
	}
	response.JSON(w, http.StatusOK, user)
}

// UpdateUser update user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	id, err := strconv.ParseUint(parameters["id"], 10, 32)
	if err != nil {
		response.Erro(w, http.StatusUnprocessableEntity, errors.New("an error occurred while converting ID"))
		return
	}

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

	db, err := database.Connection()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, errors.New("an error occurred while connecting to database"))
		return
	}

	defer db.Close()

	var sql string = "UPDATE users SET name = ?, email = ? WHERE id = ?"
	statement, err := db.Prepare(sql)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, errors.New("an error occurred while creating statement"))
		return
	}

	defer statement.Close()

	if _, err := statement.Exec(user.Name, user.Email, id); err != nil {
		response.Erro(w, http.StatusInternalServerError, errors.New("an error occurred while update user"))
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
