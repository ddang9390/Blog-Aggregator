package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}

func fetchFeeds(url string) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var rss RSS
	err = xml.Unmarshal(data, &rss)
	if err != nil {
		fmt.Println(err)
		return
	}
}
