package main

import (
	"testing"
	"time"
)

func TestParseXML(t *testing.T) {
	p := &Podcast{}
	err := p.ParseXML(`
		<rss>
			<channel>
				<title>T</title>
				<language>L</language>
				<description>D</description>
				<pubDate>Mon, 22 Apr 2019 00:00:00 -0000</pubDate>
				<image><url>https://example.com/1.jpg</url></image>
				<item>
					<title>E1</title>
					<description>Episode 1</description>
					<pubDate>Mon, 22 Apr 2019 00:00:00 -0000</pubDate>
					<image><url>https://example.com/1/1.png</url></image>
					<enclosure url="https://example.com/1"
										 length="100"
										 type="audio/mpeg"/>
				</item>
				<item>
					<title>E2</title>
					<description>Episode 2</description>
					<pubDate>Tue, 23 Apr 2019 00:00:00 -0000</pubDate>
					<image><url>https://example.com/1/2.png</url></image>
					<enclosure url="https://example.com/2"
										 length="101"
										 type="audio/ogg+vorbis"/>
				</item>
			</channel>
		</rss>`)
	if err != nil {
		t.Fatal(err)
	}
	apr22 := time.Date(2019, time.April, 22, 0, 0, 0, 0, time.UTC)
	apr23 := time.Date(2019, time.April, 23, 0, 0, 0, 0, time.UTC)
	for _, c := range []struct{ expected, actual interface{} }{
		{"T", p.Title},
		{"D", p.Description},
		{"L", p.Language},
		{2, len(p.Episodes)},
		{"https://example.com/1.jpg", p.Image.URL},
		{"image/jpeg", p.Image.Type},
		{"E1", p.Episodes[0].Title},
		{"Episode 1", p.Episodes[0].Description},
		{1, len(p.Episodes[0].Versions)},
		{apr22.String(), p.PublishedAt.String()},
		{"https://example.com/1", p.Episodes[0].Versions[0].URL},
		{"audio/mpeg", p.Episodes[0].Versions[0].Type},
		{"https://example.com/1/1.png", p.Episodes[0].Image.URL},
		{"image/png", p.Episodes[0].Image.Type},
		{apr22.String(), p.Episodes[0].PublishedAt.String()},
		{"E2", p.Episodes[1].Title},
		{"Episode 2", p.Episodes[1].Description},
		{1, len(p.Episodes[1].Versions)},
		{"https://example.com/2", p.Episodes[1].Versions[0].URL},
		{"audio/ogg+vorbis", p.Episodes[1].Versions[0].Type},
		{"https://example.com/1/2.png", p.Episodes[1].Image.URL},
		{"image/png", p.Episodes[1].Image.Type},
		{apr23.String(), p.Episodes[1].PublishedAt.String()},
	} {
		if c.actual != c.expected {
			t.Errorf("expected %q, got %q", c.expected, c.actual)
		}
	}
}
