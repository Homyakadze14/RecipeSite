package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/Homyakadze14/RecipeSite/internal/images"
	"github.com/Homyakadze14/RecipeSite/internal/jsonvalidator"
	"github.com/Homyakadze14/RecipeSite/internal/models"
	"github.com/Homyakadze14/RecipeSite/internal/repos"
	"github.com/Homyakadze14/RecipeSite/internal/session"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type RecipeService struct {
	recipeRepo     *repos.RecipeRepository
	userRepo       *repos.UserRepository
	likeRepo       *repos.LikeRepository
	sessionManager *session.SessionManager
	commentRepo    *repos.CommentRepository
	validator      *jsonvalidator.JSONValidator
}

func NewRecipeService(rr *repos.RecipeRepository, sm *session.SessionManager,
	ur *repos.UserRepository, lr *repos.LikeRepository, cr *repos.CommentRepository, v *jsonvalidator.JSONValidator) *RecipeService {
	return &RecipeService{
		recipeRepo:     rr,
		validator:      v,
		userRepo:       ur,
		sessionManager: sm,
		likeRepo:       lr,
		commentRepo:    cr,
	}
}

func (rs *RecipeService) HandlFuncs(handler *mux.Router) {
	recipe := handler.PathPrefix("/recipe").Subrouter()
	recipe.HandleFunc("", rs.getAll).Methods(http.MethodGet)
	recipe.HandleFunc("", rs.getFiltered).Methods(http.MethodPost)
	recipe.HandleFunc("/{id:[0-9]+}", rs.get).Methods(http.MethodGet)

	userRecipe := handler.PathPrefix("/user/{login}/recipe").Subrouter()
	userRecipe.Use(rs.sessionManager.AuthMiddleware)
	userRecipe.HandleFunc("", rs.create).Methods(http.MethodPost)
	userRecipe.HandleFunc("/{id:[0-9]+}", rs.update).Methods(http.MethodPut)
	userRecipe.HandleFunc("/{id:[0-9]+}", rs.delete).Methods(http.MethodDelete)
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
	json.NewEncoder(w).Encode(map[string][]models.Recipe{"recipes": recipes})
}

func (rs *RecipeService) getFiltered(w http.ResponseWriter, r *http.Request) {
	// Read request body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't parse body", http.StatusInternalServerError)
		return
	}

	// Parse json values to user
	filter := &models.RecipeFilter{}
	err = json.Unmarshal(data, &filter)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't parse json", http.StatusInternalServerError)
		return
	}

	// validate
	err = rs.validator.Struct(filter)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get filtered recipes
	recipes, err := rs.recipeRepo.GetFiltered(r.Context(), filter)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]models.Recipe{"recipes": recipes})
}

func (rs *RecipeService) get(w http.ResponseWriter, r *http.Request) {
	// Parse recipe id
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "id must be integer", http.StatusBadRequest)
		return
	}

	fullRecipe := &models.FullRecipe{}

	// Get recipe
	recipe, err := rs.recipeRepo.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error(err.Error())
			http.Error(w, "recipe not found", http.StatusNotFound)
			return
		}
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fullRecipe.Recipe = recipe

	// Get Author
	author, err := rs.userRepo.GetAuthor(r.Context(), recipe.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error(err.Error())
			http.Error(w, "author not found", http.StatusNotFound)
			return
		}
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fullRecipe.Author = author

	// Get likes count
	fullRecipe.LikesCount, err = rs.likeRepo.LikesCount(r.Context(), recipe.ID)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get comments
	fullRecipe.Comments, err = rs.commentRepo.GetCommets(r.Context(), recipe.ID, rs.userRepo)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get session
	sess, err := rs.sessionManager.GetSession(r)
	// Check auth user
	if err == nil {
		like := &models.Like{
			UserID:   sess.UserID,
			RecipeID: recipe.ID,
		}
		fullRecipe.IsLiked, err = rs.likeRepo.IsAlreadyLike(r.Context(), like)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]*models.FullRecipe{"info": fullRecipe})
}

func (rs *RecipeService) create(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	// Check who create recipe
	if sess.UserID != dbUser.ID {
		errNoPermMes := "no permission to create recipe to this user"
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

	recipe := &models.Recipe{
		UserID:      dbUser.ID,
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
		recipe.PhotosUrls += uri + ";"

		uploadedFile.Close()
	}

	// Save to storage
	err = rs.recipeRepo.Create(r.Context(), recipe)
	if err != nil {
		errImage := images.Remove(recipe.PhotosUrls)
		if errImage != nil {
			slog.Error(errImage.Error())
		}

		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rs *RecipeService) update(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	// Check who update recipe
	if sess.UserID != dbUser.ID {
		errNoPermMes := "no permission to update recipe to this user"
		slog.Error(errNoPermMes)
		http.Error(w, errNoPermMes, http.StatusBadRequest)
		return
	}

	// Parse recipe id
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "id must be integer", http.StatusBadRequest)
		return
	}

	// Get recipe
	dbRecipe, err := rs.recipeRepo.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error(err.Error())
			http.Error(w, "recipe not found", http.StatusNotFound)
			return
		}
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse form values to recipe
	if r.FormValue("complexitiy") != "" {
		dbRecipe.Complexitiy, err = strconv.Atoi(r.FormValue("complexitiy"))
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "complexitiy must be integer", http.StatusBadRequest)
			return
		}
	}
	if r.FormValue("title") != "" {
		dbRecipe.Title = r.FormValue("title")
	}
	if r.FormValue("about") != "" {
		dbRecipe.About = r.FormValue("about")
	}
	if r.FormValue("need_time") != "" {
		dbRecipe.NeedTime = r.FormValue("need_time")
	}
	if r.FormValue("ingridients") != "" {
		dbRecipe.Ingridients = r.FormValue("ingridients")
	}

	// validate
	err = rs.validator.Struct(dbRecipe)
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

	oldPhotos := dbRecipe.PhotosUrls
	dbRecipe.PhotosUrls = ""
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
		dbRecipe.PhotosUrls += uri + ";"

		uploadedFile.Close()
	}

	// Update recipe in storage
	err = rs.recipeRepo.Update(r.Context(), id, dbRecipe)
	if err != nil {
		errImage := images.Remove(dbRecipe.PhotosUrls)
		if errImage != nil {
			slog.Error(errImage.Error())
		}

		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete old photos
	if oldPhotos != "" {
		err = images.Remove(oldPhotos)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

func (rs *RecipeService) delete(w http.ResponseWriter, r *http.Request) {
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

	// Check who delete recipe
	if sess.UserID != dbUser.ID {
		errNoPermMes := "no permission to delete this recipe"
		slog.Error(errNoPermMes)
		http.Error(w, errNoPermMes, http.StatusBadRequest)
		return
	}

	// Parse recipe id
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "id must be integer", http.StatusBadRequest)
		return
	}

	// Get recipe
	recipe, err := rs.recipeRepo.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error(err.Error())
			http.Error(w, "recipe not found", http.StatusNotFound)
			return
		}
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete recipe
	err = rs.recipeRepo.Delete(r.Context(), recipe)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
