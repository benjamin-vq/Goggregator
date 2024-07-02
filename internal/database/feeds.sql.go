// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: feeds.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id, createdat, updatedat)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, url, user_id, createdat, updatedat, lastfetchedat
`

type CreateFeedParams struct {
	ID        uuid.UUID
	Name      string
	Url       string
	UserID    uuid.UUID
	Createdat time.Time
	Updatedat time.Time
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.ID,
		arg.Name,
		arg.Url,
		arg.UserID,
		arg.Createdat,
		arg.Updatedat,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.Createdat,
		&i.Updatedat,
		&i.Lastfetchedat,
	)
	return i, err
}

const createFeedFollow = `-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, user_id, feed_id, createdat, updatedat)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, feed_id, createdat, updatedat
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FeedID    uuid.UUID
	Createdat time.Time
	Updatedat time.Time
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.ID,
		arg.UserID,
		arg.FeedID,
		arg.Createdat,
		arg.Updatedat,
	)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.FeedID,
		&i.Createdat,
		&i.Updatedat,
	)
	return i, err
}

const deleteFeedFollow = `-- name: DeleteFeedFollow :one
DELETE FROM feed_follows
WHERE feed_id = $1 AND user_id = $2
RETURNING id, user_id, feed_id, createdat, updatedat
`

type DeleteFeedFollowParams struct {
	FeedID uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteFeedFollow(ctx context.Context, arg DeleteFeedFollowParams) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, deleteFeedFollow, arg.FeedID, arg.UserID)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.FeedID,
		&i.Createdat,
		&i.Updatedat,
	)
	return i, err
}

const getAllFeeds = `-- name: GetAllFeeds :many
SELECT id, name, url, user_id, createdat, updatedat, lastfetchedat FROM feeds
`

func (q *Queries) GetAllFeeds(ctx context.Context) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getAllFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.Createdat,
			&i.Updatedat,
			&i.Lastfetchedat,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFeedFollowsByUser = `-- name: GetFeedFollowsByUser :many
SELECT id, user_id, feed_id, createdat, updatedat FROM feed_follows
WHERE user_id = $1
`

func (q *Queries) GetFeedFollowsByUser(ctx context.Context, userID uuid.UUID) ([]FeedFollow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedFollow
	for rows.Next() {
		var i FeedFollow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.FeedID,
			&i.Createdat,
			&i.Updatedat,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNextFeedsToFetch = `-- name: GetNextFeedsToFetch :many
SELECT id, name, url, user_id, createdat, updatedat, lastfetchedat FROM feeds
ORDER BY lastfetchedat NULLS FIRST
LIMIT $1
`

func (q *Queries) GetNextFeedsToFetch(ctx context.Context, limit int32) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getNextFeedsToFetch, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.Createdat,
			&i.Updatedat,
			&i.Lastfetchedat,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markFeedFetched = `-- name: MarkFeedFetched :exec
UPDATE feeds
SET lastfetchedat = $1, updatedat = $2
WHERE feeds.id = $3
`

type MarkFeedFetchedParams struct {
	Lastfetchedat sql.NullTime
	Updatedat     time.Time
	ID            uuid.UUID
}

func (q *Queries) MarkFeedFetched(ctx context.Context, arg MarkFeedFetchedParams) error {
	_, err := q.db.ExecContext(ctx, markFeedFetched, arg.Lastfetchedat, arg.Updatedat, arg.ID)
	return err
}
