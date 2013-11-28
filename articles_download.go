package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

func UpdateArticles() {
	log.Println("Updating articles")
	DB, err := sql.Open("sqlite3", ExePath+"/db.db")
	if err != nil {
		log.Println("Article Download Error | Error opening database: " + err.Error())
		return
	}
	tx, err := DB.Begin()
	if err != nil {
		log.Println("Article Download Error | Error starting transaction: " + err.Error())
		return
	}
	defer DB.Close()
	rows, _ := DB.Query("select url, id from subscriptions")
	subs := make([]Subscription, 0)
	for rows.Next() {
		var sub Subscription
		rows.Scan(&sub.Url, &sub.Id)
		subs = append(subs, sub)
	}
	rows.Close()
	articleChannel := make(chan *[]Article, MAX_ARTICLE_PROCS)
	sem := make(chan int, MAX_ARTICLE_PROCS)
	// init semaphore
	for i := 0; i < MAX_ARTICLE_PROCS; i++ {
		sem <- 1
	}
	var waitGroupTwo sync.WaitGroup
	go processArticles(&articleChannel, tx, &waitGroupTwo)
	var waitGroup sync.WaitGroup
	for i := 0; i < len(subs); i++ {
		<-sem
		waitGroup.Add(1)
		go fetchAndFormatArticles(&sem, &subs[i], &articleChannel, &waitGroup)
	}
	waitGroup.Wait()
	close(articleChannel)
	waitGroupTwo.Wait() // TODO: PLEASE MAKE THIS ELEGANT YOU STUPID FUCK (if possible)
	err = tx.Commit()
	if err != nil {
		log.Println("Article Download Error | Commit error: " + err.Error())
		return
	}
	log.Println("Finished updating articles")
}

func processArticles(articleChannel *chan *[]Article, tx *sql.Tx, waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)
	for articles := range *articleChannel {
		rows, err := tx.Query("select subscription_id, datetime(published), url from articles where published >= datetime(?) and subscription_id = ? order by datetime(published) asc", (*articles)[0].Published.Format(time.RFC3339), (*articles)[0].SubscriptionId)
		var existingArticles []Article
		if err != nil {
			log.Println("Article Download Error | Error retrieving existing articles " + err.Error())
		}
		for rows.Next() {
			var article Article
			var dateString string
			rows.Scan(&article.SubscriptionId, &dateString, &article.Url)
			article.Published, err = time.Parse("2006-01-02 15:04:05", dateString)
			existingArticles = append(existingArticles, Article{})
			copy(existingArticles[0+1:], existingArticles[0:])
			existingArticles[0] = article
		}
		rows.Close()
		if len(existingArticles) == 0 {
			insertFinalArticleSlice(articles, tx)
		} else if (existingArticles[0].Published != (*articles)[len(*articles)-1].Published) && (existingArticles[0].Url != (*articles)[len(*articles)-1].Url) {
			index := sort.Search(len(*articles), func(i int) bool {
				return (*articles)[i].Published.Unix() > existingArticles[0].Published.Unix()
			})
			newSlice := ((*articles)[index:])
			insertFinalArticleSlice(&newSlice, tx)
		}
	}
	waitGroup.Done()
}

func insertFinalArticleSlice(articles *[]Article, tx *sql.Tx) {
	stmt, err := tx.Prepare("insert into articles values (null, ?, ?, ?, ?, datetime(?), ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Article Download Error | Insert error: " + err.Error())
	}
	for _, v := range *articles {
		_, err := stmt.Exec(v.SubscriptionId, v.Url, v.Title, v.Author, v.Published.Format(time.RFC3339), v.Body.Content, v.Body.Type, v.Summary.Content, v.Summary.Type, v.Read)
		if err != nil {
			log.Println("Articles error: " + err.Error())
		}
	}
	stmt.Close()
}

