-- name: DeleteFeedFollow :one
DELETE FROM feed_follows
USING users, feeds
WHERE feed_follows.user_id = users.id
  AND feed_follows.feed_id = feeds.id
  AND users.name = $1
  AND feeds.url = $2
RETURNING feed_follows.*;