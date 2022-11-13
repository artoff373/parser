package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// Страница с отчетами
func archive(w http.ResponseWriter, r *http.Request) {
	var reports []Archive
	PDD.Сonnecting()
	defer PDD.Db.Close()
	qry := `SELECT DISTINCT search_date FROM "Search"."Posts" WHERE fresh = false ORDER BY search_date DESC`
	rows, err := PDD.Db.Query(qry)
	if err != nil {
		fmt.Fprintf(w, "проблемы с получением списка постов\n%v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var rec string
		var item Archive
		err = rows.Scan(&rec)
		if err != nil {
			fmt.Fprintf(w, "проблемы с получением постов\n%v", err)
		}
		dt, err := time.Parse(time.RFC1123Z, rec)
		if err != nil {
			fmt.Fprintf(w, "проблемы с парсингом даты %s: %v", dt, err)
		}
		item.Val = template.HTML(rec)
		item.Read = template.HTML(dt.Format(time.RFC822))
		reports = append(reports, item)
	}
	var page = template.Must(template.ParseFiles("./templates/archive.html"))
	if err := page.Execute(w, reports); err != nil {
		log.Fatal(err)
	}
}
