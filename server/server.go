package main

import (
	"encoding/json"
	"io"
	"log"
	"mysearch/myparser"
	"net/http"
	"os"
)

var PDD myparser.DbData

func init() {
	f, err := os.Open("config.ini")
	if err != nil {
		log.Printf("проблемы с открытием файла настроек %v", err)
		PDD.NewDb("localhost", "Search", "postgres", "1q2w3e4r", "5432")
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		log.Printf("проблемы с чтением файла настроек %v", err)
		PDD.NewDb("localhost", "Search", "postgres", "1q2w3e4r", "5432")
	}
	err = json.Unmarshal(data, &PDD)
	if err != nil {
		log.Printf("проблемы с разбором файла настроек %v", err)
		PDD.NewDb("localhost", "Search", "postgres", "1q2w3e4r", "5432")
	}
	PDD.Сonnecting()
}

func main() {
	go myparser.Search()
	http.HandleFunc("/index", root)
	http.HandleFunc("/report", report)
	http.HandleFunc("/db", db)
	http.HandleFunc("/check_json", checkJSON)
	http.HandleFunc("/save_json", saveJSON)
	http.HandleFunc("/chksrc", checkSource)
	http.HandleFunc("/make", make)
	http.HandleFunc("/archive", archive)
	http.HandleFunc("/settings", settings)
	http.HandleFunc("/print", print)
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
