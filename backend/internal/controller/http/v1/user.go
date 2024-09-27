package v1

import (
	"errors"
	"io"
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
	u  *usecases.UserUseCase
	su *usecases.SessionUseCase
}

func NewUserRoutes(handler *gin.RouterGroup, u *usecases.UserUseCase, su *usecases.SessionUseCase) {
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
		usr.GET("/:login/icon", r.getDBIcon)
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
	sess, err := r.su.SessionFromContext(c)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	token, err := r.u.GenerateJWT(sess.UserID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
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
// @Success     200 {object} entities.AuthUser
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

	authInfo, err := r.u.Signup(c.Request.Context(), user)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrUserUnique) {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecases.ErrUserUnique.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, authInfo)
}

// @Summary     Sign in
// @Description Sign in user
// @ID          signin
// @Tags  	    auth
// @Param 		user body entities.UserLogin  true  "User params"
// @Accept      json
// @Produce     json
// @Success     200 {object} entities.AuthUser
// @Failure     400
// @Failure     500
// @Router      /auth/signin [post]
func (r *userRoutes) signin(c *gin.Context) {
	var params *entities.UserLogin
	if err := c.ShouldBindJSON(&params); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	if params.Login == "" && params.Email == "" {
		errMes := "login or email must be provide"
		slog.Error(errMes)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMes})
		return
	}

	authInfo, err := r.u.Signin(c.Request.Context(), params)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, authInfo)
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
	err := r.u.Logout(c)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "you are logout"})
}

func (r *userRoutes) getIcon(c *gin.Context) (io.ReadSeeker, error) {
	err := c.Request.ParseMultipartForm(maxFilesSize)
	if err != nil {
		return nil, common.ErrHudgeFiles
	}

	fileHeader, err := c.FormFile("icon")

	if fileHeader == nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if !strings.Contains(fileHeader.Header.Get("Content-Type"), "image") {
		return nil, common.ErrImageType
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return file, nil
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

	login, ok := c.Params.Get("login")
	if !ok {
		slog.Error(common.ErrLoginProvided.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrLoginProvided.Error()})
		return
	}

	params := &entities.UserUpdate{}
	if err := c.Bind(&params); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	var err error
	params.Icon, err = r.getIcon(c)
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

	sess, err := r.su.SessionFromContext(c)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	login, err = r.u.Update(c, login, sess.UserID, params)
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
		if errors.Is(err, common.ErrNoPermissions) {
			c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrNoPermissions.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
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
	login, ok := c.Params.Get("login")
	if !ok {
		slog.Error(common.ErrLoginProvided.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrLoginProvided.Error()})
		return
	}

	params := &entities.UserPasswordUpdate{}
	if err := c.BindJSON(&params); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.GetErrMessages(err).Error()})
		return
	}

	sess, err := r.su.SessionFromContext(c)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	err = r.u.UpdatePassword(c.Request.Context(), login, sess.UserID, params)
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
	login, ok := c.Params.Get("login")
	if !ok {
		slog.Error(common.ErrLoginProvided.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrLoginProvided.Error()})
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

	user, err := r.u.Get(c, login, userID, authorized)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": usecases.ErrUserNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, entities.JSONUserInfo{User: user})
}

// @Summary     Get icon
// @Description Get user icon
// @ID          get icon
// @Tags  	    user
// @Produce     json
// @Success     200 {object} entities.UserIcon
// @Failure     400
// @Failure     401
// @Failure     404
// @Failure     500
// @Router      /user/{login}/icon [get]
func (r *userRoutes) getDBIcon(c *gin.Context) {
	login, ok := c.Params.Get("login")
	if !ok {
		slog.Error(common.ErrLoginProvided.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": common.ErrLoginProvided.Error()})
		return
	}

	icn, err := r.u.GetIcon(c.Request.Context(), login)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, usecases.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": usecases.ErrUserNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": common.ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, icn)
}
