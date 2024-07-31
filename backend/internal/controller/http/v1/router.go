// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	_ "github.com/Homyakadze14/RecipeSite/docs"
	"github.com/Homyakadze14/RecipeSite/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter -.
// Swagger spec:
// @title       RecipeSite
// @description RestAPI for recipe site
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine,
	sess *usecases.SessionUseCase,
	user *usecases.UserUseCases,
	like *usecases.LikeUseCases,
	recipe *usecases.RecipeUseCases,
	comment *usecases.CommentUseCases,
	subscribe *usecases.SubscribeUseCases) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/v1")
	{
		NewUserRoutes(h, user, sess)
		NewLikeRoutes(h, like, sess)
		NewRecipeRoutes(h, recipe, sess)
		NewCommentRoutes(h, comment, sess, user)
		NewSubscribeRoutes(h, subscribe, sess)
	}
}
