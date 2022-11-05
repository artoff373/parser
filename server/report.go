package main

import (
	"fmt"
	"html/template"
	"log"
	"mysearch/myparser"
	"net/http"
	"strconv"
)

func report(w http.ResponseWriter, r *http.Request) {
	var posts myparser.Post
	var reports postHTML
	profile, _ := strconv.Atoi(r.FormValue("profile"))
	PDD.Сonnecting()
	defer PDD.Db.Close()
	qry := fmt.Sprintf(`SELECT "id", "title", "text", "relev", "url", "pub_date" FROM "Search"."Posts" WHERE "fresh" = TRUE AND "profile_id" = %d`, profile)
	rows, err := PDD.Db.Query(qry)
	if err != nil {
		fmt.Printf("проблемы с получением списка постов\n%v", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&posts.ID, &posts.Title, &posts.Text, &posts.Relev, &posts.Link, &posts.PubDate)
		if err != nil {
			fmt.Fprintf(w, "проблемы с получением постов\n%v", err)
		}
		var item postiks
		item.Id = template.HTML(fmt.Sprint(posts.ID))
		item.Link = template.HTML(posts.Link)
		item.Text = template.HTML(posts.Text)
		item.Title = template.HTML(posts.Title)
		item.PubDate = template.HTML(posts.PubDate)
		reports.Postik = append(reports.Postik, item)
	}
	var page = template.Must(template.ParseFiles("./templates/report.html"))
	if err := page.Execute(w, reports); err != nil {
		log.Fatal(err)
	}
}
