// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package db

import (
	"context"
)

type Querier interface {
	CreatePost(ctx context.Context, arg CreatePostParams) (Post, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeletePost(ctx context.Context, id int64) error
	GetPost(ctx context.Context, id int64) (Post, error)
	GetUser(ctx context.Context, userName string) (User, error)
	ListPosts(ctx context.Context, arg ListPostsParams) ([]Post, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdatePost(ctx context.Context, arg UpdatePostParams) (Post, error)
}

var _ Querier = (*Queries)(nil)
