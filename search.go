package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type searcher interface {
	newRequest(string) (*http.Request, error)
	parseResponse(*http.Response) ([]*Podcast, error)
}

type iTunesSearcher struct{}

func (s *iTunesSearcher) newRequest(q string) (*http.Request, error) {
	if q == "" {
		return nil, fmt.Errorf("q cannot be blank")
	}
	if len(q) > 255 {
		return nil, fmt.Errorf("q cannot be longer than 255 characters")
	}

	u := "https://itunes.apple.com/search?entity=podcast&term=%s"
	u = fmt.Sprintf(u, url.QueryEscape(q))

	req, _ := http.NewRequest(http.MethodGet, u, nil)
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (s *iTunesSearcher) parseResponse(r *http.Response) ([]*Podcast, error) {
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search returned %d (%s)", r.StatusCode, r.Status)
	}
	defer r.Body.Close()
	var results struct {
		ResultCount int `json:"result_count"`
		Results     []struct {
			CollectionName   string    `json:"collectionName"`
			TrackName        string    `json:"trackName"`
			FeedURL          string    `json:"feedUrl"`
			TrackViewURL     string    `json:"trackViewUrl"`
			ArtworkURL30     string    `json:"artworkUrl30"`
			ArtworkURL60     string    `json:"artworkUrl60"`
			ArtworkURL100    string    `json:"artworkUrl100"`
			ArtworkURL600    string    `json:"artworkUrl600"`
			ReleaseDate      time.Time `json:"releaseDate"`
			TrackCount       int       `json:"trackCount"`
			Country          string    `json:"country"`
			PrimaryGenreName string    `json:"primaryGenreName"`
			GenreIDs         []string  `json:"genreIds"`
			Genres           []string  `json:"genres"`
		} `json:"results"`
	}
	if err := json.NewDecoder(r.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %s", err)
	}
	podcasts := []*Podcast{}
	firstString := func(ss ...string) (s string) {
		for _, s = range ss {
			if s != "" {
				break
			}
		}
		return
	}
	for _, result := range results.Results {
		img := firstString(result.ArtworkURL600, result.ArtworkURL100,
			result.ArtworkURL60, result.ArtworkURL30)
		podcasts = append(podcasts, &Podcast{
			Title:       result.CollectionName,
			URL:         result.FeedURL,
			Image:       Image{URL: img},
			PublishedAt: result.ReleaseDate,
		})
	}
	return podcasts, nil
}
