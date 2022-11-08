package myparser

import (
	"encoding/xml"
	"fmt"
	"io"
	"mysearch/words"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var client *http.Client

// Метод для парсинга конкретной новости из xml
func (n *News) ParsingNews(ls string, sel string, node *words.Branch,
	ns string, p Profile) error {
	var post Post
	t, err := time.Parse(time.RFC1123Z, n.PubDate)
	if err != nil {
		return fmt.Errorf("проблемы с парсингом даты %s: %v", n.PubDate, err)
	}
	last, err := time.Parse(time.RFC1123Z, ls)
	if err != nil {
		return fmt.Errorf("проблемы с парсингом даты %s: %v", last, err)
	}
	if t.After(last) {
		post.PubDate = n.PubDate
		post.Text, err = GetNews(n.Link, sel)
		if err != nil {
			return fmt.Errorf("проблемы с парсингом страницы: %v", err)
		}
		post.Dictionary = words.SortUniq(strings.Split(post.Text, " "))
		post.Relev = words.SearchKeys(node, post.Dictionary)
		post.Title = n.Title
		post.SearchDate = ns
		post.Link = n.Link
		if post.Relev > 1 {
			err = post.insert(p.ID, p.Source.ID, SDB.Db)
			//_, err := db.Exec(fmt.Sprintf(`INSERT INTO "Search"."Posts" ("title", "text", "pub_date", "relev", "url", "is_in_report", "fresh", "profile_id", "source_id", "search_date") values ('%s', '%s', '%s', %f, '%s', '%v', '%v', %d, %d, '%s')`, n.Title, text, n.PubDate, relev, n.Link, false, false, p.ID, p.Source.ID, *ns))
			if err != nil {
				return fmt.Errorf("проблемы со вставкой новости в базу - %v", err)
			}
		}
	}
	return nil
}

// Парсим текст новости из html страницы
func GetNews(link, selector string) (string, error) {
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
	return page.Find(selector).Text(), nil
}

// Парсим RSS ленту
func ParsingSource(p Profile) (a Rss, ok bool) {
	res, err := client.Get(string(p.Source.URL))
	lg <- fmt.Sprintf("начинаем разбор %s", p.Source.URL)
	if err != nil {
		lg <- fmt.Sprintf("не смог подключиться к RSS %s - %v", p.Source.URL, err)
		return a, false
	}
	byteValue, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		lg <- fmt.Sprintf("ошибка чтения %s: статус - %s", p.Source.URL, res.Status)
		return a, false
	}
	err = xml.Unmarshal(byteValue, &a)
	if err != nil {
		lg <- fmt.Sprintf("ошибка парсинга ленты %s: %v", p.Source.URL, err)
		return a, false
	}
	return a, true
}

// метод парсинга профиля
func (p *Profile) parsingProfile() string {
	if len(p.Keys) == 0 {
		lg <- fmt.Sprintf("список ключей %s пуст", p.Name)
		return p.LastSearch
	}
	node := words.Tree(p.Keys)
	newSearch := time.Now().Format(time.RFC1123Z)
	query := fmt.Sprintf(`SELECT "id","url", "selector" FROM "Search"."Sources" WHERE profile_id = %d`, p.ID)
	sources, err := SDB.Db.Query(query)
	if err != nil {
		lg <- fmt.Sprintf("проблемы с получением списка источников -%v", err)
		return p.LastSearch
	}
	defer sources.Close()
	client = &http.Client{}
	for sources.Next() {
		err = sources.Scan(&p.Source.ID, &p.Source.URL, &p.Source.Selector)
		if err != nil {
			lg <- fmt.Sprintf("проблемы с списка источников - %v", err)
			return p.LastSearch
		}
		a, ok := ParsingSource(*p)
		if ok {
			for i := range a.Channels.News {
				news := &a.Channels.News[i]
				err = news.ParsingNews(p.LastSearch, p.Source.Selector, node, newSearch, *p)
				if err != nil {
					lg <- fmt.Sprint(err)
					continue
				}
			}
		}
	}

	return newSearch
}
