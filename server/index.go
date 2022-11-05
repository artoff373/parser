package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func root(w http.ResponseWriter, r *http.Request) {
	var root rootHTML
	if !PDD.Connect {
		if r.PostFormValue("DbName") != "" {
			PDD.DbName = r.PostFormValue("DbName")
			PDD.Host = r.PostFormValue("Host")
			PDD.Port, _ = strconv.Atoi(r.PostFormValue("Port"))
			PDD.User = r.PostFormValue("User")
			PDD.Password = r.PostFormValue("Password")
		}
	}
	err := PDD.Сonnecting()
	if err != nil {
		//fmt.Fprintf(w, "проблемы с подключением к базе - %v", err)
		http.Redirect(w, r, "http://127.0.0.1:8000/db", http.StatusSeeOther)
		return
	}
	rows, err := PDD.Db.Query(`SELECT "name", "id" FROM "Search"."Profile"`)
	if err != nil {
		fmt.Fprintf(w, "проблемы с получением списка профилей\n%v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var profile roots
		var name string
		var id int
		err = rows.Scan(&name, &id)
		if err != nil {
			fmt.Fprintf(w, "проблемы с получением списка профилей\n%v", err)
			return
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
