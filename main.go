package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"pdfsigning/utils"
	"time"
)

const (
	toBeSignedFileName = "tmp.pdf"
	signedFileName     = "signed.pdf"
	logoFileName = "lequest_logo.png"  // if you change this, change also the image file name as well
)

func SignPdfHandler(w http.ResponseWriter, r *http.Request) {
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
		_, err := w.Write([]byte(errMsg))
		if err != nil {
			log.Println(err)
		}
		return
	}

	// get working dir path
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// create tmp saving path
	savePath := filepath.Join(wd, "tmp", toBeSignedFileName)

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

	// rewind to the beginning of the uploaded file
	_, err = mf.Seek(0, 0)
	if err != nil {
		defer log.Println(removeFilesFromTmpDir())
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// we copy the data from the uploaded file to the new created file
	_, err = io.Copy(nf, mf)
	if err != nil {
		defer log.Println(removeFilesFromTmpDir())
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// we run the python script to generate the pdf with the signature
	pyScript := filepath.Join(wd, "python_scripts", "pdf_signing_process.py")
	src := savePath
	dest := filepath.Join(wd, "tmp", signedFileName)
	date := time.Now().Format("2006.01.02")
	logoPath := filepath.Join(wd, logoFileName)
	cmd := exec.Command("python", pyScript, src, dest, date, logoPath)
	err = cmd.Run()
	if err != nil {
		defer log.Println(removeFilesFromTmpDir())
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// if all goes well we send a response 200
	w.WriteHeader(http.StatusOK)
}

func SendSignedPdfHandler(w http.ResponseWriter, r *http.Request) {
	// get working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// create path where the signed filename is
	targetPath := filepath.Join(wd, "tmp", signedFileName)

	// run this func after returning to do some cleanup
	defer func() {
		err := removeFilesFromTmpDir()
		if err != nil {
			log.Println(err)
		}
	}()
	// serve file
	http.ServeFile(w, r, targetPath)
}

func removeFilesFromTmpDir() error {
	dir, err := ioutil.ReadDir("./tmp")
	if err != nil {
		log.Println(err)
		return err
	}
	for _, d := range dir {
		// we don't want to delete our ``.keepdir`` because of Git.
		if d.Name() == ".keepdir" {
			continue
		}
		err = os.RemoveAll(path.Join("tmp", d.Name()))
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", SignPdfHandler)
	mux.HandleFunc("/download-signed-pdf", SendSignedPdfHandler)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
