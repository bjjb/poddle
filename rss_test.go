package main

import (
	"strings"
	"testing"
)

func TestRSS(t *testing.T) {

	const xml = `
		<rss xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd"
				 xmlns:googleplay="http://www.google.com/schemas/play-podcasts/1.0"
				 xmlns:atom="http://www.w3.org/2005/Atom"
				 version="2.0">
			<channel>
				<atom:link href="https://example.com/1/test.xml"
									 rel="self"
									 type="application/rss+xml"/>
				<title>My Feed</title>
				<link/>
				<language>en</language>
				<copyright/>
				<description>My Feed Description</description>
				<pubDate>Mon, 22 Apr 2019 00:00:00 -0000</pubDate>
				<image>
					<url>
						https://example.com/1.jpg
					</url>
					<title>My Image</title>
					<link/>
				</image>
				<item>
					<title>My Episode One</title>
					<description>My Episode Description One</description>
					<pubDate>Mon, 22 Apr 2019 00:00:00 -0000</pubDate>
					<enclosure url="https://example.com/1/1.mp3"
										 length="100"
										 type="audio/mpeg"/>
				</item>
				<item>
				<title>My Episode Two</title>
				<description>My Episode Description Two</description>
				<pubDate>Mon, 15 Apr 2019 00:00:00 -0000</pubDate>
				<enclosure url="https://example.com/1/2.mp3"
									 length="200"
									 type="audio/mpeg"/>
				</item>
				<item>
				<title>My Episode Three</title>
				<description>My Episode Description Three</description>
				<pubDate>Mon, 08 Apr 2019 00:00:00 -0000</pubDate>
				<enclosure url="https://example.com/1/3.mp3"
									 length="300"
									 type="audio/mpeg"/>
				</item>
			</channel>
	</rss>`

	r := strings.NewReader(xml)
	feed := &RSSFeed{}

	if err := feed.Parse(r); err != nil {
		t.Fatal(err)
	}

	t.Run("content", func(t *testing.T) {
		if len(feed.Items) != 3 {
			t.Fatalf("expected 3 items, got %d", len(feed.Items))
		}

		tt := []struct{ expected, actual string }{
			{"2.0", feed.Version},
			{"My Feed", feed.Title},
			{"https://example.com/1.jpg", feed.Image.URL},
			{"My Feed Description", feed.Description},
			{"en", feed.Language},
			{"My Episode One", feed.Items[0].Title},
			{"My Episode Description One", feed.Items[0].Description},
			{"Mon, 22 Apr 2019 00:00:00 -0000", feed.Items[0].PubDate},
			{"https://example.com/1/1.mp3", feed.Items[0].Enclosure.URL},
			{"100", feed.Items[0].Enclosure.Length},
			{"audio/mpeg", feed.Items[0].Enclosure.Type},
			{"My Episode Two", feed.Items[1].Title},
			{"My Episode Description Two", feed.Items[1].Description},
			{"Mon, 15 Apr 2019 00:00:00 -0000", feed.Items[1].PubDate},
			{"https://example.com/1/2.mp3", feed.Items[1].Enclosure.URL},
			{"200", feed.Items[1].Enclosure.Length},
			{"audio/mpeg", feed.Items[1].Enclosure.Type},
			{"My Episode Three", feed.Items[2].Title},
			{"My Episode Description Three", feed.Items[2].Description},
			{"Mon, 08 Apr 2019 00:00:00 -0000", feed.Items[2].PubDate},
			{"https://example.com/1/3.mp3", feed.Items[2].Enclosure.URL},
			{"300", feed.Items[2].Enclosure.Length},
			{"audio/mpeg", feed.Items[2].Enclosure.Type},
		}
		for _, tc := range tt {
			if tc.actual != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, tc.actual)
			}
		}
	})
}