func fetchAndFormatArticles(sem *chan int, sub *Subscription, articleChannel *chan *[]Article, waitGroup *sync.WaitGroup) {
	resp, err := http.Get(sub.Url)
	if err != nil {
		log.Printf("Unable to fetch articles from %v | Error: %v ", sub.Url, err.Error()) // TODO: Manage error. How?
		*sem <- 1
		waitGroup.Done()
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	doc, err := gokogiri.ParseXml(body)
	if err != nil {
		log.Println("Error parsing XML on feed with URL " + sub.Url)
		return
	}
	defer doc.Free()
	doc.RecursivelyRemoveNamespaces()
	var articles []Article
	if doc.Root().Name() == "rss" {
		formatRSS(sub, doc, &articles)
	} else if doc.Root().Name() == "feed" {
		formatAtom(sub, doc, &articles)
	}
	*articleChannel <- &articles
	*sem <- 1
	waitGroup.Done()
}

func formatAtom(sub *Subscription, doc *xml.XmlDocument, articles *[]Article) {
	entries, err := doc.Search("//entry")
	if err != nil {
		log.Println("Error parsing Atom entries on feed: " + err.Error())
		return
	}
	for _, v := range entries {
		var article Article
		// Subscription ID
		article.SubscriptionId = sub.Id
		// Title
		titleNodes, _ := v.Search(v.Path() + "/title")
		article.Title = titleNodes[0].Content()
		// Link
		linkNodes, _ := v.Search(v.Path() + "/link")
		article.Url = linkNodes[0].Attr("href")
		// Published
		dateNodes, _ := v.Search(v.Path() + "/published")
		pub, _ := time.Parse(time.RFC3339, dateNodes[0].Content())
		article.Published = pub.UTC()
		// Author
		authorNodes, _ := v.Search(v.Path() + "/author/name")
		if len(authorNodes) > 0 {
			article.Author = authorNodes[0].Content()
		}
		// Summary
		summaryNodes, _ := v.Search(v.Path() + "/summary")
		if len(summaryNodes) > 0 {
			article.Summary.Content = summaryNodes[0].Content()
			article.Summary.Type = summaryNodes[0].Attr("type")
		}
		// Body
		bodyNodes, _ := v.Search(v.Path() + "/content")
		if len(bodyNodes) > 0 {
			article.Body.Content = bodyNodes[0].Content()
			article.Body.Type = bodyNodes[0].Attr("type")
		}
		// Read
		article.Read = false
		*articles = append(*articles, Article{})
		copy((*articles)[0+1:], (*articles)[0:])
		(*articles)[0] = article
	}
	// rfc3339
}

func formatRSS(sub *Subscription, doc *xml.XmlDocument, articles *[]Article) {
	items, err := doc.Search("//item")
	if err != nil {
		log.Println("Error parsing RSS entries on feed: " + err.Error())
		return
	}
	for _, v := range items {
		var article Article
		// Subscription ID
		article.SubscriptionId = sub.Id
		// Title
		titleNodes, _ := v.Search(v.Path() + "/title")
		article.Title = titleNodes[0].Content()
		// Link
		linkNodes, _ := v.Search(v.Path() + "/link")
		article.Url = linkNodes[0].Content()
		// Published
		dateNodes, _ := v.Search(v.Path() + "/pubDate")
		pub, err := time.Parse(time.RFC1123Z, dateNodes[0].Content())
		if err != nil {
			pub, _ = time.Parse(time.RFC1123, dateNodes[0].Content())
		}
		article.Published = pub.UTC()
		// Author
		authorNodes, _ := v.Search(v.Path() + "/author")
		if len(authorNodes) > 0 {
			article.Author = authorNodes[0].Content()
		}
		// Summary
		summaryNodes, _ := v.Search(v.Path() + "/description")
		if len(summaryNodes[0]) > 0 {
			article.Summary.Content = summaryNodes[0].Content()
			article.Summary.Type = "html"
		}
		// Body
		// No body? Body = Summary? No summary and body contains entry description?
		// Read
		article.Read = false
		*articles = append(*articles, Article{})
		copy((*articles)[0+1:], (*articles)[0:])
		(*articles)[0] = article
	}
	// rfc1123 or rfc1123z
}
