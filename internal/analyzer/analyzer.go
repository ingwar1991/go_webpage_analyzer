package analyzer

import (
    "fmt"
    "net/http"
    "time"
    "github.com/PuerkitoBio/goquery"
)

type Result struct {
    HTMLVersion    string
    Title          string
    Headings       map[string]int
    InternalLinks  int
    ExternalLinks  int
    BrokenLinks    int
    HasLoginForm   bool
}

func AnalyzePage(pageURL string) (*Result, int, error) {
    client := http.Client{
        Timeout: 30 * time.Second,
    }

    resp, err := client.Get(pageURL)
    if err != nil {
        return nil, 0, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return nil, resp.StatusCode, fmt.Errorf("HTTP Error %d", resp.StatusCode)
    }

    doc, root, err := ParseHTML(resp.Body)
    if err != nil {
        return nil, resp.StatusCode, err
    }

    result := &Result{
        Headings: make(map[string]int),
    }

    result.HTMLVersion = DetectHTMLVersion(root)  

    // Title
    result.Title = doc.Find("title").Text()

    // Headings
    doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
        tag := goquery.NodeName(s)
        result.Headings[tag]++
    })

    result.InternalLinks, result.ExternalLinks, result.BrokenLinks = AnalyzeLinks(pageURL, doc)

    // Login form detection
    doc.Find("form").EachWithBreak(func(_ int, f *goquery.Selection) bool {
        if f.Find("input[type='password']").Length() > 0 {
            result.HasLoginForm = true

            return false
        }
        return true
    })

    return result, 200, nil
}
