package main

import (
	"database/sql"
	"github.com/benjamin-vq/goggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	readinessPath        = "/v1/healthz"
	errorPath            = "/v1/err"
	createUserPath       = "POST /v1/users"
	getUserPath          = "GET /v1/users"
	createFeedPath       = "POST /v1/feeds"
	getAllFeedsPath      = "GET /v1/feeds"
	createFeedFollowPath = "POST /v1/feed_follows"
	deleteFeedFollowPath = "DELETE /v1/feed_follows/{feedFollowID}"
	getAllFFByUserPath   = "GET /v1/feed_follows"
	getPostsByUserPath   = "GET /v1/posts"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	setup()
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB")
	fetchTickString := os.Getenv("FETCH_EVERY")
	feedLimitString := os.Getenv("FEEDS_AMOUNT")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Could not open database: %q", err)
	}
	dbQueries := database.New(db)

	cfg := apiConfig{
		DB: dbQueries,
	}

	fetchTick, feedLimit := parseArgs(fetchTickString, feedLimitString)
	go cfg.fetchLeastUpdatedFeeds(time.Duration(fetchTick)*time.Second, int32(feedLimit))

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc(readinessPath, handlerReadiness)
	mux.HandleFunc(errorPath, handlerError)

	mux.HandleFunc(createUserPath, cfg.handlerUsersCreate)
	mux.HandleFunc(getUserPath, cfg.authMiddleware(cfg.handlerUsersGet))

	mux.HandleFunc(createFeedPath, cfg.authMiddleware(cfg.handlerFeedsCreate))
	mux.HandleFunc(getAllFeedsPath, cfg.handlerFeedsGetAll)
	mux.HandleFunc(createFeedFollowPath, cfg.authMiddleware(cfg.handlerFeedFollowCreate))
	mux.HandleFunc(deleteFeedFollowPath, cfg.authMiddleware(cfg.handlerFeedFollowDelete))
	mux.HandleFunc(getAllFFByUserPath, cfg.authMiddleware(cfg.handlerFFByUser))
	mux.HandleFunc(getPostsByUserPath, cfg.authMiddleware(cfg.handlerGetPostsByUser))

	log.Printf("Starting server on port :%s", port)
	log.Fatal(server.ListenAndServe())
}

func setup() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could not load environment variables: %q", err)
	}
}

func parseArgs(fetchTickString, feedLimitString string) (fetchTick, feedLimit int) {

	fetchTick, err := strconv.Atoi(fetchTickString)
	if err != nil {
		log.Fatalf("Invalid fetch tick environment variable: %s, exiting.",
			fetchTickString)
	}

	feedLimit, err = strconv.Atoi(feedLimitString)
	if err != nil {
		log.Fatalf("Invalid feed limit environment variable: %s, exiting.",
			feedLimitString)
	}

	return fetchTick, feedLimit
}
