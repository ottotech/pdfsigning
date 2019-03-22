package utils

import (
	"fmt"
	"log"
	"net/http"
	"pdfsigning/config"
)

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := config.TPL.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, fmt.Sprintf("Template error: %s", err.Error()), http.StatusInternalServerError)
	}
}
