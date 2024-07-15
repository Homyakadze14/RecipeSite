package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Homyakadze14/RecipeSite/internal/config"
	"github.com/Homyakadze14/RecipeSite/internal/images"
	"github.com/Homyakadze14/RecipeSite/internal/jsonvalidator"
	"github.com/Homyakadze14/RecipeSite/internal/models"
	"github.com/Homyakadze14/RecipeSite/internal/repos"
	"github.com/Homyakadze14/RecipeSite/internal/session"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	usrRepo        *repos.UserRepository
	likeRepo       *repos.LikeRepository
	recipeRepo     *repos.RecipeRepository
	sessionManager *session.SessionManager
	s3             *images.S3Storage
	validator      *jsonvalidator.JSONValidator
}

func NewService(ur *repos.UserRepository, sm *session.SessionManager, lr *repos.LikeRepository, rr *repos.RecipeRepository, s3 *images.S3Storage, v *jsonvalidator.JSONValidator) *UserService {
	return &UserService{
		usrRepo:        ur,
		validator:      v,
		sessionManager: sm,
		likeRepo:       lr,
		recipeRepo:     rr,
		s3:             s3,
	}
}

func (us *UserService) HandlFuncs(handler *mux.Router) {
	auth := handler.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/signup", us.signup).Methods(http.MethodPost)
	auth.HandleFunc("/signin", us.signin).Methods(http.MethodPost)

	logout := auth.PathPrefix("/logout").Subrouter()
	logout.Use(us.sessionManager.AuthMiddleware)
	logout.HandleFunc("", us.logout).Methods(http.MethodPost)

	user := handler.PathPrefix("/user").Subrouter()
	user.Use(us.sessionManager.AuthMiddleware)
	user.HandleFunc("/{login}", us.update).Methods(http.MethodPut)
	user.HandleFunc("/{login}/password", us.updatePassword).Methods(http.MethodPut)

	userWithoutAuth := handler.PathPrefix("/user").Subrouter()
	userWithoutAuth.HandleFunc("/{login}", us.get).Methods(http.MethodGet)
}

func (us *UserService) signup(w http.ResponseWriter, r *http.Request) {
	// Read request body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't parse body", http.StatusInternalServerError)
		return
	}

	// Parse json values to user
	usr := &models.User{}
	err = json.Unmarshal(data, &usr)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't parse json", http.StatusInternalServerError)
		return
	}

	// validate
	err = us.validator.Struct(usr)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate
	err = us.validator.Struct(usr)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// set default icon
	usr.Icon_URL = config.DefaultIconURL

	// Hash password
	cryptPass, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't hash password", http.StatusInternalServerError)
		return
	}
	usr.Password = string(cryptPass)

	// Save to storage
	id, err := us.usrRepo.Create(r.Context(), usr)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create session
	sess, err := us.sessionManager.Create(r.Context(), id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't create session", http.StatusInternalServerError)
		return
	}

	// Send cookie
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Expires: time.Now().Add(90 * 60 * time.Hour),
		Path:    "/",
	}

	http.SetCookie(w, cookie)
}

func (us *UserService) signin(w http.ResponseWriter, r *http.Request) {
	// Read request body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't parse body", http.StatusInternalServerError)
		return
	}

	// Parse json values to user
	usr := &models.UserLogin{}
	err = json.Unmarshal(data, &usr)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't parse json", http.StatusInternalServerError)
		return
	}

	// validate
	err = us.validator.Struct(usr)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if usr.Login == "" && usr.Email == "" {
		errMes := "login or email must be provide"
		slog.Error(errMes)
		http.Error(w, errMes, http.StatusBadRequest)
		return
	}

	// Get db user
	var dbUser *models.User
	if usr.Login != "" {
		dbUser, err = us.usrRepo.GetByLogin(r.Context(), usr.Login)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "can't get user with this login", http.StatusBadRequest)
			return
		}
	} else if usr.Email != "" {
		dbUser, err = us.usrRepo.GetByEmail(r.Context(), usr.Email)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "can't get user with this email", http.StatusBadRequest)
			return
		}
	}

	// Check passwords
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(usr.Password))
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "different passwords", http.StatusBadRequest)
		return
	}

	// Create session
	sess, err := us.sessionManager.Create(r.Context(), dbUser.ID)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't create session", http.StatusInternalServerError)
		return
	}

	// Send cookie
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Expires: time.Now().Add(90 * 60 * time.Hour),
		Path:    "/",
	}

	http.SetCookie(w, cookie)
}

