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
