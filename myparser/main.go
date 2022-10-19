package myparser

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1q2w3e4r"
	dbname   = "Search"
)

// Подключение к БД
func Сonnecting() (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, fmt.Errorf("Проблемы с подключением\n%v", err)
	}
	return db, nil
}

func loggger(f *os.File, lg *chan string) {
	log.SetOutput(f)
	for {
		rec := <-*lg
		log.Println(rec)
	}
}

// Поиск по профилю
func Search() error {
	db, _ := Сonnecting()
	defer db.Close()
	f, err := os.OpenFile("errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	lg := make(chan string)
	go loggger(f, &lg)
	for {
		//Перебор профилей
		rows, err := db.Query(`SELECT "id", "name", "last_search", "keys" FROM "Search"."Profile"`)
		if err != nil {
			lg <- fmt.Sprintf("проблемы с получением списка профилей %v", err)
			return fmt.Errorf("проблемы с получением списка профилей %v", err)
		}
		defer rows.Close()
		for rows.Next() {
			var profile Profile
			err = rows.Scan(&profile.ID, &profile.Name, &profile.LastSearch, &profile.Keys)
			if err != nil {
				lg <- fmt.Sprintf("проблемы с присвоением списка профилей %v", err)
				return fmt.Errorf("проблемы с присвоением списка профилей\n%v", err)
			}
			//Парсим профили
			t, _ := parsingProfile(profile, db, &lg)
			//Обновляем время последнего поиска
			err := profile.Update(t, db)
			if err != nil {
				lg <- fmt.Sprint(err)
			}
			lg <- fmt.Sprintf("Закончил разбор %s", profile.Name)
		}
		time.Sleep(time.Minute * 10)
	}
}