func (us *UserService) logout(w http.ResponseWriter, r *http.Request) {
	err := us.sessionManager.DestroySession(r.Context())
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "sess.ID",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}

	http.SetCookie(w, cookie)
}

func (us *UserService) update(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user from db
	dbUser, err := us.usrRepo.GetByLogin(r.Context(), mux.Vars(r)["login"])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error(err.Error())
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get session
	sess, err := us.sessionManager.GetSession(r)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "session error", http.StatusInternalServerError)
		return
	}

	// Check who update user
	if sess.UserID != dbUser.ID {
		errNoPermMes := "no permission to update this user"
		slog.Error(errNoPermMes)
		http.Error(w, errNoPermMes, http.StatusBadRequest)
		return
	}

	// Icon
	file, fileHeader, err := r.FormFile("icon")
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't read file", http.StatusInternalServerError)
		return
	}
	if !strings.Contains(fileHeader.Header.Get("Content-Type"), "image") {
		ErrFilesType := "file must be image"
		slog.Error(ErrFilesType)
		http.Error(w, ErrFilesType, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Parse form values to user
	usr := &models.UserUpdate{
		Email: r.FormValue("email"),
		Login: r.FormValue("login"),
		About: r.FormValue("about"),
	}

	if usr.Email == "" {
		usr.Email = dbUser.Email
	}
	if usr.Login == "" {
		usr.Login = dbUser.Login
	}

	// validate
	err = us.validator.Struct(usr)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// save file to storage
	uri, err := us.s3.Save(file, "image/jpeg")
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't save file", http.StatusInternalServerError)
		return
	}
	oldIconUrl := dbUser.Icon_URL
	usr.Icon_URL = uri

	// Save to storage
	err = us.usrRepo.Update(r.Context(), dbUser.ID, usr)
	if err != nil {
		imgerr := us.s3.Remove(uri)
		if imgerr != nil {
			slog.Error(imgerr.Error())
		}

		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if oldIconUrl != "" {
		err = us.s3.Remove(oldIconUrl)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

func (us *UserService) updatePassword(w http.ResponseWriter, r *http.Request) {
	// Get user from db
	dbUser, err := us.usrRepo.GetByLogin(r.Context(), mux.Vars(r)["login"])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error(err.Error())
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get session
	sess, err := us.sessionManager.GetSession(r)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "session error", http.StatusInternalServerError)
		return
	}

	// Check who update user
	if sess.UserID != dbUser.ID {
		errNoPermMes := "no permission to update this user"
		slog.Error(errNoPermMes)
		http.Error(w, errNoPermMes, http.StatusBadRequest)
		return
	}

	// Read request body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't parse body", http.StatusInternalServerError)
		return
	}

	// Parse form values to user
	updUsr := &models.UserPasswordUpdate{}
	json.Unmarshal(data, &updUsr)

	// validate
	err = us.validator.Struct(updUsr)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash password
	cryptPass, err := bcrypt.GenerateFromPassword([]byte(updUsr.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't hash password", http.StatusInternalServerError)
		return
	}
	updUsr.Password = string(cryptPass)

	// Save to storage
	err = us.usrRepo.UpdatePassword(r.Context(), dbUser.ID, updUsr)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = us.sessionManager.DestroyAllSessions(r.Context(), dbUser.ID)
	if err != nil {
		slog.Error(err.Error())
	}
}

func (us *UserService) get(w http.ResponseWriter, r *http.Request) {
	// Get user from db
	dbUser, err := us.usrRepo.GetByLogin(r.Context(), mux.Vars(r)["login"])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error(err.Error())
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Parse to new structure
	userInfo := &models.UserInfo{}
	userInfo.ID = dbUser.ID
	userInfo.Login = dbUser.Login
	userInfo.About = dbUser.About
	userInfo.Icon_URL = dbUser.Icon_URL
	userInfo.Created_at = dbUser.Created_at

	// Get session
	sess, err := us.sessionManager.GetSession(r)
	if err == nil {
		// Check who get user
		if sess.UserID == dbUser.ID {
			// Get liked recipes
			likedRecipes, err := us.likeRepo.GetLikedRecipies(r.Context(), dbUser.ID)
			if err != nil {
				slog.Error(err.Error())
				http.Error(w, "cant get liked recipes!", http.StatusInternalServerError)
				return
			}
			userInfo.LikedRecipies = likedRecipes
		}
	}

	// Get recipies
	recipies, err := us.recipeRepo.GetAllByUserID(r.Context(), userInfo.ID)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "cant get recipes!", http.StatusInternalServerError)
		return
	}
	userInfo.Recipies = recipies

	// Send json
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userInfo)
}