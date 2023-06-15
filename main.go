package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const source string = "https://www.minecraft.net/en-us/download/server/bedrock"

func checkErr(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func get(url string) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	checkErr(err)
	client := &http.Client{}
	resp, err := client.Do(req)
	checkErr(err)
	return resp
}

func getLastUrl() string {
	url, e := os.ReadFile("lasturl.txt")
	checkErr(e)
	return string(url)
}

func update() {
	log.Println("Start updating")
	resp := get(source)
	doc, err := goquery.NewDocumentFromReader(resp.Body)	
	checkErr(err)
	// fmt.Println(doc.Text())
	node := doc.Find("#main-content > div > div > div.page-section-container.aem-GridColumn.aem-GridColumn--default--12 > div > div > div > div.server-card.aem-GridColumn.aem-GridColumn--default--12 > div > div > div > div:nth-child(2) > div.card-footer > div > a")
	url, _ := node.Attr("href")
	log.Println("Fetched the url")
	log.Println(url)
	if url == getLastUrl() {
		log.Println("No update")
		return
	}
	os.WriteFile("lasturl.txt", []byte(url), 0666)
	log.Println("Downloading...")
	pack, err := ioutil.ReadAll(get(url).Body)
	checkErr(err)
	os.WriteFile("Latest_minecraft_linux.zip", pack, 0666)
	log.Println("Downloaded")
}

func main() {
	for {
		update()
		log.Println("Sleep for an hour...")
		time.Sleep(time.Hour)
	}
}