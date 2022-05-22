package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/CM-IV/mef-api/util"

	"github.com/stretchr/testify/require"
)

func createRandomPost(t *testing.T) Post {
	user := createRandomUser(t)
	arg := CreatePostParams{

		Owner:    user.UserName,
		Image:    util.RandomImage(),
		Title:    util.RandomTitle(),
		Subtitle: util.RandomSubtitle(),
		Content:  util.RandomContent(),
	}

	post, err := testQueries.CreatePost(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, post)

	require.Equal(t, arg.Owner, post.Owner)
	require.Equal(t, arg.Image, post.Image)
	require.Equal(t, arg.Title, post.Title)
	require.Equal(t, arg.Subtitle, post.Subtitle)
	require.Equal(t, arg.Content, post.Content)

	require.NotZero(t, post.ID)
	require.NotZero(t, post.CreatedAt)

	return post

}

func TestCreatePost(t *testing.T) {

	createRandomPost(t)

}

func TestGetPost(t *testing.T) {

	post1 := createRandomPost(t)
	post2, err := testQueries.GetPost(context.Background(), post1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, post2)

	require.Equal(t, post1.ID, post2.ID)
	require.Equal(t, post1.Owner, post2.Owner)
	require.Equal(t, post1.Image, post2.Image)
	require.Equal(t, post1.Title, post2.Title)
	require.Equal(t, post1.Subtitle, post2.Subtitle)
	require.Equal(t, post1.Content, post2.Content)
	require.WithinDuration(t, post1.CreatedAt, post2.CreatedAt, time.Second)

}

func TestUpdatePost(t *testing.T) {

	post1 := createRandomPost(t)

	args := UpdatePostParams{

		ID:      post1.ID,
		Content: util.RandomContent(),
	}

	post2, err := testQueries.UpdatePost(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, post2)

	require.Equal(t, post1.ID, post2.ID)
	require.Equal(t, post1.Owner, post2.Owner)
	require.Equal(t, post1.Image, post2.Image)
	require.Equal(t, post1.Title, post2.Title)
	require.Equal(t, post1.Subtitle, post2.Subtitle)
	require.Equal(t, args.Content, post2.Content)
	require.WithinDuration(t, post1.CreatedAt, post2.CreatedAt, time.Second)

}

func TestDeletePost(t *testing.T) {

	post1 := createRandomPost(t)
	err := testQueries.DeletePost(context.Background(), post1.ID)
	require.NoError(t, err)

	post2, err := testQueries.GetPost(context.Background(), post1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, post2)

}

func TestListPosts(t *testing.T) {

	for i := 0; i < 10; i++ {

		createRandomPost(t)

	}

	args := ListPostsParams{

		Limit:  5,
		Offset: 5,
	}

	posts, lastPage, totalRecords, err := testQueries.ListPosts(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, posts, 5)

	for _, post := range posts {

		require.NotEmpty(t, post)
		require.NotNil(t, lastPage)
		require.NotNil(t, totalRecords)

	}

}
