package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"mysearch/myparser"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type postHTML struct {
	Postik []postiks
}
type postiks struct {
	Id      template.HTML
	Title   template.HTML
	PubDate template.HTML
	Text    template.HTML
	Link    template.HTML
}

func report(w http.ResponseWriter, r *http.Request) {
	var bd []myparser.Post
	profile := r.FormValue("profile")
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ := mongo.Connect(context.TODO(), clientOptions)
	collection := client.Database("parser").Collection(profile)
	result, err := collection.Find(context.TODO(), bson.D{{Key: "relev", Value: bson.D{{Key: "$gt", Value: 1}}}})
	if err != nil {
		fmt.Fprintf(w, "ошибка поиска - %v", err)
		return
	}
	err = result.All(context.TODO(), &bd)
	if err != nil {
		fmt.Fprintf(w, "ошибка поиска - %v", err)
		return
	}
	if len(bd) == 0 {
		fmt.Fprint(w, "Ничего не найдено")
		return
	}
	var reports postHTML
	for i := 0; i < len(bd); i++ {
		var item postiks
		item.Id = template.HTML(fmt.Sprint(i))
		item.Link = template.HTML(bd[i].Link)
		item.Text = template.HTML(bd[i].Text)
		item.Title = template.HTML(bd[i].Title)
		item.PubDate = template.HTML(bd[i].PubDate)
		reports.Postik = append(reports.Postik, item)
	}
	var page = template.Must(template.ParseFiles("./templates/report.html"))
	if err := page.Execute(w, reports); err != nil {
		log.Fatal(err)
	}
}
