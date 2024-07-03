package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strings"
)

type Scraper struct {
    URL string
}

func NewScraper(url string) *Scraper {
    return &Scraper{
        URL: url,
    }
}

func (s *Scraper) Fetch() ([]byte, error) {
    resp, err := http.Get(s.URL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return ioutil.ReadAll(resp.Body)
}

func (s *Scraper) Parse(body []byte) error {

    data := extractData(string(body))
    return s.SaveToFile(data)
}

func (s *Scraper) SaveToFile(data string) error {
    return ioutil.WriteFile("scraped_data.txt", []byte(data), 0644)
}

func extractData(html string) string {

    return "Extracted data: " + strings.ToUpper(html[:50])
}

func main() {
    fmt.Print("Enter the URL to scrape: ")
    reader := bufio.NewReader(os.Stdin)
    url, _ := reader.ReadString('\n')
    url = strings.TrimSpace(url)

    scraper := NewScraper(url)

    body, err := scraper.Fetch()
    if err != nil {
        log.Fatalf("Error fetching web page: %v", err)
    }

    err = scraper.Parse(body)
    if err != nil {
        log.Fatalf("Error parsing web page: %v", err)
    }

    fmt.Println("Scraping completed. Data saved to scraped_data.txt")
}
