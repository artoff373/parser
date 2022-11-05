package main

import (
	"log"
	"mysearch/myparser"
	"net/http"
)

var PDD myparser.DbData

func init() {
	PDD.NewDb("localhost", "Search", "postgres", "1q2w3e4r", 5432)
	PDD.Сonnecting()
}

func main() {
	//go myparser.Search()
	http.HandleFunc("/index", root)
	http.HandleFunc("/report", report)
	http.HandleFunc("/db", db)
	http.HandleFunc("/settings", settings)
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
