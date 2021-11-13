package main

// ref: https://github.com/bradtraversy/go_restapi

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

const (
	port string = ":8000"
)

func setupRoutes() {
	// Init router
	// router := handler.NewRouter()

	// router.Get("/", handler.Handler{H: handle})
	log.Println("Sever started on port 8000")
	router := mux.NewRouter()
	// Route handles & endpoints
	router.HandleFunc("/api/v1/users", getUsers).Methods("GET")
	router.HandleFunc("/api/v1/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/api/v1/users", createUser).Methods("POST")
	router.HandleFunc("/api/v1/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/api/v1/users/{id}", deleteUser).Methods("DELETE")
	router.HandleFunc("/upload", fileUpload)
	// router.HandleFunc("/upload1", uploadFile)
	// router.HandleFunc("/index", index()
	// req.HandleFunc("/", index).Methods("GET")

	// Start server
	log.Println("Server listening on port", port)
	log.Fatal(http.ListenAndServe(port, router)) // set listen port
}

// Main function
func main() {
	setupRoutes()

	// continue with psql tutorial: https://youtu.be/zj12MYTrkdc
}

func fileUpload(resp http.ResponseWriter, r *http.Request) {
	log.Printf("%s, uploading file\n", resp)

	// 1. parse input , type multipart/form-data
	r.ParseMultipartForm(10 << 20) // set constraints on file upload size

	// 2. retrieve data from file posted form-date
	file, fileHeader, err := r.FormFile("myfile")
	if err != nil {
		log.Println(err)
		http.Error(resp, "Error Retrieving file from form-data", http.StatusInternalServerError)
		return
	}

	defer file.Close()
	// print headers to console
	log.Printf("Uploading File: %+v\n", fileHeader.Filename)
	log.Printf("File Size: %+v\n", fileHeader.Size)
	log.Printf("MIME Header: %+v\n", fileHeader.Header)

	// 3. resp temporary file on the server
	// create a temp file in our director
	// tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	contentType := fileHeader.Header["Content-Type"][0]
	log.Println("Content Type: ", contentType)

	var osFile *os.File
	if contentType == "image/jpeg" {
		osFile, err = ioutil.TempFile("images", "*.jpg")
	} else if contentType == "application/pdf" {
		osFile, err = ioutil.TempFile("SOPs", "*.pdf")
	}
	log.Println("error:", err)
	defer osFile.Close()

	// func ReadAll(r io.Reader) ([]byte, error)
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}

	osFile.Write(fileBytes)

	resp.Header().Set("Content-Type", contentType)
	log.Println(resp, `<form action="http://localhost:8000/upload" method="post" enctype="multipart/form-data">     
	upload a file<br>
	<input type="file" name="usrfile" />
	<input type="submit" />
	</form>
	<br>
	<br>
	<h1>%v</h1>`, osFile)

	// 4. return whether or not this has been successful
	log.Printf("%s: Successfully uploaded file\n", resp)
}

// func uploadFile(w http.ResponseWriter, r *http.Request) {
// 	var s string
// 	if r.Method == http.MethodPost {
// 		file, _, err := r.FormFile("usrfile")
// 		if err != nil {
// 			log.Println(err)
// 			http.Error(w, "Error uploading file", http.StatusInternalServerError)
// 			return
// 		}
// 		defer file.Close()

// 		fileBytes, err := ioutil.ReadAll(file)
// 		if err != nil {
// 			log.Println(err)
// 			http.Error(w, "Error reading file", http.StatusInternalServerError)
// 			return
// 		}
// 		s = string(fileBytes)
// 	}

// 	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
// 	log.Println(w, `<form action="http://localhost:8000/upload1" method="post" enctype="multipart/form-data">
// 	upload a file<br>
// 	<input type="file" name="usrfile" />
// 	<input type="submit" />
// 	</form>
// 	<br>
// 	<br>
// 	<h1>%v</h1>`, s)

// }

func index(w http.ResponseWriter, r *http.Request) error {

	w.WriteHeader(http.StatusOK)
	base := filepath.Join("templates", "base.html")
	index := filepath.Join("templates", "index.html")

	templ, _ := template.ParseFiles(base, index)
	templ.ExecuteTemplate(w, "base", nil)
	return nil
}
