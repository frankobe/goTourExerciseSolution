package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

func crawlHelper(url string, fetcher Fetcher, uChan chan []string) {
	body, urlArray, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("found: %s %q\n", url, body)
	}
	uChan <- urlArray
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	m := make(map[string]bool)
	urlChan := make(chan []string)

	go crawlHelper(url, fetcher, urlChan)
	m[url] = true
	running := 1
	for running > 0 {
		running--
		urlGroup := <-urlChan
		for _, urlStr := range urlGroup {
			if m[urlStr] == false {
				m[urlStr] = true
				running++
				go crawlHelper(urlStr, fetcher, urlChan)
			}
		}
	}
}
func CrawlMutex(url string, depth int, fetcher Fetcher) {
	m := map[string]bool{url: true}
	var mx sync.Mutex
	var wg sync.WaitGroup
	var c2 func(string, int)
	c2 = func(url string, depth int) {
		defer wg.Done()
		fmt.Println("depth: ", depth)
		if depth <= 0 {
			return
		}
		body, urls, err := fetcher.Fetch(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("found: %s %q\n", url, body)

		for _, u := range urls {
			mx.Lock()
			if !m[u] {
				m[u] = true
				mx.Unlock()
				wg.Add(1)
				go c2(u, depth-1)
			} else {
				mx.Unlock()
			}

		}
	}
	wg.Add(1)
	c2(url, depth)
	wg.Wait()
}
func main() {
	Crawl("http://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
