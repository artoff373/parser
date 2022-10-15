package main

import (
	"log"
	"mysearch/myparser"
	"net/http"
)

func main() {
	err := myparser.Search()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/report", report)
	http.HandleFunc("/", root)
	http.HandleFunc("/settings", settings)
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
