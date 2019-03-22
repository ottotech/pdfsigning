package main

import (
	"github.com/jung-kurt/gofpdf"
	"log"
	"net/http"
	"pdfsigning/utils"
)


func SignPdfHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet {
		utils.RenderTemplate(w, "index.gohtml", nil)
		return
	}
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

func createPDF()  {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
	err := pdf.OutputFileAndClose("hello.pdf")
	if err != nil {
		log.Println(err)
	}
}

