package api

import (
	"database/sql"
	"errors"
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

type listPostRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=15"`
}

type updatePostRequest struct {
	Content string `json:"content" binding:"required"`
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
	if err := ctx.ShouldBindQuery(&req); err != nil {

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return

	}

	args := db.ListPostsParams{

		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	posts, err := server.store.ListPosts(ctx, args)

	if err != nil {

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, posts)

}

func (server *Server) updatePost(ctx *gin.Context) {

	var req getPostRequest
	var update_post updatePostRequest

	if err := ctx.ShouldBindUri(&req); err != nil {

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return

	}

	postObj, err := server.store.GetPost(ctx, req.ID)

	if err != nil {

		if err == sql.ErrNoRows {

			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return

		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err_json := ctx.ShouldBindJSON(&update_post); err_json != nil {

		ctx.JSON(http.StatusBadRequest, errorResponse(err_json))
		return

	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if postObj.Owner != authPayload.UserName {
		err := errors.New("post does not belong to the authenticated user, cannot update it")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	args := db.UpdatePostParams{

		ID:      req.ID,
		Content: update_post.Content,
	}

	post, err_update := server.store.UpdatePost(ctx, args)

	if err_update != nil {

		ctx.JSON(http.StatusInternalServerError, errorResponse(err_update))
		return

	}

	ctx.JSON(http.StatusOK, post)

}

func (server *Server) deletePost(ctx *gin.Context) {

	var req getPostRequest

	if err := ctx.ShouldBindUri(&req); err != nil {

		ctx.JSON(http.StatusBadRequest, errorResponse(err))

	}

	_, err := server.store.GetPost(ctx, req.ID)

	if err != nil {

		if err == sql.ErrNoRows {

			ctx.JSON(http.StatusNotFound, errorResponse(err))

		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	del_err := server.store.DeletePost(ctx, req.ID)

	if del_err != nil {

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, del_err)

}
