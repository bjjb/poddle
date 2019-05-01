package main

import (
	"encoding/json"
	"encoding/xml"
	"testing"
)

func TestImage(t *testing.T) {
	t.Run("from XML", func(t *testing.T) {
		i := &Image{}
		err := xml.Unmarshal([]byte(`<image title='T'><url>U</url></image>`), i)
		if err != nil {
			t.Fatal(err)
		}
		if i.Title != "T" {
			t.Errorf("expected %q, got %q", "T", i.Title)
		}
		if i.URL != "U" {
			t.Errorf("expected %q, got %q", "U", i.URL)
		}
	})
	t.Run("to JSON", func(t *testing.T) {
		i := &Image{Title: "T", URL: "U", Type: "t"}
		b, err := json.Marshal(i)
		if err != nil {
			t.Fatal(err)
		}
		expected := `{"src":"U","title":"T","type":"t"}`
		if string(b) != expected {
			t.Errorf("expected %q, got %q", expected, b)
		}
	})
}
