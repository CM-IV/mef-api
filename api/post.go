package api

import (
	"database/sql"
	"math"
	"net/http"

	db "github.com/CM-IV/mef-api/db/sqlc"
	"github.com/CM-IV/mef-api/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createPostRequest struct {
	Image    string `json:"image" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Subtitle string `json:"subtitle" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type getPostRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type deletePostRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type listPostRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=15"`
}

type updatePostRequest struct {
	Content string `json:"content" binding:"required"`
}

type updatePostRequestID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) createPost(ctx *gin.Context) {

	var req createPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return

	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreatePostParams{

		Owner:    authPayload.UserName,
		Image:    req.Image,
		Title:    req.Title,
		Subtitle: req.Subtitle,
		Content:  req.Content,
	}

	post, err := server.store.CreatePost(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}

	ctx.JSON(http.StatusCreated, post)

}

func (server *Server) getPost(ctx *gin.Context) {

	var req getPostRequest
	if err := ctx.ShouldBindUri(&req); err != nil {

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return

	}

	post, err := server.store.GetPost(ctx, req.ID)

	if err != nil {

		if err == sql.ErrNoRows {

			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return

		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, post)

}

func (server *Server) listPost(ctx *gin.Context) {

	var req listPostRequest
	var lastPage int64

	if err := ctx.ShouldBindQuery(&req); err != nil {

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return

	}

	args := db.ListPostsParams{

		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	totalRecords, err := server.store.CountPosts(ctx)

	if err != nil {

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	posts, err := server.store.ListPosts(ctx, args)

	if err != nil {

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var resp struct {
		TotalRecords int64     `json:"total_records"`
		LastPage     int64     `json:"last_page"`
		Posts        []db.Post `json:"posts"`
	}

	lastPage = int64(math.Ceil(float64(totalRecords) / float64(args.Limit)))

	resp.TotalRecords = totalRecords
	resp.LastPage = lastPage
	resp.Posts = posts

	ctx.JSON(http.StatusOK, resp)

}

func (server *Server) updatePost(ctx *gin.Context) {

	var req updatePostRequest
	var id updatePostRequestID

	if err_json := ctx.ShouldBindJSON(&req); err_json != nil {

		ctx.JSON(http.StatusBadRequest, errorResponse(err_json))
		return

	}

	if err := ctx.ShouldBindUri(&id); err != nil {

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return

	}

	args := db.UpdatePostParams{

		ID:      id.ID,
		Content: req.Content,
	}

	post, err_update := server.store.UpdatePost(ctx, args)

	if err_update != nil {

		if err_update == sql.ErrNoRows {

			ctx.JSON(http.StatusNotFound, errorResponse(err_update))
			return

		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err_update))
		return
	}

	ctx.JSON(http.StatusOK, post)

}

func (server *Server) deletePost(ctx *gin.Context) {

	var req deletePostRequest

	if err := ctx.ShouldBindUri(&req); err != nil {

		ctx.JSON(http.StatusBadRequest, errorResponse(err))

	}

	del_err := server.store.DeletePost(ctx, req.ID)

	if del_err != nil {

		if del_err == sql.ErrNoRows {

			ctx.JSON(http.StatusNotFound, errorResponse(del_err))

		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(del_err))
	}

	ctx.JSON(http.StatusOK, del_err)

}
