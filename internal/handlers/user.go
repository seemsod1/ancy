package handlers

import (
	rend "github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/seemsod1/ancy/internal/lib/api/response"
	"github.com/seemsod1/ancy/internal/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (m *Repository) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var req []models.User

	if err := m.App.DB.Find(&req).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rend.JSON(w, r, req)

}

func (m *Repository) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.User

	if err := rend.DecodeJSON(r.Body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("invalid request"))
		return
	} else if err = validator.New().Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("invalid request"))
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to hash password"))
		return
	}
	req.Password = string(pass)

	err = m.App.DB.Create(&req).Error
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		rend.JSON(w, r, response.Error("user with this email already exists or failed to create user"))
		return
	}

	rend.JSON(w, r, response.OK())
}
