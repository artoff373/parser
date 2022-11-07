package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

func db(w http.ResponseWriter, r *http.Request) {
	var result template.HTML
	PDD.Db.Close()
	PDD.Connect = false
	var start = template.Must(template.ParseFiles("./templates/db.html"))
	if err := start.Execute(w, result); err != nil {
		log.Fatal(err)
	}
}
func check(w http.ResponseWriter, r *http.Request) {
	psqlconn, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Fprintf(w, "проблемы с передачей настроек сервера\n%v", err)
		return
	}
	db, _ := sql.Open("postgres", string(psqlconn))
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Fprintf(w, "проблемы с подключением\n%v", err)
		return
	}
	fmt.Fprint(w, "Успешное подключение\n")
	rows, err := db.Query(`select table_name from information_schema."tables" where table_schema='Search'`)
	if err != nil {
		fmt.Fprintf(w, "проблемы с получением списка таблиц\n%v", err)
		return
	}
	defer rows.Close()
	fmt.Fprint(w, "Список таблиц:")
	//var found = false
	for rows.Next() {
		var table string
		//found = true
		rows.Scan(&table)
		fmt.Fprintf(w, "\n%s", table)
	}
	// if !found {
	// 	fmt.Fprint(w, " нет\nПодключение выполнено, но таблиц не найдено.\nНачинаем инициализацию схемы БД.")
	// 	_, err = db.Exec(`CREATE SCHEMA IF NOT EXISTS "Sarch" AUTHORIZATION postgres;`)
	// 	if err != nil {
	// 		fmt.Fprintf(w, "проблемы с созданием схемы\n%v", err)
	// 		return
	// 	}
	// 	fmt.Fprint(w, "Схема создана")
	// 	_, err = db.Exec(`
	// 	CREATE SEQUENCE IF NOT EXISTS "Sarch"."Profile_id_seq"
	//   INCREMENT 1
	//   START 1
	//   MINVALUE 1
	//   MAXVALUE 2147483647
	//   CACHE 1
	//   OWNED BY "Profile".id;

	// 	ALTER SEQUENCE "Sarch"."Profile_id_seq"
	//   OWNER TO postgres;
	// 	CREATE TABLE IF NOT EXISTS "Sarch"."Profile"
	// 	(
	//   id integer NOT NULL DEFAULT nextval('"Sarch"."Profile_id_seq"'::regclass),
	//   name text COLLATE pg_catalog."default",
	//   keys text[] COLLATE pg_catalog."default",
	//   last_search text COLLATE pg_catalog."default",
	//   CONSTRAINT "Profile_pkey" PRIMARY KEY (id)
	// 		)

	// 	TABLESPACE pg_default;

	// 		ALTER TABLE IF EXISTS "Sarch"."Profile"
	//   	OWNER to postgres;
	// 	`)
	// 	if err != nil {
	// 		fmt.Fprintf(w, "проблемы с созданием таблицы профилей\n%v", err)
	// 		return
	// 	}
	// 	fmt.Fprint(w, "Таблица профилей создана")
	// }

}
