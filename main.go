package main

// ref: https://github.com/bradtraversy/go_restapi

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

const (
	port string = ":8000"
)

// func setupRoutes() {
// 	// Init router
// 	// router := handler.NewRouter()

// 	// router.Get("/", handler.Handler{H: handle})
// 	fmt.Println("Sever started on port 8000")
// 	router := mux.NewRouter()
// 	// Route handles & endpoints
// 	router.HandleFunc("/users", getUsers).Methods("GET")
// 	router.HandleFunc("/users/{id}", getUser).Methods("GET")
// 	router.HandleFunc("/users", createUser).Methods("POST")
// 	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
// 	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
// 	router.HandleFunc("/upload", fileUpload)
// 	router.HandleFunc("/upload1", uploadFile)
// 	// router.HandleFunc("/index", index()
// 	// req.HandleFunc("/", index).Methods("GET")

// 	// Start server
// 	log.Println("Server listening on port", port)
// 	log.Fatal(http.ListenAndServe(port, router)) // set listen port
// }

// Main function
func main() {
	// setupRoutes()

	// continue with psql tutorial: https://youtu.be/zj12MYTrkdc

	fmt.Println("Sever started on port 8000")
	router := mux.NewRouter()
	// Route handles & endpoints
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	router.HandleFunc("/upload", fileUpload)
	router.HandleFunc("/upload1", uploadFile)
	// router.HandleFunc("/index", index()
	// req.HandleFunc("/", index).Methods("GET")

	// Start server
	log.Println("Server listening on port", port)
	log.Fatal(http.ListenAndServe(port, router)) // set listen port
}

func fileUpload(resp http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(resp, "uploading file\n")

	// 1. parse input , type multipart/form-data
	r.ParseMultipartForm(10 << 20)

	// 2. retrieve file from posted form-date
	file, handler, err := r.FormFile("myfile")
	if err != nil {
		log.Println(err)
		http.Error(resp, "Error Retrieving file from form-data", http.StatusInternalServerError)
		return
	}

	defer file.Close()
	// print headers to console
	fmt.Printf("Uploading File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// 3. resp temporary file on the server
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
	fmt.Fprintf(resp, "Successfully uploaded file\n")
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

func index(w http.ResponseWriter, r *http.Request) error {

	w.WriteHeader(http.StatusOK)
	base := filepath.Join("templates", "base.html")
	index := filepath.Join("templates", "index.html")

	templ, _ := template.ParseFiles(base, index)
	templ.ExecuteTemplate(w, "base", nil)
	return nil
}
