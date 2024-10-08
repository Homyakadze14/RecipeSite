package usecases

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Homyakadze14/RecipeSite/internal/common"
	"github.com/gin-gonic/gin"
)

func (u *SessionUseCase) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, err := u.GetSession(c.Request)
		if err != nil {
			slog.Error(err.Error())
			if errors.Is(err, ErrUnauth) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrUnauth.Error()})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
			}
			return
		}
		c.Set(sessionKey, sess)
		c.Next()
	}
}
