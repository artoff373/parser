package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"mysearch/myparser"
	"net/http"
	"strconv"
	"time"
)

// Страница с найденными новостями
func report(w http.ResponseWriter, r *http.Request) {
	var posts myparser.Post
	var reports postHTML
	var source string
	profile, _ := strconv.Atoi(r.FormValue("profile"))
	PDD.Сonnecting()
	defer PDD.Db.Close()
	qry := fmt.Sprintf(`SELECT "Posts"."id", "Posts"."title", 
	"Posts"."text", "Posts"."relev", "Posts"."url", "Posts"."pub_date", 
	"Sources"."name" FROM "Search"."Posts" INNER JOIN "Search"."Sources" 
	ON "Posts".source_id = "Sources".id WHERE "Posts"."fresh" = TRUE 
	AND "Posts"."profile_id" = %d ORDER BY id ASC `, profile)
	rows, err := PDD.Db.Query(qry)
	if err != nil {
		fmt.Printf("проблемы с получением списка постов\n%v", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&posts.ID, &posts.Title, &posts.Text, &posts.Relev, &posts.Link, &posts.PubDate, &source)
		if err != nil {
			fmt.Fprintf(w, "проблемы с получением постов\n%v", err)
		}
		var item postiks
		var relev = fmt.Sprintf("%.2f", posts.Relev)
		item.Id = template.HTML(fmt.Sprint(posts.ID))
		item.Link = template.HTML(posts.Link)
		item.Text = template.HTML(posts.Text)
		item.Title = template.HTML(posts.Title)
		item.PubDate = template.HTML(posts.PubDate)
		item.Relev = template.HTML(relev)
		item.Source = template.HTML(source)
		reports.Postik = append(reports.Postik, item)
	}
	var page = template.Must(template.ParseFiles("./templates/report.html"))
	if err := page.Execute(w, reports); err != nil {
		log.Fatal(err)
	}
}

// Обработка результата отбора новостей
func make(w http.ResponseWriter, r *http.Request) {
	PDD.Сonnecting()
	reportDate := time.Now().Format(time.RFC1123Z)
	defer PDD.Db.Close()
	get, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "проблемы с получением постов\n%v", err)
		return
	}
	qry := fmt.Sprintf(`UPDATE "Search"."Posts" SET "is_in_report" = true, "fresh" = false, "search_date" = '%s' WHERE "id" IN (%s)`, reportDate, string(get))
	_, err = PDD.Db.Exec(qry)
	if err != nil {
		fmt.Fprintf(w, "проблемы с обработкой списка постов\n%v", err)
		return
	}
	/*
		qry = fmt.Sprintf(`DELETE FROM "Search"."Posts"  WHERE "id"  NOT IN (%s)`, string(get))
		_, err = PDD.Db.Exec(qry)
		if err != nil {
			fmt.Fprintf(w, "проблемы с обработкой списка постов\n%v", err)
		}*/
	fmt.Fprint(w, "НЕТ проблемы с обработкой списка постов\n")
}
