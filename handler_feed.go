package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/Jesbr/BlogAggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	// get arguments
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	ctx := context.Background()

	// create feed
	now := time.Now()
	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		// handle duplicate URL (unique constraint)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("feed with this URL already exists")
		}
		return err
	}

	// automatically follow the feed
	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		// edge case: already followed
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			// ignore or warn instead of failing
			fmt.Println("Feed created, but you were already following it")
			return nil
		}
		return err
	}

	fmt.Println("Feed created successfully:")
	fmt.Printf("Name: %s\n", feed.Name)
	fmt.Printf("URL: %s\n", feed.Url)

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't list feeds: %w", err)
	}
	if len(feeds) == 0 {
		fmt.Errorf("No feeds found")
		return nil
	}
	for _, feed := range feeds {
		fmt.Printf("Feed: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("Added by: %s\n", feed.UserName)
		fmt.Println("------------------------------------")
	}
	return nil
}