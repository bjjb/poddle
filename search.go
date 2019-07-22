package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// A Searcher can be used to configure how to search for podcasts.
type Searcher interface {
	NewRequest(string) (*http.Request, error)
	ParseResponse(*http.Response) ([]*Podcast, error)
}

type iTunesSearcher struct{}

var defaultSearcher Searcher

var defaultSearchClient *SearchClient

// A SearchClient can search for podcasts.
type SearchClient struct {
	Searcher Searcher
	Client   *http.Client
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a podcast",
	Long:  `Search for a podcast matching the given query.`,
	Run: func(c *cobra.Command, args []string) {
		if args == nil || len(args) < 1 {
			fmt.Fprintf(c.OutOrStderr(), "Missing search query\n")
			return
		}
		podcasts, err := defaultSearchClient.Search(strings.Join(args, " "))
		if err != nil {
			fmt.Fprintf(c.OutOrStderr(), "Error: %q", err)
		}
		for _, p := range podcasts {
			if p.URL == "" {
				continue
			}
			fmt.Fprintf(c.OutOrStdout(), "%s [%s]\n", p.Title, p.URL)
		}
	},
}

// Search performs a search for podcasts.
func (sc *SearchClient) Search(q string) ([]*Podcast, error) {
	s := sc.Searcher
	c := sc.Client

	if c == nil {
		c = &http.Client{}
	}

	if s == nil {
		s = defaultSearcher
	}

	req, err := s.NewRequest(q)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	results, err := s.ParseResponse(resp)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *iTunesSearcher) NewRequest(q string) (*http.Request, error) {
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

func (s *iTunesSearcher) ParseResponse(r *http.Response) ([]*Podcast, error) {
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

func init() {
	defaultSearcher = &iTunesSearcher{}
	defaultSearchClient = &SearchClient{}
	cmd.AddCommand(searchCmd)
}
