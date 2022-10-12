package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type templeHTML struct {
	Profiles []profileHTML
}
type profileHTML struct {
	Name    template.HTML
	Keys    []template.HTML
	Last    template.HTML
	Sources []template.HTML
}
type profileBD struct {
	Name string `bson:"name"`
}

func root(w http.ResponseWriter, r *http.Request) {
	var teHTML templeHTML
	var proBD []profileBD
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ := mongo.Connect(context.TODO(), clientOptions)
	collection := client.Database("parser").Collection("profile")
	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		fmt.Fprintf(w, "Проблемы с получением коллекции профилей%v", err)
	}
	err = cur.All(context.TODO(), &proBD)
	if err != nil {
		fmt.Fprintf(w, "Проблемы с присвоением коллекции профилей%v", err)
	}
	for i := 0; i < len(proBD); i++ {
		var name profileHTML
		name.Name = template.HTML(proBD[i].Name)
		teHTML.Profiles = append(teHTML.Profiles, name)
	}
	var start = template.Must(template.ParseFiles("./templates/index.html"))
	if err := start.Execute(w, teHTML); err != nil {
		log.Fatal(err)
	}
}
