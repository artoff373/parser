package myparser

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type DbData struct {
	Host     string `json:"host"`     //= "localhost"
	Port     int    `json:"port"`     //= 5432
	User     string `json:"user"`     //= "postgres"
	Password string `json:"password"` //= "1q2w3e4r"
	DbName   string `json:"dbName"`   //= "Search"
	Db       *sql.DB
	Connect  bool
}

// Инициализация БД
func (DD *DbData) NewDb(host, dbname, user, password string, port int) {
	DD.Host = host
	DD.Port = port
	DD.DbName = dbname
	DD.Password = password
	DD.User = user
}

// Подключение к БД
func (DD *DbData) Сonnecting() error {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DD.Host, DD.Port, DD.User, DD.Password, DD.DbName)
	db, _ := sql.Open("postgres", psqlconn)
	err := db.Ping()
	if err != nil {
		return fmt.Errorf("проблемы с подключением\n%v", err)
	}
	DD.Db = db
	DD.Connect = true
	return nil
}

type Post struct {
	ID         int
	Title      string
	PubDate    string
	Link       string
	Text       string
	Dictionary pq.StringArray
	Relev      float64
	IsInReport bool
	Fresh      bool
	SearchDate string
}

// Вставка поста в БД
func (post *Post) insert(pID int, sID int, db *sql.DB) error {
	insert := fmt.Sprintf(`INSERT INTO "Search"."Posts" ("title", "text", "pub_date", "relev", "url", "is_in_report", "fresh", "profile_id", "source_id", "search_date", "dictionary") VALUES ('%s', '%s', '%s', %f, '%s', false, true, %d, %d, '%s', '{"%s"}')`, post.Title, post.Text, post.PubDate, post.Relev, post.Link, pID, sID, post.SearchDate, strings.Join(post.Dictionary, "\", \""))
	_, err := db.Exec(insert)
	if err != nil {
		return fmt.Errorf("ошибка вставки - %v", err)
	}
	return nil
}

type Profile struct {
	ID         int
	Name       string
	LastSearch string
	Keys       pq.StringArray
	Source     Sources
}
type Sources struct {
	ID       int
	Name     string
	URL      string
	Selector string
}

// Обновление времени последнего поиска в профиля
func (p *Profile) Update(newSearch string, db *sql.DB) error {
	_, err := db.Exec(`UPDATE "Search"."Profile" SET "last_search" = '%s' WHERE "id"=%d`, newSearch, p.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления - %v", err)
	}
	return nil
}

// Добавление ключа в профиль
func (p *Profile) AddKey(newKey string, db *sql.DB) error {
	add := fmt.Sprintf(`UPDATE "Search"."Profile" SET "keys" = array_append("keys", '%s') WHERE "id" = %d`, newKey, p.ID)
	_, err := db.Exec(add)
	if err != nil {
		return fmt.Errorf("ошибка добавления ключа - %v", err)
	}
	return nil
}

// Добавление источника в профиль
func (p *Profile) AddSource(s Sources, db *sql.DB) error {
	insert := fmt.Sprintf(`INSERT INTO "Search"."Sources" ("name", "url", "selector", "profile_id") VALUES ('%s', '%s', '%s', %d)`, s.Name, s.URL, s.Selector, p.ID)
	_, err := db.Exec(insert)
	if err != nil {
		return fmt.Errorf("ошибка добавления источника - %v\n %s", err, insert)
	}
	return nil
}

type Rss struct {
	XMLName  xml.Name `xml:"rss"`
	Channels Channel  `xml:"channel"`
}
type Channel struct {
	XMLName xml.Name `xml:"channel"`
	News    []News   `xml:"item"`
}
type News struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	PubDate string   `xml:"pubDate"`
	Link    string   `xml:"link"`
	Content string   `xml:"description"`
}
