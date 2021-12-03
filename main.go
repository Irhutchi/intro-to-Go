package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

const (
	port       string = ":8000"
	apiVersion string = "/api/v1/"
)

func setupRouter() *mux.Router {
	// Initialise new router
	router := mux.NewRouter() //net/http lib provides a defualt multiplexer
	glog.Infof("Sever started on port %s", port)

	// Route handles & endpoints
	router.HandleFunc(apiVersion+"users", getUsers).Methods("GET")
	router.HandleFunc(apiVersion+"users/{id}", getUser).Methods("GET")
	router.HandleFunc(apiVersion+"users", createUser).Methods("POST")
	router.HandleFunc(apiVersion+"users/{id}", updateUser).Methods("PUT")
	router.HandleFunc(apiVersion+"users/{id}", deleteUser).Methods("DELETE")
	router.HandleFunc("/upload", fileUpload)
	// router.HandleFunc("/index", index).Methods("GET")
	// router.HandleFunc("/upload1", uploadFile)
	files := http.Dir("./public/")
	fileHandler := http.StripPrefix("/public/", http.FileServer(files))
	router.PathPrefix("/public/").Handler(fileHandler).Methods("GET")
	router.Handle("/", router)
	// router.HandleFunc("/hello", handler).Methods("GET")

	return router

	// glog.Fatal(http.ListenAndServe(port, router)) // set listen port
}

func main() {

	// This is needed to make `glog` believe that the flags have already been parsed, otherwise
	// every log messages is prefixed by an error message stating the the flags haven't been
	// parsed.
	_ = flag.CommandLine.Parse([]string{})

	// Always log to stderr by default
	if err := flag.Set("logtostderr", "true"); err != nil {
		glog.Infof("Unable to set logtostderr to true")
	}

	router := setupRouter()
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// enforce server timeout
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		glog.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func fileUpload(resp http.ResponseWriter, r *http.Request) {
	glog.Infof("%s, uploading file\n", resp)

	// 1. parse input , type multipart/form-data
	r.ParseMultipartForm(10 << 20) // set constraints on file upload size

	// 2. retrieve data from file posted form-date
	file, fileHeader, err := r.FormFile("myfile")
	if err != nil {
		glog.Warning(err)
		http.Error(resp, "Error Retrieving file from form-data", http.StatusInternalServerError)
		return
	}

	defer file.Close()
	// print headers to console
	glog.Infof("Uploading File: %+v\n", fileHeader.Filename)
	glog.Infof("File Size: %+v\n", fileHeader.Size)
	glog.Infof("MIME Header: %+v\n", fileHeader.Header)

	// 3. resp temporary file on the server
	// create a temp file in our director
	// tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	contentType := fileHeader.Header["Content-Type"][0]
	glog.Info("Content Type: ", contentType)

	var osFile *os.File
	if contentType == "image/jpeg" {
		osFile, err = ioutil.TempFile("images", "*.jpg")
	} else if contentType == "application/pdf" {
		osFile, err = ioutil.TempFile("SOPs", "*.pdf")
	}
	glog.Error("error:", err)
	defer osFile.Close()

	// func ReadAll(r io.Reader) ([]byte, error)
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		glog.Error(err)
	}

	//Stream to response
	// if _, err := io.Copy(w, f); err != nil {
	// 	glog.Info(err)
	// 	w.WriteHeader(500)
	// } // https://stackoverflow.com/questions/45302083/serving-up-pdfs-using-golang

	// osFile.Write(fileBytes)
	glog.Info(fileBytes)

	resp.Header().Set("Content-Type", contentType)
	glog.Info(resp, `<form action="http://localhost:8000/upload" method="post" enctype="multipart/form-data">     
	upload a file<br>
	<input type="file" name="usrfile" />
	<input type="submit" />
	</form>
	<br>
	<br>
	<h1>%v</h1>`, osFile)

	// 4. return whether or not this has been successful
	glog.Infof("%s: Successfully uploaded file\n", resp)

	// after file upload is done redirect to new page.
	// display pdf in iframe or similar
	// modify upload post show pdf in the front-end
	// what pkg will view pdf on front end.
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

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
	// w.WriteHeader(http.StatusOK)
	// base := filepath.Join("templates", "base.html")
	// index := filepath.Join("templates", "index.html")

	// templ, _ := template.ParseFiles(base, index)
	// templ.ExecuteTemplate(w, "base", nil)
	// return nil
}
