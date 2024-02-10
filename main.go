package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

var qMutex sync.RWMutex
var quotes []string
var wg sync.WaitGroup

func main() {
	var start = time.Now().Second()
	var urls = []string{"https://quotes.toscrape.com/", "https://quotes.toscrape.com/page/2", "https://quotes.toscrape.com/page/3", "https://quotes.toscrape.com/4"}
	wg.Add(len(urls))

	fo, err := os.Create("quotes.json")
	if err != nil {
		log.Fatal(err)
	}

	for _, url := range urls {
		go crawlQuotes(&wg, url)
	}

	wg.Wait()
	qjson, err := json.Marshal(map[string][]string{"quotes": quotes})
	if err != nil {
		log.Fatal(err)
	}
	fo.Write(qjson)
	fmt.Println(time.Now().Second() - start)
}

func crawlQuotes(wg *sync.WaitGroup, url string) {

	defer wg.Done()

	c := colly.NewCollector(colly.AllowURLRevisit(), colly.MaxDepth(100))

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Url of page is ", r.URL)
	})

	c.OnHTML(".quote", func(h *colly.HTMLElement) {
		quote := h.ChildText("span.text")
		qMutex.Lock()
		quotes = append(quotes, quote)
		qMutex.Unlock()
	})
	c.Visit(url)

}
