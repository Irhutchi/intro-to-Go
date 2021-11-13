package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// user represents data about a record user.
type User struct {
	Id    string `json:"id"`
	Fname string `json:"fname"`
	Lname string `json:"lname"`
	Uname string `json:"uname"`
}

// Init users var as a slice User struct
var (
	users []User
)

func init() {
	users = []User{
		// Hardcoded data - @todo: add database
		{Id: "1", Fname: "John", Lname: "Doe", Uname: "Jdoe"},
		{Id: "2", Fname: "Steve", Lname: "Smith", Uname: "Steve"},
	}
}

// Get all users
func getUsers(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	result, err := json.Marshal(users)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{error: "Error marshalling the users array"}`))
		return
	}
	resp.WriteHeader(http.StatusOK)
	//write the result of the users array
	resp.Write(result)
	// json.NewEncoder(resp).Encode(users)
}

// Get single user
func getUser(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req) // Gets params
	// Loop through users and find one with the id from the params
	for _, item := range users {
		if item.Id == params["id"] {
			json.NewEncoder(resp).Encode(item)
			return
		}
	}
	json.NewEncoder(resp).Encode(&User{})
}

// Add new user
func createUser(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	var user User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{error: "Error marshalling the request to add user"}`))
		return
	}
	// convert random Int Id to a string
	user.Id = strconv.Itoa(rand.Intn(100000000)) // Mock Id - not safe
	users = append(users, user)
	resp.WriteHeader(http.StatusOK)
	result, err := json.Marshal(user)
	resp.Write(result) //write the result of the users array

	// add user to server memory
	// json.NewEncoder(resp).Encode(user)
}

// Update user
func updateUser(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	for index, item := range users {
		if item.Id == params["id"] {
			users = append(users[:index], users[index+1:]...)
			var user User
			_ = json.NewDecoder(req.Body).Decode(&user)
			user.Id = params["id"]
			users = append(users, user)
			json.NewEncoder(resp).Encode(user)
			return
		}
	}
}

// Delete user
func deleteUser(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	for index, item := range users {
		if item.Id == params["id"] {
			users = append(users[:index], users[index+1:]...)
			break
		}
	}
	json.NewEncoder(resp).Encode(users)
}
