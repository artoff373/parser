package myparser

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"strings"

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
