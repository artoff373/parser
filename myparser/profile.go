package myparser

import (
	"database/sql"
	"encoding/xml"
	"fmt"

	"github.com/lib/pq"
)

type Post struct {
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

func (post *Post) insert(pID int, sID int, db *sql.DB) error {
	res, err := db.Query(fmt.Sprintf(`INSERT INTO "Search"."Posts" ("title") VALUES ('%s') RETURNING id`, post.Title))
	if err != nil {
		return fmt.Errorf("ошибка вставки - %v", err)
	}
	defer res.Close()
	var id int
	for res.Next() {
		err = res.Scan(&id)
		if err != nil {
			return fmt.Errorf("ошибка получения id вставки - %v", err)
		}
	}
	_, err = db.Exec(fmt.Sprintf(`UPDATE "Search"."Posts" SET is_in_report=false, fresh=true,  text='%s', relev=%f, url='%s' WHERE "id" = %d`, post.Text, post.Relev, post.Link, id))
	if err != nil {
		return fmt.Errorf("ошибка вставки - %v", err)
	}
	_, err = db.Exec(fmt.Sprintf(`UPDATE "Search"."Posts" SET profile_id=%d, source_id=%d, pub_date='%s', search_date='%s' WHERE "id" = %d`, pID, sID, post.PubDate, post.SearchDate, id))
	if err != nil {
		return fmt.Errorf("ошибка вставки - %v", err)
	}
	_, err = db.Exec(fmt.Sprintf(`UPDATE "Search"."Posts" SET dictionary=%vWHERE "id" = %d`, post.Dictionary, id))
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
