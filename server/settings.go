package main

import (
	"fmt"
	"html/template"
	"log"
	"mysearch/myparser"
	"net/http"
	"strconv"
)

func settings(w http.ResponseWriter, r *http.Request) {
	//собираем данные форм
	var profile = new(myparser.Profile)
	profile.ID, _ = strconv.Atoi(r.FormValue("sets"))
	var newSource myparser.Sources
	newKey := r.PostFormValue("addKey")
	newSource.Name = r.PostFormValue("addSource")
	newSource.URL = r.PostFormValue("addLink")
	newSource.Selector = r.PostFormValue("addSelector")
	//Ищем профиль в БД
	db, err := myparser.Сonnecting()
	if err != nil {
		fmt.Fprintf(w, "ошибка подключения к базе %v", err)
	}
	defer db.Close()
	if newKey != "" {
		err = profile.AddKey(newKey, db)
		if err != nil {
			fmt.Fprintf(w, "ошибка добавления ключа %s - %v", newKey, err)
		}
	}
	if (newSource.URL != "") && (newSource.Selector != "") && (newSource.Name != "") {
		err = profile.AddSource(newSource, db)
		if err != nil {
			fmt.Fprintf(w, "ошибка добавления источника %s - %v", newSource.Name, err)
		}
	}
	qry := fmt.Sprintf(`SELECT "name", "last_search", "keys" FROM "Search"."Profile" WHERE "id" = %d`, profile.ID)
	rows, err := db.Query(qry)
	if err != nil {
		fmt.Printf("проблемы с получением списка профилей\n%v", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&profile.Name, &profile.LastSearch, &profile.Keys)
		if err != nil {
			fmt.Fprintf(w, "проблемы с получением профиля\n%v", err)
		}
		var htmlProfile profileHTML
		htmlProfile.ID = template.HTML(fmt.Sprint(profile.ID))
		htmlProfile.Name = template.HTML(profile.Name)
		htmlProfile.Last = template.HTML(profile.LastSearch)
		for i := 0; i < len(profile.Keys); i++ {
			htmlProfile.Keys = append(htmlProfile.Keys, template.HTML(profile.Keys[i]))
		}
		query := fmt.Sprintf(`SELECT "id", "name","url", "selector" FROM "Search"."Sources" WHERE profile_id = %d`, profile.ID)
		sources, err := db.Query(query)
		if err != nil {
			fmt.Fprintf(w, "проблемы с получением списка источников -%v", err)
		}
		defer sources.Close()
		for sources.Next() {
			err = sources.Scan(&profile.Source.ID, &profile.Source.Name, &profile.Source.URL, &profile.Source.Selector)
			if err != nil {
				fmt.Fprintf(w, "проблемы с присвоением источника %s - %v", profile.Name, err)
			}
			var htmlSource sourceHTML
			htmlSource.ID = template.HTML(fmt.Sprint(profile.Source.ID))
			htmlSource.Name = template.HTML(profile.Source.Name)
			htmlSource.URL = template.HTML(profile.Source.URL)
			htmlSource.Selector = template.HTML(profile.Source.Selector)
			htmlProfile.Sources = append(htmlProfile.Sources, htmlSource)
		}

		var page = template.Must(template.ParseFiles("./templates/settings.html"))
		if err := page.Execute(w, htmlProfile); err != nil {
			log.Fatal(err)
		}
	}
}
