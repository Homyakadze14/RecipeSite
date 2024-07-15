package services

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Homyakadze14/RecipeSite/internal/models"
	"github.com/Homyakadze14/RecipeSite/internal/repos"
	"github.com/Homyakadze14/RecipeSite/internal/session"
	"github.com/gorilla/mux"
)

type LikeService struct {
	likeRepo       *repos.LikeRepository
	sessionManager *session.SessionManager
}

func NewLikeService(lr *repos.LikeRepository, sm *session.SessionManager) *LikeService {
	return &LikeService{
		likeRepo:       lr,
		sessionManager: sm,
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

	// Form like
	like := &models.Like{
		UserID:   sess.UserID,
		RecipeID: id,
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

	// Form like
	like := &models.Like{
		UserID:   sess.UserID,
		RecipeID: id,
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
