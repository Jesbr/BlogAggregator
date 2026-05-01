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

func handlerAddFeed(s *state, cmd command) error {
	// make sure a user is logged in
	if s.cfg.CurrentUserName == "" {
		return fmt.Errorf("no current user set")
	}
	
	// get arguments
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	// get the user from DB
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("current user does not exist")
		}
		return err
	}

	// create feed
	now := time.Now()
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
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

	fmt.Println("Feed created successfully:")
	fmt.Printf("Name: %s\n", feed.Name)
	fmt.Printf("URL: %s\n", feed.Url)

	return nil
}