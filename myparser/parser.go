package myparser

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"mysearch/words"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Метод для парсинга конкретной новости из xml
func (n *News) parsingNews(ls *string, sel *string, node *words.Branch,
	ns *string, db *sql.DB, p *Profile) error {
	t, err := time.Parse(time.RFC1123Z, n.PubDate)
	if err != nil {
		return fmt.Errorf("проблемы с парсингом даты %s: %v", n.PubDate, err)
	}
	last, err := time.Parse(time.RFC3339, *ls)
	if err != nil {
		return fmt.Errorf("проблемы с парсингом даты %s: %v", last, err)
	}
	if t.After(last) {
		text, err := getNews(n.Link, *sel)
		if err != nil {
			return fmt.Errorf("проблемы с парсингом страницы: %v", err)
		}
		dictionary := words.SortUniq(strings.Split(text, " "))
		relev := words.SearchKeys(node, dictionary)
		if relev > 1 {
			_, err := db.Exec(fmt.Sprintf(`INSERT INTO "Search"."Posts" ("title", "text", "pub_date", "relev", "url", "is_in_report", "fresh", "profile_id", "source_id", "search_date") values ('%s', '%s', '%s', %f, '%s', '%v', '%v', %d, %d, '%s')`, n.Title, text, n.PubDate, relev, n.Link, false, false, p.ID, p.Source.ID, *ns))
			if err != nil {
				return fmt.Errorf("проблемы со вставкой новости в базу - %v", err)
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
func parsingSource(URL string) (a Rss, e error) {
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
func parsingProfile(p Profile, db *sql.DB) (string, error) {
	if len(p.Keys) == 0 {
		return p.LastSearch, fmt.Errorf("список ключей %s пуст", p.Name)
	}
	query := fmt.Sprintf(`SELECT "id","url", "selector" FROM "Search"."Sources" WHERE profile_id = %d`, p.ID)
	sources, err := db.Query(query)
	if err != nil {
		return p.LastSearch, fmt.Errorf("проблемы с получением списка источников -%v", err)
	}
	defer sources.Close()
	node := words.Tree(p.Keys)
	var result error
	newSearch := time.Now().GoString()
	for sources.Next() {
		err = sources.Scan(&p.Source.ID, &p.Source.URL, &p.Source.Selector)
		if err != nil {
			return p.LastSearch, fmt.Errorf("проблемы с списка источников - %v", err)
		}

		answer, err := parsingSource(p.Source.URL)
		if err != nil {
			result = fmt.Errorf("%v - %v", result, err)
			continue
		}
		for i := range answer.Channels.News {
			news := &answer.Channels.News[i]
			err = news.parsingNews(&p.LastSearch, &p.Source.Selector, node, &newSearch, db, &p)
			if err != nil {
				result = fmt.Errorf("%v - %v", result, err)
				continue
			}
		}
	}
	return newSearch, result
}
