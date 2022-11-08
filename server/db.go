package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"mysearch/myparser"
	"net/http"
)

// Отображение страницы с настройками подключения к БД
func db(w http.ResponseWriter, r *http.Request) {
	var result template.HTML
	PDD.Db.Close()
	PDD.Connect = false
	var start = template.Must(template.ParseFiles("./templates/db.html"))
	if err := start.Execute(w, result); err != nil {
		log.Fatal(err)
	}
}

// Проверка настроек подключения к БД
func checkJSON(w http.ResponseWriter, r *http.Request) {
	var setDb myparser.DbData
	set, _ := io.ReadAll(r.Body)
	err := setDb.JsonParse(set)
	if err != nil {
		fmt.Fprintf(w, "проблемы с парсингом настроек базе данных\n%v", err)
		return
	}
	err = setDb.Сonnecting()
	if err != nil {
		fmt.Fprintf(w, "проблемы с подключением к базе данных\n%v", err)
		return
	}
	fmt.Fprint(w, "Успешное подключение\n")
	rows, err := setDb.Db.Query(`select table_name from information_schema."tables" where table_schema='Search'`)
	if err != nil {
		fmt.Fprintf(w, "проблемы с получением списка таблиц\n%v", err)
		return
	}
	defer rows.Close()
	fmt.Fprint(w, "Список таблиц:")
	for rows.Next() {
		var table string
		rows.Scan(&table)
		fmt.Fprintf(w, "\n%s", table)
	}
}
