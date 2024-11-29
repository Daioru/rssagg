-- name: CreatePost :one
INSERT INTO posts (
    id, 
    created_at, 
    updated_at, 
    title, 
    description,
    published_at,
    url,
    feed_id
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.* from posts
join feed_follows on posts.feed_id = feed_follows.feed_id
where feed_follows.user_id = $1
order by posts.published_at desc
limit $2;