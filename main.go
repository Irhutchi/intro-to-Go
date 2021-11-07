package main

// ref: https://github.com/bradtraversy/go_restapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// user represents data about a record user.
type User struct {
	ID    string `json:"id"`
	Fname string `json:"fname"`
	Lname string `json:"lname"`
	Uname string `json:"uname"`
}

// albums slice to seed record user data.
// var albums = []user{
// 	{ID: "1", Lname: "Blue Train", Artist: "John Coltrane", Price: 56.99},
// 	{ID: "2", Lname: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
// 	{ID: "3", Lname: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
// }

// Init users var as a slice User struct
var users []User

func setupRoutes() {
	// Init router
	fmt.Println("Sever started on port 8000")
	request := mux.NewRouter()
	// Route handles & endpoints
	request.HandleFunc("/users", getUsers).Methods("GET")
	request.HandleFunc("/users/{id}", getUser).Methods("GET")
	request.HandleFunc("/users", createUser).Methods("POST")
	request.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	request.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	request.HandleFunc("/upload", fileUpload)
	request.HandleFunc("/upload1", uploadFile)

	// Start server
	log.Fatal(http.ListenAndServe(":8000", request)) // set listen port
}

// Main function
func main() {

	// Hardcoded data - @todo: add database
	users = append(users, User{ID: "1", Fname: "John", Lname: "Doe", Uname: "Jdoe"})
	users = append(users, User{ID: "2", Fname: "Steve", Lname: "Smith", Uname: "Steve"})

	setupRoutes()

}

// Get all users
func getUsers(write http.ResponseWriter, request *http.Request) {
	write.Header().Set("Content-Type", "application/json")
	json.NewEncoder(write).Encode(users)
}

// Get single user
func getUser(write http.ResponseWriter, request *http.Request) {
	write.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request) // Gets params
	// Loop through users and find one with the id from the params
	for _, item := range users {
		if item.ID == params["id"] {
			json.NewEncoder(write).Encode(item)
			return
		}
	}
	json.NewEncoder(write).Encode(&User{})
}

// Add new user
func createUser(write http.ResponseWriter, request *http.Request) {
	write.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(request.Body).Decode(&user)
	// convert random Int ID to a string
	user.ID = strconv.Itoa(rand.Intn(100000000)) // Mock ID - not safe
	users = append(users, user)
	// add user to server memory
	json.NewEncoder(write).Encode(user)
}

// Update user
func updateUser(write http.ResponseWriter, request *http.Request) {
	write.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for index, item := range users {
		if item.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			var user User
			_ = json.NewDecoder(request.Body).Decode(&user)
			user.ID = params["id"]
			users = append(users, user)
			json.NewEncoder(write).Encode(user)
			return
		}
	}
}

// Delete user
func deleteUser(write http.ResponseWriter, request *http.Request) {
	write.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for index, item := range users {
		if item.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			break
		}
	}
	json.NewEncoder(write).Encode(users)
}

func fileUpload(write http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(write, "uploading file\n")

	// 1. parse input , type multipart/form-data
	r.ParseMultipartForm(10 << 20)

	// 2. retrieve file from posted form-date
	file, handler, err := r.FormFile("myfile")
	if err != nil {
		log.Println(err)
		http.Error(write, "Error Retrieving file from form-data", http.StatusInternalServerError)
		return
	}

	defer file.Close()
	// print headers to console
	fmt.Printf("Uploading File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// 3. write temporary file on the server
	// create a temp file in our director
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)

	// 4. return whether or not this has been successful
	fmt.Fprintf(write, "Successfully uploaded file\n")
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	var s string
	if r.Method == http.MethodPost {
		file, _, err := r.FormFile("usrfile")
		if err != nil {
			log.Println(err)
			http.Error(w, "Error uploading file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		s = string(fileBytes)
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	fmt.Fprintf(w, `<form action="http://localhost:8000/upload1" method="post" enctype="multipart/form-data">     
	upload a file<br>
	<input type="file" name="usrfile" />
	<input type="submit" />
	</form>
	<br>
	<br>
	<h1>%v</h1>`, s)

}
