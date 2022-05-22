// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: post.sql

package db

import (
	"context"
)

const countPosts = `-- name: CountPosts :one
SELECT COUNT(*) as total_posts FROM posts
`

func (q *Queries) CountPosts(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countPosts)
	var total_posts int64
	err := row.Scan(&total_posts)
	return total_posts, err
}

const createPost = `-- name: CreatePost :one
INSERT INTO posts (
  owner,
  image,
  title,
  subtitle,
  content
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, owner, image, title, subtitle, content, created_at
`

type CreatePostParams struct {
	Owner    string `json:"owner"`
	Image    string `json:"image"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Content  string `json:"content"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.Owner,
		arg.Image,
		arg.Title,
		arg.Subtitle,
		arg.Content,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Image,
		&i.Title,
		&i.Subtitle,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}

const deletePost = `-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1
`

func (q *Queries) DeletePost(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deletePost, id)
	return err
}

const getPost = `-- name: GetPost :one
SELECT id, owner, image, title, subtitle, content, created_at FROM posts
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPost(ctx context.Context, id int64) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPost, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Image,
		&i.Title,
		&i.Subtitle,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}

const listPosts = `-- name: ListPosts :many
SELECT id, owner, image, title, subtitle, content, created_at FROM posts
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListPostsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListPosts(ctx context.Context, arg ListPostsParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, listPosts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Post{}
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Image,
			&i.Title,
			&i.Subtitle,
			&i.Content,
			&i.CreatedAt,
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

const updatePost = `-- name: UpdatePost :one
UPDATE posts
SET content = $2
WHERE id = $1
RETURNING id, owner, image, title, subtitle, content, created_at
`

type UpdatePostParams struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
}

func (q *Queries) UpdatePost(ctx context.Context, arg UpdatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, updatePost, arg.ID, arg.Content)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Image,
		&i.Title,
		&i.Subtitle,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}
