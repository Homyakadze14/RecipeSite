package like

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Homyakadze14/RecipeSite/internal/recipe"
	"github.com/Homyakadze14/RecipeSite/internal/session"
	"github.com/gorilla/mux"
)

type LikeService struct {
	likeRepo       *LikeRepository
	recipeRepo     *recipe.RecipeRepository
	sessionManager *session.SessionManager
}

func NewService(lr *LikeRepository, sm *session.SessionManager, rr *recipe.RecipeRepository) *LikeService {
	return &LikeService{
		likeRepo:       lr,
		sessionManager: sm,
		recipeRepo:     rr,
	}
}

func (ls *LikeService) HandlFuncs(handler *mux.Router) {
	like := handler.PathPrefix("/recipe/{id:[0-9]+}").Subrouter()
	like.Use(ls.sessionManager.AuthMiddleware)
	like.HandleFunc("/like", ls.like).Methods(http.MethodPost)
	like.HandleFunc("/unlike", ls.unlike).Methods(http.MethodPost)
}

func (ls *LikeService) like(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess, err := ls.sessionManager.GetSession(r)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	recipe, err := ls.recipeRepo.Get(r.Context(), id)
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

	// Form like
	like := &Like{
		UserID:   sess.UserID,
		RecipeID: recipe.ID,
	}

	// Check
	alreadyLike, err := ls.likeRepo.IsAlreadyLike(r.Context(), like)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if alreadyLike {
		ErrAlreadyLike := "This recipe already liked"
		slog.Error(ErrAlreadyLike)
		http.Error(w, ErrAlreadyLike, http.StatusInternalServerError)
		return
	}

	// Update db
	err = ls.likeRepo.Like(r.Context(), like)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ls *LikeService) unlike(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess, err := ls.sessionManager.GetSession(r)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	recipe, err := ls.recipeRepo.Get(r.Context(), id)
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

	// Form like
	like := &Like{
		UserID:   sess.UserID,
		RecipeID: recipe.ID,
	}

	// Check
	alreadyLike, err := ls.likeRepo.IsAlreadyLike(r.Context(), like)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !alreadyLike {
		ErrAlreadyLike := "This recipe isn't liked yet"
		slog.Error(ErrAlreadyLike)
		http.Error(w, ErrAlreadyLike, http.StatusInternalServerError)
		return
	}

	// Update db
	err = ls.likeRepo.Unlike(r.Context(), like)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
