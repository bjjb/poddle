package main

import "encoding/xml"

// An Image contains data from an <image>'s attributes
type Image struct {
	XMLName xml.Name `xml:"image" json:"-"`
	URL     string   `xml:"url" json:"src"`
	Title   string   `xml:"title,attr" json:"title"`
	Type    string   `xml:"-" json:"type"`
}
