package handlers

import (
	"github.com/go-chi/chi/v5"
	rend "github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/seemsod1/ancy/internal/lib/api/response"
	"github.com/seemsod1/ancy/internal/models"
	"net/http"
)

func (m *Repository) GetAllUserRoles(w http.ResponseWriter, r *http.Request) {
	var req []models.UserRole

	if err := m.App.DB.Find(&req).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//render all roles
	rend.JSON(w, r, req)
}

func (m *Repository) CreateUserRole(w http.ResponseWriter, r *http.Request) {
	var req models.UserRole

	if err := rend.DecodeJSON(r.Body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("invalid request"))
		return
	}

	err := m.App.DB.Create(&req).Error
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		rend.JSON(w, r, response.Error("role already exists or failed to create role"))
		return
	}

	rend.JSON(w, r, response.OK())
}

func (m *Repository) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	var req models.UserRole

	if err := rend.DecodeJSON(r.Body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("invalid request"))
		return
	}
	if err := validator.New().Var(req.ID, "required,numeric"); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("invalid request"))
		return
	}
	if err := validator.New().Var(req.Name, "required"); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("invalid request"))
		return
	}
	if err := m.App.DB.First(&models.UserRole{}, req.ID).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		rend.JSON(w, r, response.NotFound("role not found"))
		return
	}

	err := m.App.DB.Save(&req).Error
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		rend.JSON(w, r, response.Error("failed to update role"))
		return
	}

	rend.JSON(w, r, response.OK())
}

func (m *Repository) DeleteUserRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("invalid request"))
		return
	} else {
		if err := validator.New().Var(id, "required,numeric"); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			rend.JSON(w, r, response.Error("invalid request"))
			return
		}
	}

	var req models.UserRole
	if err := m.App.DB.First(&req, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		rend.JSON(w, r, response.NotFound("role not found"))
		return
	}

	err := m.App.DB.Delete(&req).Error
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		rend.JSON(w, r, response.Error("failed to delete role"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
