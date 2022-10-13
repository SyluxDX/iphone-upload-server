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
	Lines []string
}

// Configs Global variable containing Configurations
var (
	Configs utils.Configurations
)

func uploadFiles(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	files := r.MultipartForm.File["myFiles"]

	uploadedFiles := []string{}
	for _, fileHeader := range files {
		// open file
		buffFile, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer buffFile.Close()

		uploadFile, err := os.Create(filepath.Join("upload", fileHeader.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer uploadFile.Close()

		_, err = io.Copy(uploadFile, buffFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		uploadedFiles = append(uploadedFiles, fileHeader.Filename)
	}
	// present confirmation page
	confirm := confimation{Lines: uploadedFiles}
	template, err := template.ParseFiles("html/confirm.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, confirm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
	http.HandleFunc("/confirm", uploadFiles)

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
