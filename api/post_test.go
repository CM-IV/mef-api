package api

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/CM-IV/mef-api/db/mock"
	db "github.com/CM-IV/mef-api/db/sqlc"
	"github.com/CM-IV/mef-api/token"
	"github.com/CM-IV/mef-api/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func TestGetPostAPI(t *testing.T) {
	user, _ := randomUser(t)
	post := randomPost(user.UserName)

	testCases := []struct {
		title         string
		postID        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{

		{

			title:  "OK",
			postID: post.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPost(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return(post, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {

				//check response
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPost(t, recorder.Body, post)

			},
		},
		{

			title:  "NotFound",
			postID: post.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPost(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return(db.Post{}, sql.ErrNoRows)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusNotFound, recorder.Code)

			},
		},
		{

			title:  "InternalError",
			postID: post.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPost(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return(db.Post{}, sql.ErrConnDone)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
		{

			title:  "InvalidID",
			postID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPost(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
	}

	for i := range testCases {

		tc := testCases[i]

		t.Run(tc.title, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			//start test http server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/posts/%d", tc.postID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)

		})

	}

}

func TestCreatePostAPI(t *testing.T) {
	user, _ := randomUser(t)
	post := randomPost(user.UserName)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"image":    post.Image,
				"title":    post.Title,
				"subtitle": post.Subtitle,
				"content":  post.Content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreatePostParams{
					Owner:    post.Owner,
					Image:    post.Image,
					Title:    post.Title,
					Subtitle: post.Subtitle,
					Content:  post.Content,
				}

				store.EXPECT().
					CreatePost(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(post, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchPost(t, recorder.Body, post)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"owner":    post.Owner,
				"image":    post.Image,
				"title":    post.Title,
				"subtitle": post.Subtitle,
				"content":  post.Content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					CreatePost(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				requireBodyMatchPost(t, recorder.Body, post)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"owner":    post.Owner,
				"image":    post.Image,
				"title":    post.Title,
				"subtitle": post.Subtitle,
				"content":  post.Content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreatePost(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Post{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidTitle",
			body: gin.H{
				"owner":    post.Owner,
				"image":    post.Image,
				"title":    "",
				"subtitle": post.Subtitle,
				"content":  post.Content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreatePost(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidOwner",
			body: gin.H{
				"owner":    "",
				"image":    post.Image,
				"title":    post.Title,
				"subtitle": post.Subtitle,
				"content":  post.Content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreatePost(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			json := jsoniter.ConfigCompatibleWithStandardLibrary

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/posts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListPostsAPI(t *testing.T) {
	user, _ := randomUser(t)

	n := 5
	posts := make([]db.Post, n)
	for i := 0; i < n; i++ {
		posts[i] = randomPost(user.UserName)
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListPostsParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListPosts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(posts, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPosts(t, recorder.Body, posts)
			},
		},
		{
			name: "InternalError",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListPosts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Post{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListPosts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				pageID:   1,
				pageSize: 100000,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListPosts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/api/posts"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomPost(owner string) db.Post {

	return db.Post{

		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Image:    util.RandomImage(),
		Title:    util.RandomTitle(),
		Subtitle: util.RandomSubtitle(),
		Content:  util.RandomContent(),
	}

}

func requireBodyMatchPost(t *testing.T, body *bytes.Buffer, post db.Post) {

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	var gotPost db.Post

	err := json.NewDecoder(body).Decode(&gotPost)
	require.NoError(t, err)
	require.Equal(t, post, gotPost)

}

func requireBodyMatchPosts(t *testing.T, body *bytes.Buffer, posts []db.Post) {

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	var gotPosts []db.Post

	err := json.NewDecoder(body).Decode(&gotPosts)
	require.NoError(t, err)
	require.Equal(t, posts, gotPosts)

}
