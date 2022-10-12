package myparser

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"mysearch/words"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/mongo"
)

// Метод для парсинга конкретной новости из xml
func (n *News) parsingNews(ls *time.Time, sel *string, node *words.Branch,
	ns *string, client *mongo.Client, p *string) error {
	t, err := time.Parse(time.RFC1123Z, n.PubDate)
	if err != nil {
		return fmt.Errorf("проблемы с парсингом даты %s: %v", n.PubDate, err)
	}

	if t.After(*ls) {
		var post Post
		post.Text, err = getNews(n.Link, *sel)
		if err != nil {
			return fmt.Errorf("проблемы с парсингом страницы: %v", err)
		}
		post.Dictionary = words.SortUniq(strings.Split(post.Text, " "))
		post.Relev = words.SearchKeys(node, post.Dictionary)
		if post.Relev > 1 {
			post.PubDate = n.PubDate
			post.Title = n.Title
			post.Link = n.Link
			post.SearchDate = *ns
			post.Fresh = true
			collection := client.Database("parser").Collection(*p)
			_, err := collection.InsertOne(context.TODO(), post)
			if err != nil {
				return fmt.Errorf("проблемы со вставкой документа в базу%v", err)
			}
		}
	}
	return nil
}

// Парсим текст новости из html страницы
func getNews(link, selector string) (string, error) {
	result, err := http.Get(link)
	if err != nil {
		return "", fmt.Errorf("проблемы с получением страницы новостей%s: %v", link, err)
	}
	defer result.Body.Close()
	if result.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error on the link: %s status - %s", link, result.Status)
	}
	page, err := goquery.NewDocumentFromReader(result.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка получения документа из страницы %s: %v", link, err)
	}
	return page.Find(selector).Find("p").Text(), nil
}

// Парсим RSS ленту
func parsingSource(s string) (a Rss, e error) {
	ind := strings.Index(s, "!")
	URL := s[:ind]
	selector := s[ind+1:]
	a.Selector = selector
	client := &http.Client{}
	res, err := client.Get(string(URL))
	if err != nil {
		return a, fmt.Errorf("не смог подключиться к RSS %s - %v", URL, err)
	}
	byteValue, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return a, fmt.Errorf("ошибка чтения %s: статус - %s", URL, res.Status)
	}
	err = xml.Unmarshal(byteValue, &a)
	if err != nil {
		return a, fmt.Errorf("ошибка парсинга ленты %s: %v", URL, err)
	}
	return a, nil
}

// парсим профиль
func parsingProfile(p Profile, mongoClient *mongo.Client) (string, error) {
	if len(p.Keys) == 0 {
		return p.Last, fmt.Errorf("список ключей %s пуст", p.Name)
	}
	node := words.Tree(p.Keys)
	var result error
	LastSearch, err := time.Parse(time.RFC1123Z, p.Last)
	if err != nil {
		result = fmt.Errorf("проблемы с парсингом даты последнего поиска %s - %s: %v", p.Name, p.Last, err)
		LastSearch = time.Now().Add(-(time.Hour * 24))
	}
	newSearch := time.Now().Format(time.RFC1123Z)
	for i := range p.Source {
		answer, err := parsingSource(p.Source[i])
		if err != nil {
			result = fmt.Errorf("%v - %v", result, err)
		}
		for i := range answer.Channels.News {
			news := &answer.Channels.News[i]
			news.parsingNews(&LastSearch, &answer.Selector, node, &newSearch, mongoClient, &p.Name)
		}
	}
	return newSearch, result
}
