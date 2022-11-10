package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"mysearch/myparser"
	"net/http"
	"os"
)

// Отображение страницы с настройками подключения к БД
func db(w http.ResponseWriter, r *http.Request) {
	var settings DbHtml
	if PDD.Connect {
		PDD.Db.Close()
		PDD.Connect = false
	}
	settings.Host = template.HTML(PDD.Host)
	settings.Port = template.HTML(PDD.Port)
	settings.User = template.HTML(PDD.User)
	settings.Password = template.HTML(PDD.Password)
	settings.DbName = template.HTML(PDD.DbName)
	var start = template.Must(template.ParseFiles("./templates/db.html"))
	if err := start.Execute(w, settings); err != nil {
		log.Fatal(err)
	}
}

// Сохранение настроек БД
func saveJSON(w http.ResponseWriter, r *http.Request) {
	var setDb myparser.DbData
	set, _ := io.ReadAll(r.Body)
	err := setDb.JsonParse(set)
	if err != nil {
		fmt.Fprintf(w, "проблемы с парсингом настроек базе данных\n%v", err)
		return
	}
	f, _ := os.OpenFile("config.ini", os.O_WRONLY, 0755)
	defer f.Close()
	data, err := json.Marshal(setDb)
	if err != nil {
		fmt.Fprintf(w, "проблемы с парсингом настроек базе данных\n%v", err)
	}
	f.WriteAt(data, 0)
	fmt.Fprintf(w, "настройки сохранены\n")
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
