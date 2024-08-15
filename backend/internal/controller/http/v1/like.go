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

type likeRoutes struct {
	u  *usecases.LikeUseCase
	su *usecases.SessionUseCase
}

func NewLikeRoutes(handler *gin.RouterGroup, u *usecases.LikeUseCase, su *usecases.SessionUseCase) {
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

	sess, err := r.su.SessionFromContext(c)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	like := &entities.Like{
		UserID:   sess.UserID,
		RecipeID: recipeID,
	}

	err = r.u.Like(c.Request.Context(), like)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError})
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

	sess, err := r.su.SessionFromContext(c)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	like := &entities.Like{
		UserID:   sess.UserID,
		RecipeID: recipeID,
	}

	err = r.u.Unlike(c.Request.Context(), like)
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
