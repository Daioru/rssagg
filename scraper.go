package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/daioru/rssagg/internal/database"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBerweenRequest time.Duration,
) {
	log.Printf("Scraping on %v gorutines every %s duration", concurrency, timeBerweenRequest)
	ticker := time.NewTicker(timeBerweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("error fetching feeds:", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetch(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("Found post", item.Title, "on feed", feed.Name)
	}
	log.Printf("Feed %s collected, %v posts found\n", feed.Name, len(rssFeed.Channel.Item))
}