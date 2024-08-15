package v1

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Homyakadze14/RecipeSite/internal/common"
	"github.com/Homyakadze14/RecipeSite/internal/entities"
	"github.com/Homyakadze14/RecipeSite/internal/usecases"
	"github.com/gin-gonic/gin"
)

var (
	ErrContentType = errors.New("content type must be multipart/form-data")
)

type userRoutes struct {
	u  *usecases.UserUseCases
	su *usecases.SessionUseCase
}

func NewUserRoutes(handler *gin.RouterGroup, u *usecases.UserUseCases, su *usecases.SessionUseCase) {
	r := &userRoutes{u, su}

	h := handler.Group("/auth")
	{
		h.POST("/signup", r.signup)
		h.POST("/signin", r.signin)
		h.POST("/checktgtoken", r.checkTGToken)
	}

	a := handler.Group("/auth")
	{
		a.Use(su.Auth())
		a.POST("/logout", r.logout)
		a.GET("/tgtoken", r.tgToken)
	}

	usr := handler.Group("/user")
	{
		usr.Use(su.Auth())
		usr.PUT("/:login", r.update)
		usr.PUT("/:login/password", r.updatePassword)
	}

	us := handler.Group("/user")
	{
		us.GET("/:login", r.get)
	}
}

// @Summary     Generate user telegram token
// @Description Generate user telegram token
// @ID          generate user telegram token
// @Tags  	    auth
// @Produce     json
// @Success     200 {object} entities.JWTToken
// @Failure     401
// @Failure     500
// @Router      /auth/tgtoken [get]
func (r *userRoutes) tgToken(c *gin.Context) {
	token, err := r.u.GenerateJWT(c.Request)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, token)
}

// @Summary     Check user telegram token
// @Description Check user telegram token
// @ID          Check user telegram token
// @Tags  	    auth
// @Param 		token body entities.JWTToken  true  "token"
// @Accept 		json
// @Produce     json
// @Success     200 {object} entities.JWTData
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /auth/checktgtoken [post]
func (r *userRoutes) checkTGToken(c *gin.Context) {
	var token *entities.JWTToken
	if err := c.ShouldBindJSON(&token); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	data, err := r.u.GetDataFromJWT(token)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrBadToken) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// @Summary     Sign up
// @Description Sign up user
// @ID          signup
// @Tags  	    auth
// @Param 		user body entities.UserLogin  true  "User params"
// @Accept      json
// @Produce     json
// @Success     200
// @Failure     400
// @Failure     500
// @Router      /auth/signup [post]
func (r *userRoutes) signup(c *gin.Context) {
	var user *entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	cookie, err := r.u.Signup(c.Request.Context(), user)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrUserUnique) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrUserUnique.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, gin.H{"status": "you are signed up"})
}

// @Summary     Sign in
// @Description Sign in user
// @ID          signin
// @Tags  	    auth
// @Param 		user body entities.UserLogin  true  "User params"
// @Accept      json
// @Produce     json
// @Success     200 {string} string "login"
// @Failure     400
// @Failure     500
// @Router      /auth/signin [post]
func (r *userRoutes) signin(c *gin.Context) {
	var user *entities.UserLogin
	if err := c.ShouldBindJSON(&user); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	if user.Login == "" && user.Email == "" {
		errMes := "login or email must be provide"
		slog.Error(errMes)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMes})
		return
	}

	cookie, login, err := r.u.Signin(c.Request.Context(), user)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": usecases.ErrUserNotFound.Error()})
			return
		}
		if errors.Is(err, usecases.ErrUserWrongPassword) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrUserWrongPassword.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, gin.H{"login": login})
}

// @Summary     Logout
// @Description Logout user
// @ID          Logout
// @Tags  	    auth
// @Produce     json
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /auth/logout [post]
func (r *userRoutes) logout(c *gin.Context) {
	cookie, err := r.u.Logout(c)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, gin.H{"status": "you are logout"})
}

// @Summary     Update user
// @Description Update user
// @ID          update user
// @Tags  	    user
// @Param 		icon formData file false "Icon"
// @Param 		user formData entities.UserUpdate false "User params"
// @Accept      mpfd
// @Produce     json
// @Success     200 {string} string "login"
// @Failure     400
// @Failure     404
// @Failure     401
// @Failure     500
// @Router      /user/{login} [put]
func (r *userRoutes) update(c *gin.Context) {
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

	// Parse form values to user
	user := &entities.UserUpdate{}
	if err := c.Bind(&user); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	// Update
	login, err = r.u.Update(c, login, user, c.Request)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": usecases.ErrUserNotFound.Error()})
			return
		}
		if errors.Is(err, usecases.ErrUserUnique) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrUserUnique.Error()})
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

	c.JSON(http.StatusOK, gin.H{"login": login})
}

// @Summary     Update user password
// @Description Update user password
// @ID          update user password
// @Tags  	    user
// @Param 		user body entities.UserPasswordUpdate false "User params"
// @Accept      json
// @Produce     json
// @Success     200
// @Failure     400
// @Failure     404
// @Failure     401
// @Failure     500
// @Router      /user/{login}/password [put]
func (r *userRoutes) updatePassword(c *gin.Context) {
	// Get user login
	login, ok := c.Params.Get("login")
	if !ok {
		errLogin := "Login must be provided in url"
		slog.Error(errLogin)
		c.JSON(http.StatusBadRequest, gin.H{"error": errLogin})
		return
	}

	// Parse form values to user
	user := &entities.UserPasswordUpdate{}
	if err := c.BindJSON(&user); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	// Update
	err := r.u.UpdatePassword(c.Request.Context(), login, user, c.Request)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": usecases.ErrUserNotFound.Error()})
			return
		}
		if errors.Is(err, usecases.ErrNoPermissions) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrNoPermissions.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "user password updated"})
}

// @Summary     Get user info
// @Description Get user info
// @ID          get user info
// @Tags  	    user
// @Produce     json
// @Success     200 {object} entities.JSONUserInfo
// @Failure     404
// @Failure     500
// @Router      /user/{login} [get]
func (r *userRoutes) get(c *gin.Context) {
	// Get user login
	login, ok := c.Params.Get("login")
	if !ok {
		errLogin := "Login must be provided in url"
		slog.Error(errLogin)
		c.JSON(http.StatusBadRequest, gin.H{"error": errLogin})
		return
	}

	// Get user
	user, err := r.u.Get(c, login)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": usecases.ErrUserNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, entities.JSONUserInfo{User: user})
}
