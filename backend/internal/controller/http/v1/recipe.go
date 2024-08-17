package v1

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/Homyakadze14/RecipeSite/internal/common"
	"github.com/Homyakadze14/RecipeSite/internal/entities"
	"github.com/Homyakadze14/RecipeSite/internal/usecases"
	"github.com/gin-gonic/gin"
)

const (
	maxFilesSize    = 10 << 20
	photosArrLenght = 5
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError})
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

	sess, err := r.su.GetSession(c.Request)
	authorized := true
	userID := 0
	if err != nil {
		if !errors.Is(err, usecases.ErrUnauth) {
			slog.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
			return
		}
		authorized = false
	}

	if authorized {
		userID = sess.UserID
	}

	recipe, err := r.u.Get(c.Request.Context(), recipeID, userID, authorized)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, entities.RecipeInfo{Info: recipe})
}

func (r *recipeRoutes) getPhotos(c *gin.Context) ([]io.ReadSeeker, error) {
	err := c.Request.ParseMultipartForm(maxFilesSize)
	if err != nil {
		return nil, common.ErrHudgeFiles
	}

	multipartFormData, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}
	files := multipartFormData.File["photos"]

	photos := make([]io.ReadSeeker, 0, photosArrLenght)

	for _, fileHeader := range files {
		if !strings.Contains(fileHeader.Header.Get("Content-Type"), "image") {
			return nil, common.ErrImageType
		}

		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()
		photos = append(photos, file)
	}

	return photos, nil
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

	login, ok := c.Params.Get("login")
	if !ok {
		slog.Error(common.ErrLoginProvided.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrLoginProvided.Error()})
		return
	}

	params := &entities.CreateRecipe{}
	if err := c.Bind(&params); err != nil {
		slog.Error(err.Error())
		if errors.Is(err, strconv.ErrSyntax) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrComplexityMustBeInt.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	photos, err := r.getPhotos(c)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, common.ErrHudgeFiles) {
			c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrHudgeFiles.Error()})
			return
		}
		if errors.Is(err, common.ErrImageType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrImageType.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}
	params.Photos = photos

	sess, err := r.su.SessionFromContext(c)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	err = r.u.Create(c, login, sess.UserID, params)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": usecases.ErrUserNotFound.Error()})
			return
		}
		if errors.Is(err, common.ErrNoPermissions) {
			c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrNoPermissions.Error()})
			return
		}
		if errors.Is(err, usecases.ErrEmptyPhotos) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrEmptyPhotos.Error()})
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

	login, ok := c.Params.Get("login")
	if !ok {
		slog.Error(common.ErrLoginProvided.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrLoginProvided.Error()})
		return
	}

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

	params := &entities.UpdateRecipe{}
	if err := c.Bind(&params); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	photos, err := r.getPhotos(c)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, common.ErrHudgeFiles) {
			c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrHudgeFiles.Error()})
			return
		}
		if errors.Is(err, common.ErrImageType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrImageType.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}
	params.Photos = photos

	sess, err := r.su.SessionFromContext(c)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	err = r.u.Update(c, login, sess.UserID, recipeID, params)
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
		if errors.Is(err, common.ErrNoPermissions) {
			c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrNoPermissions.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
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
	login, ok := c.Params.Get("login")
	if !ok {
		slog.Error(common.ErrLoginProvided.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrLoginProvided.Error()})
		return
	}

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

	err = r.u.Delete(c, login, sess.UserID, recipeID)
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
		if errors.Is(err, common.ErrNoPermissions) {
			c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrNoPermissions.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "recipe deleted"})
}
