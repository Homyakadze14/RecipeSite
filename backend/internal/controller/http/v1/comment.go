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
	u  *usecases.CommentUseCases
	us *usecases.UserUseCases
	su *usecases.SessionUseCase
}

func NewCommentRoutes(handler *gin.RouterGroup, u *usecases.CommentUseCases, su *usecases.SessionUseCase, us *usecases.UserUseCases) {
	r := &commentRoutes{u, us, su}

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
	// Get recipe id
	strRecipeID, ok := c.Params.Get("id")
	if !ok {
		errRecipeID := "ID must be provided"
		slog.Error(errRecipeID)
		c.JSON(http.StatusBadRequest, gin.H{"error": errRecipeID})
		return
	}

	recipeID, err := strconv.Atoi(strRecipeID)
	if err != nil {
		errRecipeID := "ID must be integer"
		slog.Error(errRecipeID)
		c.JSON(http.StatusBadRequest, gin.H{"error": errRecipeID})
		return
	}

	// Parse json values to comment
	crComment := &entities.CommentCreate{}
	if err := c.BindJSON(&crComment); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	comment := &entities.Comment{Text: crComment.Text}

	// Get user id from session
	sess, err := r.su.GetSession(c.Request)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	comment.UserID = sess.UserID

	// Set recipe id
	comment.RecipeID = recipeID

	// Save comment
	err = r.u.Save(c.Request.Context(), comment)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
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
	// Parse json values to comment
	comment := &entities.CommentUpdate{}
	if err := c.BindJSON(&comment); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	// Update comment
	err := r.u.Update(c.Request.Context(), c.Request, comment)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrCommentNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrCommentNotFound.Error()})
			return
		}
		if errors.Is(err, usecases.ErrUserNoPermisions) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrUserNoPermisions.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
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
	// Parse json values to comment
	comment := &entities.CommentDelete{}
	if err := c.BindJSON(&comment); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	// Delete comment
	err := r.u.Delete(c.Request.Context(), c.Request, comment)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrCommentNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrCommentNotFound.Error()})
			return
		}
		if errors.Is(err, usecases.ErrUserNoPermisions) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrUserNoPermisions.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "comment deleted"})
}
