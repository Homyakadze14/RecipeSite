package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/Homyakadze14/RecipeSite/RecipeSite/internal/images"
	"github.com/Homyakadze14/RecipeSite/RecipeSite/internal/jsonvalidator"
	"github.com/Homyakadze14/RecipeSite/RecipeSite/internal/session"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

const defaultIconURL = ""

type UserService struct {
	usrRepo        *UserRepository
	validator      *jsonvalidator.JSONValidator
	sessionManager *session.SessionManager
}

func NewService(usrRepo *UserRepository, validator *jsonvalidator.JSONValidator, sm *session.SessionManager) *UserService {
	return &UserService{
		usrRepo:        usrRepo,
		validator:      validator,
		sessionManager: sm,
	}
}

func (us *UserService) HandlFuncs(handler *mux.Router) {
	auth := handler.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/signup", us.signup).Methods(http.MethodPost)
	auth.HandleFunc("/signin", us.signin).Methods(http.MethodPost)

	logout := auth.Path("/logout").Subrouter()
	logout.Use(us.sessionManager.AuthMiddleware)
	logout.HandleFunc("", us.logout).Methods(http.MethodPost)

	user := handler.PathPrefix("/user").Subrouter()
	user.Use(us.sessionManager.AuthMiddleware)
	user.HandleFunc("/{login}", us.update).Methods(http.MethodPut)
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
	usr := &User{}
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
	usr.Icon_URL = defaultIconURL

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
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	usr := &UserLogin{}
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
	var dbUser *User
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
	// Get user from db
	dbUser, err := us.usrRepo.GetByLogin(r.Context(), mux.Vars(r)["login"])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error(err.Error())
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
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
	file, _, err := r.FormFile("icon")
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't read file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Parse form values to user
	usr := &UserUpdate{
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
	uri, err := images.Save(fmt.Sprintf("%v", dbUser.ID), file)
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
		imgerr := images.Remove(uri)
		if imgerr != nil {
			slog.Error(imgerr.Error())
			http.Error(w, imgerr.Error(), http.StatusBadRequest)
		}

		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if oldIconUrl != "" {
		err = images.Remove(oldIconUrl)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
