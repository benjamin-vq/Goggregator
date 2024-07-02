package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"github.com/benjamin-vq/goggregator/internal/database"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func (cfg *apiConfig) fetchLeastUpdatedFeeds(fetchTick time.Duration, feeds int32) {
	log.Printf("Starting goroutine to fetch the %d least updated feeds every %.1f seconds.",
		feeds, fetchTick.Seconds())
	ticker := time.NewTicker(fetchTick)

	for range ticker.C {
		cfg.feedRequests(feeds)
	}
}

func (cfg *apiConfig) feedRequests(limit int32) {
	feeds, err := cfg.DB.GetNextFeedsToFetch(context.Background(), limit)
	if err != nil {
		log.Printf("Could not query next feeds to fetch: %q", err)
	}
	assert(int32(len(feeds)) == limit, "Feeds length and feeds limit do not match")
	var wg sync.WaitGroup

	responses := make(chan FeedRss, len(feeds))
	for i, feed := range feeds {
		wg.Add(1)
		log.Printf("Fetching feed number %d in goroutine", i)
		go func(feedUrl string, feedId uuid.UUID) {

			defer wg.Done()
			rss, err := fetchFeed(feedUrl)
			if err != nil {
				log.Printf("Error in goroutine of fetch feed url %s -> %q",
					feedUrl, err)
				return
			}
			log.Println("Putting processed RSS response in channel")
			responses <- FeedRss{FeedId: feedId, Rss: rss}
		}(feed.Url, feed.ID)
	}

	log.Println("Waiting on work group")
	wg.Wait()
	log.Println("Work group finished, closing channel")
	close(responses)

	log.Println("Processing channel elements")
	for feedRss := range responses {
		cfg.markAsFetched(feedRss)
		log.Printf("Posts fetched from %s:", feedRss.Channel.Title)
		for _, posts := range feedRss.Channel.Item {
			log.Print(posts.Title)
		}
	}
	log.Printf("All %d RSS feeds processed.", len(feeds))
}

func (cfg *apiConfig) markAsFetched(feedRss FeedRss) {
	log.Printf("Marking as fetched RSS feed of: %s (feed id: %s)",
		feedRss.Channel.Title, feedRss.FeedId)
	err := cfg.DB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		Lastfetchedat: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Updatedat: time.Now(),
		ID:        feedRss.FeedId,
	})

	if err != nil {
		log.Printf("Could not mark feed %s (feed id: %s) as fetched: %q",
			feedRss.Channel.Title, feedRss.FeedId, err)
	}
}

func fetchFeed(url string) (Rss, error) {
	log.Printf("Making a GET request to feed: %s", url)
	res, err := http.Get(url)
	log.Printf("Call to feed url %s finished", url)
	if err != nil {
		log.Printf("Error during GET request: %q", err)
		return Rss{}, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if res.StatusCode > 299 || err != nil {
		log.Printf("Bad status code or error: %q", err)
		return Rss{}, errors.New("invalid api response")
	}

	rss := Rss{}
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		log.Printf("Could not unmarshal response: %q", err)
		return Rss{}, err
	}

	return rss, nil
}
