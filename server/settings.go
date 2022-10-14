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

func settings(w http.ResponseWriter, r *http.Request) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ := mongo.Connect(context.TODO(), clientOptions)
	collection := client.Database("parser").Collection("profile")
	profileName := r.PostFormValue("sets")
	newKey := r.PostFormValue("addKey")
	newSource := r.PostFormValue("addSource")
	newSelector := r.PostFormValue("addSelector")
	if newKey != "" {
		_, err := collection.UpdateOne(context.TODO(), bson.D{{Key: "name", Value: profileName}},
			bson.D{{Key: "$push", Value: bson.D{{Key: "keys", Value: newKey}}}})
		if err != nil {
			fmt.Fprintf(w, "ошибка добавления ключа %s - %v", newKey, err)
		}
	}
	if (newSource != "") && (newSelector != "") {
		push := newSource + "!" + newSelector
		_, err := collection.UpdateOne(context.TODO(), bson.D{{Key: "name", Value: profileName}},
			bson.D{{Key: "$push", Value: bson.D{{Key: "keys", Value: push}}}})
		if err != nil {
			fmt.Fprintf(w, "ошибка добавления источника %s - %v", newSource, err)
		}
	}
	result := collection.FindOne(context.TODO(), bson.D{{Key: "name", Value: profileName}})
	var profile = new(myparser.Profile)
	err := result.Decode(&profile)
	if err != nil {
		fmt.Fprintf(w, "ошибка получения настроек профиля %s - %v", profileName, err)
	}
	var htmlProfile profileHTML
	htmlProfile.Name = template.HTML(profile.Name)
	htmlProfile.Last = template.HTML("заглушка") //profile.LastSearch)
	for i := 0; i < len(profile.Keys); i++ {
		htmlProfile.Keys = append(htmlProfile.Keys, template.HTML(profile.Keys[i]))
		htmlProfile.Sources = append(htmlProfile.Sources, template.HTML("заглушка"))
	}
	/*for i := 0; i < len(profile.Source); i++ {
		htmlProfile.Sources = append(htmlProfile.Sources, template.HTML(profile.Source[i]))
	}*/

	var page = template.Must(template.ParseFiles("./templates/settings.html"))
	if err := page.Execute(w, htmlProfile); err != nil {
		log.Fatal(err)
	}
}
