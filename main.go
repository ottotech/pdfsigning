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
)

const (
	toBeSignedFileName = "tmp.pdf"
	signedFileName     = "signed.pdf"
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

	// create paths
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
	// we run or great python script to generate the pdf with the signature
	pyScript := filepath.Join(wd, "python_scripts", "pdf_signing_process.py")
	src := savePath
	dest := filepath.Join(wd, "tmp", signedFileName)
	cmd := exec.Command("python", pyScript, src, dest)
	err = cmd.Run()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// if all goes well we send a response 200
	w.WriteHeader(http.StatusOK)
}

func SendSignedPdf(w http.ResponseWriter, r *http.Request) {
	// get working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// create path where the signed filename is
	targetPath := filepath.Join(wd, "tmp", signedFileName)
	defer func() {
		err := removeFilesFromTmpDir()
		if err != nil {
			log.Println(err)
		}
	}()
	http.ServeFile(w, r, targetPath)
}

func removeFilesFromTmpDir() error {
	dir, err := ioutil.ReadDir("./tmp")
	if err != nil {
		log.Println(err)
		return err
	}
	for _, d := range dir {
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
	mux.HandleFunc("/download-signed-pdf", SendSignedPdf)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
