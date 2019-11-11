package metascraper

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Scrape creates a new page and populates its fields from the content found at
// the given URL.
func Scrape(url string) (*Page, error) {
	p := &Page{
		URL:          url,
		MetaReader:   &MetaReader{},
		SchemaReader: &SchemaReader{},
	}
	// Unlike the other TokenReaders, the PageReader must manipulate its parent.
	p.PageReader = &PageReader{page: p}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return p, err
	}
	defer resp.Body.Close()
	htmlBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return p, err
	}
	p.HTML = string(htmlBytes)
	if err = p.Read(htmlBytes); err != io.EOF {
		return p, err
	}
	return p, nil
}
