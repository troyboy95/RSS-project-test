package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/troyboy95/RSS-project-test/internal/database"
)

func startScraper(db *database.Queries, concurrency int, timeBtwRequest time.Duration) {
	// Implementatsion for starting the scraper
	log.Printf("Starting RSS feed scraper on %v goroutines every %v seconds", concurrency, timeBtwRequest.Seconds())

	ticker := time.NewTicker(timeBtwRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("Error fetching feeds to scrape: %v", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched: %v", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching RSS feed from %s: %v", feed.Url, err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		// log.Println("Found item:", item.Title, " from feed:", feed.Name)
		description := sql.NullString{}
		if item.Description != "" {
			description = sql.NullString{String: item.Description, Valid: true}
		}

		t, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error parsing publication date %s for item %s: %v", item.PubDate, item.Title, err)
			continue
		}

		_, e := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Description: description,
			PublishedAt: t,
			Url:         item.Link,
			FeedID:      feed.ID,
		})

		if e != nil {
			if strings.Contains(e.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Error creating post for feed %s: %v", feed.Name, e)
		}
	}
	log.Printf("Successfully scraped feed: %s, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}
