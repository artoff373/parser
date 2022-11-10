package myparser

import (
	"fmt"
	"log"
	"os"
	"time"
)

func loggger(f *os.File) {
	defer f.Close()
	log.SetOutput(f)
	for {
		rec := <-lg
		log.Println(rec)
	}
}

var SDB DbData
var lg = make(chan string)

// Поиск по профилю
func Search() error {
	f, err := os.OpenFile("errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	go loggger(f)
	SDB.NewDb("localhost", "Search", "postgres", "1q2w3e4r", "5432")
	for {
		err = SDB.Сonnecting()
		if err != nil {
			log.Fatal(err)
		}
		//Перебор профилей
		rows, err := SDB.Db.Query(`SELECT "id", "name", "last_search", "keys" FROM "Search"."Profile"`)
		if err != nil {
			lg <- fmt.Sprintf("проблемы с получением списка профилей %v", err)
			return fmt.Errorf("проблемы с получением списка профилей %v", err)
		}

		for rows.Next() {
			var profile Profile
			err = rows.Scan(&profile.ID, &profile.Name, &profile.LastSearch, &profile.Keys)
			if err != nil {
				lg <- fmt.Sprintf("проблемы с присвоением списка профилей %v", err)
				return fmt.Errorf("проблемы с присвоением списка профилей\n%v", err)
			}
			//Парсим профили
			t := profile.parsingProfile()
			//Обновляем время последнего поиска
			err := profile.Update(t, SDB.Db)
			if err != nil {
				lg <- fmt.Sprint(err)
			}
			lg <- fmt.Sprintf("Закончил разбор %s", profile.Name)
		}
		rows.Close()
		SDB.Db.Close()
		time.Sleep(time.Minute * 10)
	}
}
