package main

import (
	"fmt"
	"html/template"
	"log"
	"mysearch/myparser"
	"net/http"
)

func print(w http.ResponseWriter, r *http.Request) {
	var posts myparser.Post
	var reports postHTML
	var source string
	report := r.FormValue("data")
	PDD.Сonnecting()
	defer PDD.Db.Close()
	qry := fmt.Sprintf(`SELECT "Posts"."id", "Posts"."title", 
	"Posts"."text", "Posts"."url", "Posts"."pub_date", 
	"Sources"."name" FROM "Search"."Posts" INNER JOIN "Search"."Sources" 
	ON "Posts".source_id = "Sources".id WHERE "Posts"."fresh" = false 
	AND "Posts"."search_date" = '%s' ORDER BY id ASC `, report)
	rows, err := PDD.Db.Query(qry)
	if err != nil {
		fmt.Fprintf(w, "проблемы с получением списка постов\n%v", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&posts.ID, &posts.Title, &posts.Text, &posts.Link, &posts.PubDate, &source)
		if err != nil {
			fmt.Fprintf(w, "проблемы с получением постов\n%v", err)
		}
		var item postiks
		item.Id = template.HTML(fmt.Sprint(posts.ID))
		item.Link = template.HTML(posts.Link)
		item.Text = template.HTML(posts.Text)
		item.Title = template.HTML(posts.Title)
		item.PubDate = template.HTML(posts.PubDate)
		item.Source = template.HTML(source)
		reports.Postik = append(reports.Postik, item)
	}
	var page = template.Must(template.ParseFiles("./templates/print.html"))
	if err := page.Execute(w, reports); err != nil {
		log.Fatal(err)
	}
}
