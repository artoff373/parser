package main

import (
	"html/template"
	"log"
	"net/http"
)

func db(w http.ResponseWriter, r *http.Request) {
	var result template.HTML
	PDD.Db.Close()
	PDD.Connect = false
	var start = template.Must(template.ParseFiles("./templates/db.html"))
	if err := start.Execute(w, result); err != nil {
		log.Fatal(err)
	}
}
