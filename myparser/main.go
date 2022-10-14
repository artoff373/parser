package myparser

import (
	"database/sql"
	"fmt"
	"log"
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
func connecting() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Printf("Проблемы с подключением\n%v", err)
	}
	return db
}

// Поиск по профилю
func Search() error {
	for {
		db := connecting()
		defer db.Close()
		//Перебор профилей
		rows, err := db.Query(`SELECT "id", "name", "last_search", "keys" FROM "Search"."Profile"`)
		if err != nil {
			return fmt.Errorf("проблемы с получением коллекции профилей\n%v", err)
		}
		defer rows.Close()
		for rows.Next() {
			var profile Profile
			err = rows.Scan(&profile.ID, &profile.Name, &profile.LastSearch, &profile.Keys)
			if err != nil {
				return fmt.Errorf("проблемы с присвоением коллекции профилей\n%v", err)
			}
			//Парсим источники
			t, err := parsingProfile(profile, db)
			if err != nil {
				return fmt.Errorf("проблемы с парсингом профиля %s: %v", profile.Name, err)
			}
			//Обновляем время последнего поиска
			fmt.Printf("Закончил разбор %s - %v", profile.Name, t)
		}
		time.Sleep(time.Minute)
	}
}
