package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"

    "golang.org/x/net/html"
)

type Scraper struct {
    URL string
}

func NewScraper(url string) *Scraper {
    return &Scraper{
        URL: url,
    }
}

func (s *Scraper) Fetch() (io.Reader, error) {
    resp, err := http.Get(s.URL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
    }

    return resp.Body, nil
}

func (s *Scraper) Parse(r io.Reader) error {
    doc, err := html.Parse(r)
    if err != nil {
        return err
    }

    var data strings.Builder
    var extractData func(*html.Node)
    extractData = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "p" {
            for c := n.FirstChild; c != nil; c = c.NextSibling {
                if c.Type == html.TextNode {
                    data.WriteString(c.Data)
                    data.WriteString("\n")
                }
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            extractData(c)
        }
    }
    extractData(doc)

    err = s.SaveToFile(data.String())
    if err != nil {
        return err
    }

    return nil
}

func (s *Scraper) SaveToFile(data string) error {
    file, err := os.Create("scraped_data.txt")
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.WriteString(data)
    if err != nil {
        return err
    }

    fmt.Println("Data saved to 'scraped_data.txt'")
    return nil
}
