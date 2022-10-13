package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"iphone-upload-server/utils"
)

type confimation struct {
	Name  string
	Lines []string
}

// Configs Global variable containing Configurations
var (
	Configs utils.Configurations
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("myFile")
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "File Not found")
		return
	}
	defer file.Close()
	// Create a temporary file
	tempFile, err := os.CreateTemp("temp", "*.upload")
	if err != nil {
		log.Println(err)
	}
	defer tempFile.Close()
	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	// request Confirmation page
	confirm := confimation{tempFile.Name(), []string{}}
	// check uploaded file
	confirm.Lines = utils.ParseUpload(fileBytes)

	template, err := template.ParseFiles("html/confirm.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, confirm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func processFile(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("filename")
	filename := filepath.Base(path)
	// move file to upload folder
	os.Rename(path, filepath.Join("upload", filename))
	// Serve final page
	http.ServeFile(w, r, "html/thankyou.html")
}

func getSelfIP(configIP string) string {
	if configIP != "0.0.0.0" {
		return configIP
	}
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		return configIP
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/upload.html")
	})
	http.HandleFunc("/confirm", uploadFile)
	http.HandleFunc("/thankyou", processFile)

	log.Printf("Serving on %[1]s port %[2]d (http://%[1]s:%[2]d/)\n", getSelfIP(Configs.ServerURL), Configs.ServerPort)
	http.ListenAndServe(fmt.Sprintf("%s:%d", Configs.ServerURL, Configs.ServerPort), nil)
}

func main() {
	log.Println("Preping server")
	var err error
	Configs, err = utils.GetConfigs()
	if err != nil {
		log.Fatalln(err)
	}
	// Create folders if not exist
	os.Mkdir(Configs.UploadFolder, os.ModePerm)
	// Setup Server
	setupRoutes()
}
