package main

import (
	"archive/zip"
	"black_hat_go/bing-metadata/metadata"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func handler(i int, s *goquery.Selection) {
	url, ok := s.Find("a").Attr("href")
	if !ok {
		log.Println("Not ok")
		return
	}
	fmt.Printf("%d: %s\n", i, url)
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	r, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		log.Println(err)
		return
	}
	cp, ap, err := metadata.NewProperties(r)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf(
		"%21s %s - %s %s\n",
		cp.Creator,
		cp.LastModifiedBy,
		ap.Application,
		ap.GetMajorVersion())
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalln("Missing required argument. Usage: main.go <domain> <ext>")
	}
	domain := os.Args[1]
	fileType := os.Args[2]
	// query dosn't work
	q := fmt.Sprintf("site:%s && filetype:%s && instreamset:(url title):%s", domain, fileType, fileType)
	search := fmt.Sprintf("http://bing.com/search?q=%s", url.QueryEscape(q))

	log.Println(search)
	res, err := http.Get(search)
	if err != nil {
		log.Panic(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Panic(err)
	}
	s := "html body div#b_content ol#b_results li.b_algo div.b_title h2"
	doc.Find(s).Each(handler)
}
