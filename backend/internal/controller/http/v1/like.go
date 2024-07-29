package v1

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Homyakadze14/RecipeSite/internal/usecases"
	"github.com/gin-gonic/gin"
)

type likeRoutes struct {
	u  *usecases.LikeUseCases
	su *usecases.SessionUseCase
}

func NewLikeRoutes(handler *gin.RouterGroup, u *usecases.LikeUseCases, su *usecases.SessionUseCase) {
	r := &likeRoutes{u, su}

	h := handler.Group("/recipe/:id")
	{
		h.Use(su.Auth())
		h.POST("/like", r.like)
		h.POST("/unlike", r.unlike)
	}
}

// @Summary     Like
// @Description Like recipe
// @ID          like
// @Tags  	    likes
// @Produce     json
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /recipe/{id}/like [post]
func (r *likeRoutes) like(c *gin.Context) {
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

	// Like recipe
	err = r.u.Like(c.Request.Context(), c.Request, recipeID)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrAlreadyLike) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrAlreadyLike.Error()})
			return
		}
		if errors.Is(err, usecases.ErrRecipeNotExist) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrRecipeNotExist.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "recipe liked"})
}

// @Summary     Unlike
// @Description Unlike recipe
// @ID          unlike
// @Tags  	    likes
// @Produce     json
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /recipe/{id}/unlike [post]
func (r *likeRoutes) unlike(c *gin.Context) {
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

	// Like recipe
	err = r.u.Unlike(c.Request.Context(), c.Request, recipeID)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrNotLikedYet) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrNotLikedYet.Error()})
			return
		}
		if errors.Is(err, usecases.ErrRecipeNotExist) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrRecipeNotExist.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "recipe unliked"})
}
