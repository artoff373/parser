package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
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
func check(w http.ResponseWriter, r *http.Request) {
	psqlconn, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Fprintf(w, "проблемы с передачей настроек сервера\n%v", err)
		return
	}
	db, _ := sql.Open("postgres", string(psqlconn))
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Fprintf(w, "проблемы с подключением\n%v", err)
		return
	}
	fmt.Fprint(w, "Успешное подключение")
}
