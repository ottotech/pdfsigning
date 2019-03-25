package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"pdfsigning/utils"
)


func SignPdfHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet {
		utils.RenderTemplate(w, "index.gohtml", nil)
		return
	}

	// get file from request
	mf, fh, fileErr := r.FormFile("nf")
	if fileErr != nil {
		if fileErr == http.ErrMissingFile {
			errMsg := "Error: You need to upload a pdf file!"
			w.WriteHeader(http.StatusForbidden)
			log.Println(w.Write([]byte(errMsg)))
			return
		}
		log.Println(fileErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		err := mf.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	// check content type of uploaded file
	ct := fh.Header.Get("Content-Type")

	// if ct is not pdf, we send an error
	if ct != "application/pdf" {
		errMsg := fmt.Sprintf("Error: You need to upload a pdf file! Not a file of type %v", ct)
		w.WriteHeader(http.StatusForbidden)
		log.Println(w.Write([]byte(errMsg)))
		return
	}

	// create filename of file to be signed
	filename := "file_to_be_signed.pdf"

	// create paths
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	savePath := filepath.Join(wd, "tmp", filename)

	// we create a new file to copy all the data from the uploaded file
	nf, err := os.Create(savePath)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		err := nf.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	// we rewind the uploaded file
	_, err = mf.Seek(0, 0)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// we copy the data from the uploaded file to the new created file
	_, err = io.Copy(nf, mf)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// if all goes well we send a response 200
	w.WriteHeader(http.StatusOK)

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", SignPdfHandler)
	server := http.Server{
		Addr:              ":8080",
		Handler:           mux,
	}
	log.Fatal(server.ListenAndServe())
}

