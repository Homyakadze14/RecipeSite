package v1

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Homyakadze14/RecipeSite/internal/common"
	"github.com/Homyakadze14/RecipeSite/internal/entities"
	"github.com/Homyakadze14/RecipeSite/internal/usecases"
	"github.com/gin-gonic/gin"
)

type commentRoutes struct {
	u  *usecases.CommentUseCase
	su *usecases.SessionUseCase
}

func NewCommentRoutes(handler *gin.RouterGroup, u *usecases.CommentUseCase, su *usecases.SessionUseCase) {
	r := &commentRoutes{u, su}

	h := handler.Group("/recipe/:id/comment")
	{
		h.Use(su.Auth())
		h.POST("", r.create)
		h.PUT("", r.update)
		h.DELETE("", r.delete)
	}
}

// @Summary     Create comment
// @Description Create comment
// @ID          create comment
// @Tags  	    comments
// @Accept      json
// @Param 		comment body entities.CommentCreate  true  "Comment params"
// @Produce     json
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /recipe/{id}/comment [post]
func (r *commentRoutes) create(c *gin.Context) {
	urlParam, ok := c.Params.Get("id")
	if !ok {
		slog.Error(common.ErrUrlParam.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrUrlParam.Error()})
		return
	}

	recipeID, err := strconv.Atoi(urlParam)
	if err != nil {
		slog.Error(common.ErrRecipeIDType.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrRecipeIDType.Error()})
		return
	}

	response := &entities.CommentCreate{}
	if err := c.BindJSON(&response); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	comment := &entities.Comment{Text: response.Text}

	sess, err := r.su.SessionFromContext(c)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	comment.UserID = sess.UserID
	comment.RecipeID = recipeID

	err = r.u.Save(c.Request.Context(), comment)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "comment saved"})
}

// @Summary     Update comment
// @Description Update comment
// @ID          update comment
// @Tags  	    comments
// @Accept      json
// @Param 		comment body entities.CommentUpdate  true  "Comment params"
// @Produce     json
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /recipe/{id}/comment [put]
func (r *commentRoutes) update(c *gin.Context) {
	comment := &entities.CommentUpdate{}
	if err := c.BindJSON(&comment); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	sess, err := r.su.SessionFromContext(c)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	err = r.u.Update(c.Request.Context(), comment, sess.UserID)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrCommentNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrCommentNotFound.Error()})
			return
		}
		if errors.Is(err, usecases.ErrNoPermissions) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrNoPermissions.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "comment updated"})
}

// @Summary     Delete comment
// @Description Delete comment
// @ID          delete comment
// @Tags  	    comments
// @Accept      json
// @Param 		comment body entities.CommentDelete  true  "Comment params"
// @Produce     json
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /recipe/{id}/comment [delete]
func (r *commentRoutes) delete(c *gin.Context) {
	comment := &entities.CommentDelete{}
	if err := c.BindJSON(&comment); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	sess, err := r.su.SessionFromContext(c)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	err = r.u.Delete(c.Request.Context(), comment, sess.UserID)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrCommentNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrCommentNotFound.Error()})
			return
		}
		if errors.Is(err, usecases.ErrNoPermissions) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrNoPermissions.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "comment deleted"})
}
