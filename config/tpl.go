package config

import "html/template"

var TPL *template.Template

func init() {
	TPL = template.Must(template.ParseFiles("index.gohtml"))
}
