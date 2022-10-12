package myparser

import (
	"encoding/xml"
)

type Post struct {
	Title      string   `bson:"title"`
	PubDate    string   `bson:"pub_date"`
	Link       string   `bson:"link"`
	Text       string   `bson:"text"`
	Dictionary []string `bson:"dictionary"`
	Relev      float64  `bson:"relev"`
	IsInReport bool     `bson:"is_in_report"`
	Fresh      bool     `bson:"fresh"`
	SearchDate string   `bson:"search_date"`
	Processed  bool     `bson:"processed"`
}

type Profile struct {
	Name   string   `bson:"name"`
	Keys   []string `bson:"keys"`
	Last   string   `bson:"last"`
	Source Sources  `bson:"sources"`
}
type Sources []string

type Rss struct {
	XMLName  xml.Name `xml:"rss"`
	Channels Channel  `xml:"channel"`
	Selector string
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
