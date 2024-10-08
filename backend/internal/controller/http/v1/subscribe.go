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

	sb := handler.Group("/user/:login")
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
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /user/{login}/subscribe [post]
func (r *subscribeRoutes) subscribe(c *gin.Context) {
	login, ok := c.Params.Get("login")
	if !ok {
		slog.Error(common.ErrLoginProvided.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrLoginProvided.Error()})
		return
	}

	sess, err := r.su.SessionFromContext(c)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	info := &entities.SubscribeInfo{
		CreatorLogin: login,
		SubscriberID: sess.UserID,
	}

	err = r.u.Subscribe(c.Request.Context(), info)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "you subscribe to this user"})
}

// @Summary     Unsubscribe from user
// @Description Unsubscribe from user
// @ID          unsubscribe from user
// @Tags  	    subscription
// @Produce     json
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /user/{login}/unsubscribe [post]
func (r *subscribeRoutes) unsubscribe(c *gin.Context) {
	login, ok := c.Params.Get("login")
	if !ok {
		slog.Error(common.ErrLoginProvided.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrLoginProvided.Error()})
		return
	}

	sess, err := r.su.SessionFromContext(c)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	info := &entities.SubscribeInfo{
		CreatorLogin: login,
		SubscriberID: sess.UserID,
	}

	err = r.u.Unsubscribe(c.Request.Context(), info)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "you unsubscribe to this user"})
}
