package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func root(w http.ResponseWriter, r *http.Request) {
	var root rootHTML
	if !PDD.Connect {
		if r.PostFormValue("dbName") != "" {
			PDD.DbName = r.PostFormValue("dbName")
			PDD.Host = r.PostFormValue("host")
			PDD.Port = r.PostFormValue("port")
			PDD.User = r.PostFormValue("user")
			PDD.Password = r.PostFormValue("password")
		}
	}
	err := PDD.Сonnecting()
	if err != nil {
		http.Redirect(w, r, "http://127.0.0.1:8000/db", http.StatusSeeOther)
		return
	}
	name := r.PostFormValue("new")
	if name != "" {
		ls := time.Now().Add(-time.Hour * 24).Format(time.RFC1123Z)
		query := fmt.Sprintf(`INSERT INTO "Search"."Profile"( name, last_search ) VALUES ('%s', '%s')`, name, ls)
		_, err := PDD.Db.Exec(query)
		if err != nil {
			fmt.Fprintf(w, "проблемы с добавлением профиля\n%v", err)
			return
		}
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
