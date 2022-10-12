package myparser

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Подключение к БД
func connecting() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Printf("Проблемы с подключением\n%v", err)
	}
	return client
}

// Поиск по профилю
func Search() error {
	for {
		client := connecting()
		defer client.Disconnect(context.TODO())
		//Перебор профилей
		collection := client.Database("parser").Collection("profile")
		cur, err := collection.Find(context.TODO(), bson.D{{}})
		if err != nil {
			return fmt.Errorf("проблемы с получением коллекции профилей\n%v", err)
		}
		var profiles []Profile
		err = cur.All(context.TODO(), &profiles)
		if err != nil {
			return fmt.Errorf("проблемы с присвоением коллекции профилей\n%v", err)
		}
		for _, profile := range profiles {
			//Парсим источники
			t, err := parsingProfile(profile, client)
			if err != nil {
				return fmt.Errorf("проблемы с парсингом профиля %s: %v", profile.Name, err)
			}
			//Обновляем время последнего поиска
			filter := bson.D{{Key: "name", Value: profile.Name}}
			update := bson.D{{Key: "$set", Value: bson.D{{Key: "last", Value: t}}}}
			collection.FindOneAndUpdate(context.TODO(), filter, update)
		}
		time.Sleep(time.Minute)
	}
}
