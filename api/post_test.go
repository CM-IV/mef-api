package api

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/CM-IV/mef-api/db/mock"
	db "github.com/CM-IV/mef-api/db/sqlc"
	"github.com/CM-IV/mef-api/util"
	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func TestGetPostAPI(t *testing.T) {

	post := randomPost()

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
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/posts/%d", tc.postID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)

		})

	}

}

func randomPost() db.Post {

	return db.Post{

		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Image:    util.RandomImage(),
		Title:    util.RandomTitle(),
		Subtitle: util.RandomSubtitle(),
		Content:  util.RandomContent(),
	}

}

func requireBodyMatchPost(t *testing.T, body *bytes.Buffer, post db.Post) {

	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotPost db.Post

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	err = json.Unmarshal(data, &gotPost)
	require.NoError(t, err)
	require.Equal(t, post, gotPost)

}
