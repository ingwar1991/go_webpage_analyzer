package analyzer

import (
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
    "github.com/PuerkitoBio/goquery"
)

func AnalyzeLinks(pageURL string, doc *goquery.Document) (internal int, external int, broken int) {
    baseURL, _ := url.Parse(pageURL)
    links := []string{}
    doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
        href, _ := s.Attr("href")
        link, err := url.Parse(href)
        if err != nil || href == "" {
            return
        }

        if link.IsAbs() {
            external++
        } else {
            internal++
        }

        // check if accessible
        fullURL := link
        if !link.IsAbs() {
            fullURL = baseURL.ResolveReference(link)
        }
        links = append(links, fullURL.String())
    })

    return internal, external, checkBrokenLinks(links)
}

func checkBrokenLinks(links []string) int {
	var brokenCount int64
	var wg sync.WaitGroup

	// limit number of ongoing routines
	semaphore := make(chan struct{}, 10)
	for _, link := range links {
		wg.Add(1)

		go func(link string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			client := &http.Client{
				Timeout: 5 * time.Second,
			}

			resp, err := client.Head(link)
			if err != nil || resp.StatusCode >= 400 {
				atomic.AddInt64(&brokenCount, 1)
			}
		}(link)
	}
	wg.Wait()

	return int(brokenCount)
}
