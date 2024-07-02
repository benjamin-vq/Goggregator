package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"github.com/benjamin-vq/goggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
		return
	}
	assert(int32(len(feeds)) == limit, "Feeds length and feeds limit do not match")
	var wg sync.WaitGroup

	responses := make(chan FeedRss, len(feeds))
	for _, feed := range feeds {
		wg.Add(1)
		go func(feedUrl string, feedId uuid.UUID) {

			defer wg.Done()
			rss, err := fetchFeed(feedUrl)
			if err != nil {
				log.Printf("Error in goroutine of fetch feed url %s -> %q",
					feedUrl, err)
				return
			}
			responses <- FeedRss{FeedId: feedId, Rss: rss}
		}(feed.Url, feed.ID)
	}

	wg.Wait()
	close(responses)

	for feedRss := range responses {
		cfg.markAsFetched(feedRss)
		for _, posts := range feedRss.Channel.Item {
			parsedTime, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", posts.PubDate)
			if err != nil {
				log.Printf("Could not parse publish date of post %s: %q",
					posts.Title, err)
				continue
			}
			// TODO: Batch insert
			err = cfg.DB.CreatePost(context.Background(), database.CreatePostParams{
				ID:          uuid.New(),
				Title:       posts.Title,
				Url:         posts.Link,
				Description: sql.NullString{String: posts.Description, Valid: true},
				Publishedat: sql.NullTime{Time: parsedTime, Valid: true},
				FeedID:      feedRss.FeedId,
				Createdat:   time.Now(),
				Updatedat:   time.Now(),
			})
			if err != nil {
				var pqerr *pq.Error
				if errors.As(err, &pqerr) && pqerr.Code == "23505" {
					//log.Printf("Post with title %s already exists in database.",
					//	posts.Title)
					continue
				}
				log.Printf("Could not create post titled %s: %q (error type: %T)",
					posts.Title, err, err)
			} else {
				log.Printf("Successfully created post titled %s in database",
					posts.Title)
			}
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
	res, err := http.Get(url)
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
