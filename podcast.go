package main

import (
	"io"
	"log"
	"mime"
	"path/filepath"
	"strings"
	"time"
)

// A Podcast is a podcast
type Podcast struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Language    string    `json:"language"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Image       Image     `json:"image"`
	Episodes    []Episode `json:"episodes"`
}

// An Episode is a Podcast episode
type Episode struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Image       Image     `json:"image"`
	Versions    []Version `json:"versions"`
}

// A Version is a specific version of an episode.
type Version struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

// ParseXML parses an XML string into a new Podcast
func (p *Podcast) ParseXML(xml string) error {
	return p.ParseRSS(strings.NewReader(xml))
}

// ParseRSS parses an RSS feed into a new Podcast
func (p *Podcast) ParseRSS(r io.Reader) (err error) {
	f := &RSSFeed{}
	if err = f.Parse(r); err == nil {
		p.FromRSSFeed(f)
	}
	return
}

// FromRSSFeed injects the feed data into the Podcast. The episodes will be
// set to the items of the feed, and each item will have 1 version (defined by
// the feed item's <enclosure>.
func (p *Podcast) FromRSSFeed(f *RSSFeed) {
	trim := func(s string) string {
		return strings.TrimSpace(s)
	}
	parseTime := func(s string) time.Time {
		t, err := time.Parse(time.RFC1123Z, trim(s))
		if err != nil {
			log.Print(err)
		}
		return t.UTC()
	}
	mimeTypeOf := func(s string) string {
		return mime.TypeByExtension(filepath.Ext(trim(s)))
	}
	p.Title = trim(f.Title)
	p.Description = trim(f.Description)
	p.Language = trim(f.Language)
	p.Image.URL = trim(f.Image.URL)
	p.Image.Type = mimeTypeOf(f.Image.URL)
	p.Image.Title = trim(f.Image.Title)
	p.PublishedAt = parseTime(f.PubDate)
	p.Episodes = make([]Episode, len(f.Items))
	for i, item := range f.Items {
		p.Episodes[i] = Episode{
			Title:       trim(item.Title),
			Description: trim(item.Description),
			PublishedAt: parseTime(item.PubDate),
			Image: Image{
				URL:   trim(item.Image.URL),
				Title: trim(item.Image.Title),
				Type:  mimeTypeOf(item.Image.URL),
			},
			Versions: []Version{
				{item.Enclosure.URL, item.Enclosure.Type},
			},
		}
	}
}
