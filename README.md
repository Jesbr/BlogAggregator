# Blog Aggregator CLI (Gator)

A simple RSS blog aggregator written in Go.
It allows you to follow feeds, scrape posts, and browse content directly from your terminal.

---

## Requirements

Before running this program, make sure you have:

* **Go** (1.20 or newer recommended)
* **PostgreSQL** (running locally or remotely)

---

## Installation

Install the CLI using `go install`:

```bash
go install github.com/Jesbr/BlogAggregator@latest
```

Make sure your `$GOPATH/bin` is in your system `PATH`, so you can run the binary globally.

---

## Configuration

The app uses a config file located at:

```bash
~/.gatorconfig.json
```

### Example config:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

### Steps:

1. Create a PostgreSQL database (e.g. `gator`)
2. Run your migrations (using goose or your preferred tool)
3. Create the config file above in your home directory
4. Update `db_url` with your database credentials

---

## Usage

All commands are run through the CLI:

```bash
gator <command> [arguments]
```

---

## User Commands

### Register a new user

```bash
gator register <username>
```

Creates a new user and sets them as the current user.

---

### Login as a user

```bash
gator login <username>
```

Switches the current active user, requires user to have registered beforehand.

---

### List users

```bash
gator users
```

Shows all users. The current user will be marked with [user] (current).

---

## Feed Commands

### Add a feed

```bash
gator addfeed "<name>" "<url>"
```

Adds a new RSS feed and automatically follows it for the current user.

---

### Follow an existing feed

```bash
gator follow <feed_url>
```

Follows an RSS feed already followed by another user.

---

### Unfollow a feed

```bash
gator unfollow <feed_url>
```

---

### List feeds

```bash
gator feeds
```

shows all feeds being followed with the Aggregator in the form of:
- ID:
- Created:
- Updated:
- Name:
- URL:
- User:
- LastFetchedAt:
- ==============

---

### View followed feeds

```bash
gator following
```

Shows all the names of the feeds the user is currently following.

---

## Aggregation

### Start scraping feeds

```bash
gator agg 10s
```

This will:

* Continuously fetch feeds
* Store posts in the database
* Run every `10s` (you can use `5s`, `1m`, etc.)

---

## Browse Posts

```bash
gator browse [limit]
```

Examples:

```bash
gator browse        # default 2 posts
gator browse 25     # show 25 posts
```

Displays posts from feeds you follow, newest first.

---

## Reset Database

```bash
gator reset
```

Deletes all users (and associated data via cascade).

---

## Notes

* Duplicate posts are automatically ignored
* Feeds are fetched in a round-robin fashion based on last fetch time
* Some feeds may have inconsistent date formats — the parser handles common cases

---

## Example Workflow

```bash
gator register jesbr
gator addfeed "Boot.dev Blog" "https://blog.boot.dev/index.xml"
gator agg 1m
gator browse
```

---

