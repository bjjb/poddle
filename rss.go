package main

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strings"
)

// An RSSFeed captures data from an RSS feed. The zero value has no items. You
// can populate an RSSFeed by setting its URL and calling its Fetch method,
// or from an io.Reader (such as a http.Response's Body) with Parse.
type RSSFeed struct {
	XMLName     xml.Name `xml:"rss"`
	URL         string   `xml:"-"`
	Version     string   `xml:"version,attr"`
	Title       string   `xml:"channel>title"`
	Language    string   `xml:"channel>language"`
	Description string   `xml:"channel>description"`
	PubDate     string   `xml:"channel>pubDate"`
	Items       []struct {
		XMLName     xml.Name `xml:"item"`
		Title       string   `xml:"title"`
		Description string   `xml:"description"`
		PubDate     string   `xml:"pubDate"`
		Enclosure   struct {
			XMLName xml.Name `xml:"enclosure"`
			URL     string   `xml:"url,attr"`
			Length  string   `xml:"length,attr"`
			Type    string   `xml:"type,attr"`
		} `xml:"enclosure"`
		Image Image `xml:"image" json:"image"`
	} `xml:"channel>item"`
	Image Image `xml:"channel>image"`
}

// Fetch fetches and parses the RSS feed specified by the feed's URL. It will
// return an error if the URL is blank or invalid, or if the feed is
// unparsable.
func (f *RSSFeed) Fetch() (err error) {
	resp, err := http.Get(f.URL)
	if err != nil {
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	return f.Parse(resp.Body)
}

// Parse parses the content read from the reader, and returns a *Feed, in
// which the strings are trimmed of surrounding whitespace.
func (f *RSSFeed) Parse(r io.Reader) (err error) {
	err = xml.NewDecoder(r).Decode(f)
	if err != nil {
		return
	}
	f.Version = strings.TrimSpace(f.Version)
	f.Title = strings.TrimSpace(f.Title)
	f.Language = strings.TrimSpace(f.Language)
	f.Description = strings.TrimSpace(f.Description)
	f.Image.URL = strings.TrimSpace(f.Image.URL)
	f.Image.Title = strings.TrimSpace(f.Image.Title)
	if err != nil {
		return
	}
	for _, i := range f.Items {
		i.Title = strings.TrimSpace(i.Title)
		i.Description = strings.TrimSpace(i.Description)
		i.PubDate = strings.TrimSpace(i.PubDate)
		i.Enclosure.URL = strings.TrimSpace(i.Enclosure.URL)
		i.Enclosure.Type = strings.TrimSpace(i.Enclosure.Type)
		i.Image.URL = strings.TrimSpace(i.Image.URL)
		i.Image.Title = strings.TrimSpace(i.Image.Title)
		if err != nil {
			return
		}
	}
	return
}
