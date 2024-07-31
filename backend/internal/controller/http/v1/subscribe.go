package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Homyakadze14/RecipeSite/internal/common"
	"github.com/Homyakadze14/RecipeSite/internal/entities"
	"github.com/Homyakadze14/RecipeSite/internal/usecases"
	"github.com/gin-gonic/gin"
)

type subscribeRoutes struct {
	u  *usecases.SubscribeUseCases
	su *usecases.SessionUseCase
}

func NewSubscribeRoutes(handler *gin.RouterGroup, u *usecases.SubscribeUseCases, su *usecases.SessionUseCase) {
	r := &subscribeRoutes{u, su}

	sb := handler.Group("/user")
	{
		sb.Use(su.Auth())
		sb.POST("/subscribe", r.subscribe)
		sb.POST("/unsubscribe", r.unsubscribe)
	}
}

// @Summary     Subscribe to user
// @Description Subscribe to user
// @ID          subscribe to user
// @Tags  	    subscription
// @Produce     json
// @Param 		creator body entities.SubscribeCreator  true  "User id to whom we subscribe"
// @Produce     json
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /user/subscribe [post]
func (r *subscribeRoutes) subscribe(c *gin.Context) {
	var creator *entities.SubscribeCreator
	if err := c.ShouldBindJSON(&creator); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	err := r.u.Subscribe(c.Request.Context(), creator, c.Request)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrUserNotFound.Error()})
			return
		}
		if errors.Is(err, usecases.ErrSubscribe) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrSubscribe.Error()})
			return
		}
		if errors.Is(err, usecases.ErrYourselfSubscribe) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrYourselfSubscribe.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "you subscribe to this user"})
}

// @Summary     Unsubscribe from user
// @Description Unsubscribe from user
// @ID          unsubscribe from user
// @Tags  	    subscription
// @Produce     json
// @Param 		creator body entities.SubscribeCreator  true  "User id to whom we unsubscribe"
// @Accept      json
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /user/unsubscribe [post]
func (r *subscribeRoutes) unsubscribe(c *gin.Context) {
	var creator *entities.SubscribeCreator
	if err := c.ShouldBindJSON(&creator); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	err := r.u.Unsubscribe(c.Request.Context(), creator, c.Request)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrUserNotFound.Error()})
			return
		}
		if errors.Is(err, usecases.ErrUnsubscribe) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrUnsubscribe.Error()})
			return
		}
		if errors.Is(err, usecases.ErrYourselfUnsubscribe) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrYourselfUnsubscribe.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "you unsubscribe to this user"})
}
