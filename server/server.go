package main

import (
	"log"
	"mysearch/myparser"
	"net/http"
)

func main() {

	http.HandleFunc("/report", report)
	http.HandleFunc("/", root)
	http.HandleFunc("/settings", settings)
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	log.Fatal(http.ListenAndServe(":8000", nil))
	go myparser.Search()
}
