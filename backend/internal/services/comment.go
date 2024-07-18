package services

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Homyakadze14/RecipeSite/internal/jsonvalidator"
	"github.com/Homyakadze14/RecipeSite/internal/models"
	"github.com/Homyakadze14/RecipeSite/internal/repos"
	"github.com/Homyakadze14/RecipeSite/internal/session"
	"github.com/gorilla/mux"
)

type Service struct {
	cr *repos.CommentRepository
	sm *session.SessionManager
	vd *jsonvalidator.JSONValidator
}

func NewCommentService(cr *repos.CommentRepository, sm *session.SessionManager, vd *jsonvalidator.JSONValidator) *Service {
	return &Service{
		cr: cr,
		sm: sm,
		vd: vd,
	}
}

func (cs *Service) HandlFuncs(handler *mux.Router) {
	comment := handler.PathPrefix("/recipe/{id:[0-9]+}/comment").Subrouter()
	comment.Use(cs.sm.AuthMiddleware)
	comment.HandleFunc("", cs.addComment).Methods(http.MethodPost)
	comment.HandleFunc("", cs.updateComment).Methods(http.MethodPut)
	comment.HandleFunc("", cs.deleteComment).Methods(http.MethodDelete)
}

func (cs *Service) addComment(w http.ResponseWriter, r *http.Request) {
	// Parse json comment
	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't read request", http.StatusInternalServerError)
		return
	}

	comment := &models.Comment{}
	err = json.Unmarshal(data, &comment)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't parse json", http.StatusInternalServerError)
		return
	}

	// Validate
	err = cs.vd.Struct(comment)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user id from session
	sess, err := cs.sm.GetSession(r)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	comment.UserID = sess.UserID

	// Get recipe_id from url
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "id must be integer", http.StatusBadRequest)
		return
	}
	comment.RecipeID = id

	// Save to database
	err = cs.cr.Save(r.Context(), comment)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't save this comment", http.StatusInternalServerError)
		return
	}
}

func (cs *Service) updateComment(w http.ResponseWriter, r *http.Request) {
	// Parse json comment
	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't read request", http.StatusInternalServerError)
		return
	}

	comment := &models.CommentUpdate{}
	err = json.Unmarshal(data, &comment)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't parse json", http.StatusInternalServerError)
		return
	}

	// Validate
	err = cs.vd.Struct(comment)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update in db
	err = cs.cr.Update(r.Context(), comment)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't update this comment", http.StatusInternalServerError)
		return
	}
}

func (cs *Service) deleteComment(w http.ResponseWriter, r *http.Request) {
	// Parse json comment
	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't read request", http.StatusInternalServerError)
		return
	}

	comment := &models.CommentDelete{}
	err = json.Unmarshal(data, &comment)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't parse json", http.StatusInternalServerError)
		return
	}

	// Validate
	err = cs.vd.Struct(comment)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Delete in db
	err = cs.cr.Delete(r.Context(), comment)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "can't delete this comment", http.StatusInternalServerError)
		return
	}
}
