package main

import (
	"blog-aggregator/backend/internal/database"
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sync"
	"time"
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
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(url string) (*RSS, error) {
	status_codes := []int{200, 201, 204}

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if response.Body == nil {
		return nil, err
	}
	if !slices.Contains(status_codes, response.StatusCode) {
		return nil, err
	}

	data, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return nil, err
	}

	var rss RSS
	err = xml.Unmarshal(data, &rss)
	if err != nil {
		return nil, err
	}

	return &rss, nil
}

func fetchWorker(cfg *apiConfig, numFeeds int, t time.Duration) error {
	timer := time.NewTicker(t)

	for _ = range timer.C {
		ctx := context.Background()
		feeds, err := cfg.DB.GetNextFeedsToFetch(ctx, int32(numFeeds))

		if err != nil {
			return err
		}
		var wg sync.WaitGroup
		fmt.Println(feeds)
		for _, feed := range feeds {
			wg.Add(1)
			url := feed.Url
			feedId := feed.ID

			go func(url string, feedId string) {
				defer wg.Done()
				rss, err := fetchFeed(url)
				if err != nil {
					fmt.Println("Error getting feeds from the url")
				} else {
					_, err = cfg.DB.MarkFeedFetched(ctx, url)
					if err != nil {
						fmt.Println("Error getting feeds from the url")
					} else {
						for _, item := range rss.Channel.Items {
							timeLayout := "2006-01-02T15:04:05.999999Z"

							var t2 time.Time
							if item.PubDate != "" {
								t2, err = time.Parse(timeLayout, item.PubDate)
								if err != nil {
									fmt.Println("Error parsing time")
								}
							}
							_, err = cfg.DB.CreatePost(ctx, database.CreatePostParams{
								Title:       sql.NullString{String: item.Title, Valid: false},
								Url:         item.Link,
								Description: sql.NullString{String: item.Description, Valid: false},
								PublishedAt: sql.NullTime{Time: t2, Valid: true},
								FeedID:      feedId,
							})
							if err != nil {
								fmt.Println(err)
							}

						}
					}
				}
			}(url, feedId)
		}
		wg.Wait()
	}
	return nil
}
