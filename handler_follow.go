package main

import (
	"context"
	"fmt"
	"time"
	"database/sql"

	"github.com/lib/pq"
	"github.com/Jesbr/BlogAggregator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	 if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: follow <feed_url>")
	 }

	url := cmd.Args[0]

	ctx := context.Background()

	// get feed by URL
	feed, err := s.db.GetFeedByURL(ctx, url)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("feed not found")
		}
		return err
	}

	now := time.Now()

	// create feed follow
	follow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		// duplicate follow (unique constraint)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("already following this feed")
		}
		return err
	}

	fmt.Printf("User %s is now following %s\n", follow.UserName, follow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	// get feed follows
	follows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get follows: %w", err)
	}

	if len(follows) == 0 {
		fmt.Println("You are not following any feeds")
		return nil
	}

	for _, f := range follows {
		fmt.Println(f.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: unfollow <feed_url>")
	 }

	url := cmd.Args[0]

	ctx := context.Background()
	_, err := s.db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		Name: user.Name,
		Url: url,	
		})

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("you are not following this feed")
		}
		return fmt.Errorf("couldn't unfollow feed: %w", err)
	}

	fmt.Printf("%s unfollowed feed %s\n",user.Name, url)
	return nil
}