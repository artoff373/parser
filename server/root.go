package main

import (
	"fmt"
	"html/template"
	"log"
	"mysearch/myparser"
	"net/http"
)

func root(w http.ResponseWriter, r *http.Request) {
	var root rootHTML
	db, err := myparser.Сonnecting()
	if err != nil {
		fmt.Printf("проблемы с подключением к базе%v", err)
	}
	defer db.Close()
	rows, err := db.Query(`SELECT "name", "id" FROM "Search"."Profile"`)
	if err != nil {
		fmt.Printf("проблемы с получением списка профилей\n%v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var profile roots
		var name string
		var id int
		err = rows.Scan(&name, &id)
		if err != nil {
			fmt.Fprintf(w, "проблемы с получением списка профилей\n%v", err)
		}
		profile.Name = template.HTML(name)
		profile.ID = template.HTML(fmt.Sprint(id))
		root = append(root, profile)
	}
	var start = template.Must(template.ParseFiles("./templates/index.html"))
	if err := start.Execute(w, root); err != nil {
		log.Fatal(err)
	}
}
