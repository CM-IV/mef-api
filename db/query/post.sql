-- name: CreatePost :one
INSERT INTO posts (
  owner,
  image,
  title,
  subtitle,
  content
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;


-- name: GetPost :one
SELECT * FROM posts
WHERE id = $1 LIMIT 1;

-- name: CountPosts :one
SELECT COUNT(*) as total_posts FROM posts;

-- name: ListPosts :many
SELECT * FROM posts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdatePost :one
UPDATE posts
SET content = $2
WHERE id = $1
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;
