package v1

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/Homyakadze14/RecipeSite/internal/common"
	"github.com/Homyakadze14/RecipeSite/internal/entities"
	"github.com/Homyakadze14/RecipeSite/internal/usecases"
	"github.com/gin-gonic/gin"
)

type recipeRoutes struct {
	u  *usecases.RecipeUseCases
	su *usecases.SessionUseCase
}

func NewRecipeRoutes(handler *gin.RouterGroup, u *usecases.RecipeUseCases, su *usecases.SessionUseCase) {
	r := &recipeRoutes{u, su}

	h := handler.Group("/recipe")
	{
		h.GET("", r.getAll)
		h.POST("", r.getFiltered)
		h.GET("/:id", r.get)
	}

	ur := handler.Group("/user/:login/recipe")
	{
		ur.Use(su.Auth())
		ur.POST("", r.create)
		ur.PUT("/:id", r.update)
		ur.DELETE("/:id", r.delete)
	}
}

// @Summary     Get all recipe
// @Description Get all recipe
// @ID          get all recipe
// @Tags  	    recipe
// @Produce     json
// @Success     200 {object} []entities.Recipe
// @Failure     500
// @Router      /recipe [get]
func (r *recipeRoutes) getAll(c *gin.Context) {
	recipes, err := r.u.GetAll(c.Request.Context())
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recipes": recipes})
}

// @Summary     Get filtered recipe
// @Description Get filtered recipe
// @ID          get filtered recipe
// @Tags  	    recipe
// @Accept      json
// @Param 		filter body entities.RecipeFilter false "filter"
// @Produce     json
// @Success     200 {object} []entities.Recipe
// @Failure     400
// @Failure     500
// @Router      /recipe [post]
func (r *recipeRoutes) getFiltered(c *gin.Context) {
	var filter *entities.RecipeFilter
	if err := c.ShouldBindJSON(&filter); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	recipes, err := r.u.GetFiltered(c.Request.Context(), filter)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrBadOrderField) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrBadOrderField.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recipes": recipes})
}

// @Summary     Get recipe
// @Description Get recipe
// @ID          get recipe
// @Tags  	    recipe
// @Produce     json
// @Success     200 {object} entities.RecipeInfo
// @Failure     400
// @Failure     404
// @Failure     500
// @Router      /recipe/{id} [get]
func (r *recipeRoutes) get(c *gin.Context) {
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

	// Get recipe
	recipes, err := r.u.Get(c.Request.Context(), c.Request, recipeID)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrRecipeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": usecases.ErrRecipeNotFound.Error()})
			return
		}
		if errors.Is(err, usecases.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrUserNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, entities.RecipeInfo{Info: recipes})
}

// @Summary     Create recipe
// @Description Create recipe
// @ID          create recipe
// @Tags  	    recipe
// @Param 		photos formData file false "Photos"
// @Param 		recipe formData entities.CreateRecipe false "Recipe params"
// @Accept      mpfd
// @Produce     json
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /user/{login}/recipe [post]
func (r *recipeRoutes) create(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		slog.Error(ErrContentType.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrContentType.Error()})
		return
	}

	// Get user login
	login, ok := c.Params.Get("login")
	if !ok {
		errLogin := "Login must be provided in url"
		slog.Error(errLogin)
		c.JSON(http.StatusBadRequest, gin.H{"error": errLogin})
		return
	}

	// Check file size
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is too huge"})
		return
	}

	// Parse form values to recipe
	recipe := &entities.CreateRecipe{}
	if err := c.Bind(&recipe); err != nil {
		slog.Error(err.Error())
		if errors.Is(err, strconv.ErrSyntax) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrComplexityMustBeInt.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	// Create
	err = r.u.Create(c, login, recipe)
	if err != nil {
		slog.Error(err.Error())
		if strings.Contains(err.Error(), "RMQ") {
			// SKIP error
		}
		if errors.Is(err, usecases.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": usecases.ErrUserNotFound.Error()})
			return
		}
		if errors.Is(err, usecases.ErrNoPermissions) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrNoPermissions.Error()})
			return
		}
		if errors.Is(err, usecases.ErrEmptyPhotos) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrEmptyPhotos.Error()})
			return
		}
		if errors.Is(err, usecases.ErrUserNotImage) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrUserNotImage.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "recipe created"})
}

// @Summary     Update recipe
// @Description Update recipe
// @ID          update recipe
// @Tags  	    recipe
// @Param 		photos formData file false "Photos"
// @Param 		recipe formData entities.UpdateRecipe false "Recipe params"
// @Accept      mpfd
// @Produce     json
// @Success     200
// @Failure     400
// @Failure     401
// @Failure    	404
// @Failure     500
// @Router      /user/{login}/recipe/{id} [put]
func (r *recipeRoutes) update(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		slog.Error(ErrContentType.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrContentType.Error()})
		return
	}

	// Get user login
	login, ok := c.Params.Get("login")
	if !ok {
		errLogin := "Login must be provided in url"
		slog.Error(errLogin)
		c.JSON(http.StatusBadRequest, gin.H{"error": errLogin})
		return
	}

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

	// Check file size
	err = c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is too huge"})
		return
	}

	// Parse form values to recipe
	recipe := &entities.UpdateRecipe{}
	if err := c.Bind(&recipe); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	// Update
	err = r.u.Update(c, login, recipeID, recipe)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": usecases.ErrUserNotFound.Error()})
			return
		}
		if errors.Is(err, usecases.ErrRecipeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": usecases.ErrRecipeNotFound.Error()})
			return
		}
		if errors.Is(err, usecases.ErrNoPermissions) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrNoPermissions.Error()})
			return
		}
		if errors.Is(err, usecases.ErrUserNotImage) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrUserNotImage.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "recipe updated"})
}

// @Summary     Delete recipe
// @Description Delete recipe
// @ID          delete recipe
// @Tags  	    recipe
// @Produce     json
// @Success     200
// @Failure     400
// @Failure     401
// @Failure    	404
// @Failure     500
// @Router      /user/{login}/recipe/{id} [delete]
func (r *recipeRoutes) delete(c *gin.Context) {
	// Get user login
	login, ok := c.Params.Get("login")
	if !ok {
		errLogin := "Login must be provided in url"
		slog.Error(errLogin)
		c.JSON(http.StatusBadRequest, gin.H{"error": errLogin})
		return
	}

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

	// Delete
	err = r.u.Delete(c, login, recipeID)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": usecases.ErrUserNotFound.Error()})
			return
		}
		if errors.Is(err, usecases.ErrRecipeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": usecases.ErrRecipeNotFound.Error()})
			return
		}
		if errors.Is(err, usecases.ErrNoPermissions) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrNoPermissions.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "recipe deleted"})
}
