package recipe

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/Homyakadze14/RecipeSite/internal/images"
	"github.com/Homyakadze14/RecipeSite/internal/jsonvalidator"
	"github.com/Homyakadze14/RecipeSite/internal/session"
	"github.com/Homyakadze14/RecipeSite/internal/user"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type RecipeService struct {
	recipeRepo     *RecipeRepository
	userRepo       *user.UserRepository
	validator      *jsonvalidator.JSONValidator
	sessionManager *session.SessionManager
}

func NewService(recipeRepo *RecipeRepository, validator *jsonvalidator.JSONValidator, sm *session.SessionManager, ur *user.UserRepository) *RecipeService {
	return &RecipeService{
		recipeRepo:     recipeRepo,
		validator:      validator,
		userRepo:       ur,
		sessionManager: sm,
	}
}

func (rs *RecipeService) HandlFuncs(handler *mux.Router) {
	recipe := handler.PathPrefix("/recipe").Subrouter()
	recipe.Use(rs.sessionManager.AuthMiddleware)
	recipe.HandleFunc("", rs.getAll).Methods(http.MethodGet)
	recipe.HandleFunc("/{login}", rs.create).Methods(http.MethodPost)
}

func (rs *RecipeService) getAll(w http.ResponseWriter, r *http.Request) {
	recipes, err := rs.recipeRepo.GetAll(r.Context())
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]Recipe{"recipes": recipes})
}

func (rs *RecipeService) create(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// Get user from db
	dbUser, err := rs.userRepo.GetByLogin(r.Context(), mux.Vars(r)["login"])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error(err.Error())
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get session
	sess, err := rs.sessionManager.GetSession(r)
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

	// Parse form values to user
	complexitiy, err := strconv.Atoi(r.FormValue("complexitiy"))
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "complexitiy must be integer", http.StatusBadRequest)
		return
	}

	recipe := &Recipe{
		User_ID:     dbUser.ID,
		Title:       r.FormValue("title"),
		About:       r.FormValue("about"),
		Complexitiy: complexitiy,
		NeedTime:    r.FormValue("need_time"),
		Ingridients: r.FormValue("ingridients"),
	}

	// validate
	err = rs.validator.Struct(recipe)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse photos and save
	multipartFormData := r.MultipartForm
	files := multipartFormData.File["photos"]

	if len(files) == 0 {
		ErrFiles := "photos must be provided"
		slog.Error(ErrFiles)
		http.Error(w, ErrFiles, http.StatusBadRequest)
		return
	}

	uid := uuid.New().String()
	for _, v := range files {
		if !strings.Contains(v.Header.Get("Content-Type"), "image") {
			ErrFilesType := "files must be images"
			slog.Error(ErrFilesType)
			http.Error(w, ErrFilesType, http.StatusBadRequest)
			return
		}

		uploadedFile, err := v.Open()
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "can't read file", http.StatusInternalServerError)
			return
		}

		uri, err := images.Save(fmt.Sprintf("%v/recipes/%s", dbUser.ID, uid), uploadedFile)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "can't save files", http.StatusInternalServerError)
			return
		}
		recipe.Photos_URLS += uri + ";"

		uploadedFile.Close()
	}

	// Save to storage
	err = rs.recipeRepo.Create(r.Context(), recipe)
	if err != nil {
		errImage := images.Remove(recipe.Photos_URLS)
		if errImage != nil {
			slog.Error(errImage.Error())
		}

		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
