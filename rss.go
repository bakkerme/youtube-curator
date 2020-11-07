package main

type Thumbnail struct {
	URL    string `xml:"url,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

type MediaGroup struct {
	Title       string    `xml:"title"`
	Thumbnail   Thumbnail `xml:"thumbnail"`
	Description string    `xml:"description"`
}

type Link struct {
	Href string `xml:"href,attr"`
}

type Entry struct {
	ID         string     `xml:"id"`
	Title      string     `xml:"title"`
	Link       Link       `xml:"link"`
	Published  string     `xml:"published"`
	Updated    string     `xml:"updated"`
	MediaGroup MediaGroup `xml:"group"`
}

type RSS struct {
	ID    string  `xml:"id"`
	Title string  `xml:"title"`
	Entry []Entry `xml:"entry"`
}
