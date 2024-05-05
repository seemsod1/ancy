package handlers

import (
	"github.com/go-chi/chi/v5"
	rend "github.com/go-chi/render"
	"github.com/seemsod1/ancy/internal/lib/api/response"
	"github.com/seemsod1/ancy/internal/models"
	"net/http"
)

func (m *Repository) ApproveExhibit(w http.ResponseWriter, r *http.Request) {
	exhibitID := chi.URLParam(r, "id")
	if exhibitID == "" {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("exhibit ID is required"))
		return
	}

	var statusID int
	if err := m.App.DB.Table("exhibit_statuses").Where("name = ?", "Approved").Pluck("id", &statusID).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to get status"))
		return
	}

	if err := m.App.DB.Model(&models.Exhibit{}).Where("id = ?", exhibitID).Update("status_id", statusID).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to approve exhibit"))
		return
	}

	w.WriteHeader(http.StatusOK)
	rend.JSON(w, r, response.OK())
}

func (m *Repository) RejectExhibit(w http.ResponseWriter, r *http.Request) {
	exhibitID := chi.URLParam(r, "id")
	if exhibitID == "" {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("exhibit ID is required"))
		return
	}

	var statusID int
	if err := m.App.DB.Table("exhibit_statuses").Where("name = ?", "Rejected").Pluck("id", &statusID).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to get status"))
		return
	}

	if err := m.App.DB.Model(&models.Exhibit{}).Where("id = ?", exhibitID).Update("status_id", statusID).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to reject exhibit"))
		return
	}

	w.WriteHeader(http.StatusOK)
	rend.JSON(w, r, response.OK())
}

func (m Repository) MakeAdmin(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("user ID is required"))
		return
	}

	var roleID int
	if err := m.App.DB.Table("user_roles").Where("name = ?", "Admin").Pluck("id", &roleID).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to get role"))
		return
	}

	if err := m.App.DB.Model(&models.User{}).Where("id = ?", userID).Update("role_id", roleID).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to make user admin"))
		return
	}

	w.WriteHeader(http.StatusOK)
	rend.JSON(w, r, response.OK())
}

func (m Repository) RemoveAdmin(w http.ResponseWriter, r *http.Request) {
	uID := m.App.Session.Get(r.Context(), "user_id").(string)

	userID := chi.URLParam(r, "id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("user ID is required"))
		return
	}
	if uID != "1" {
		w.WriteHeader(http.StatusForbidden)
		rend.JSON(w, r, response.Error("you cannot remove the main admin user"))
		return

	}

	var roleID int
	if err := m.App.DB.Table("user_roles").Where("name = ?", "User").Pluck("id", &roleID).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to get role"))
		return
	}

	if err := m.App.DB.Model(&models.User{}).Where("id = ?", userID).Update("role_id", roleID).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to remove user admin"))
		return
	}

	w.WriteHeader(http.StatusOK)
	rend.JSON(w, r, response.OK())
}

func (m Repository) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("user ID is required"))
		return
	}
	if userID == "1" {
		w.WriteHeader(http.StatusForbidden)
		rend.JSON(w, r, response.Error("you cannot delete the main admin user"))
		return

	}

	withExhibits := r.URL.Query().Get("withExhibits")
	if withExhibits == "true" {
		if err := m.App.DB.Where("author_id = ?", userID).Delete(&models.Exhibit{}).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			rend.JSON(w, r, response.Error("failed to delete user exhibits"))
			return
		}
	}

	if err := m.App.DB.Where("id = ?", userID).Delete(&models.User{}).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to delete user"))
		return
	}

	w.WriteHeader(http.StatusOK)
	rend.JSON(w, r, response.OK())
}

func (m Repository) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := m.App.DB.Preload("Role").Find(&users).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to get users"))
		return
	}

	for i := range users {
		users[i].Password = ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	rend.JSON(w, r, users)
}
