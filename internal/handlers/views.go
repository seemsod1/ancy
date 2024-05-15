package handlers

import (
	rend "github.com/go-chi/render"
	"github.com/seemsod1/ancy/internal/lib/api/response"
	"github.com/seemsod1/ancy/internal/models"
	"github.com/seemsod1/ancy/internal/render"
	"net/http"
)

func (m *Repository) Search(w http.ResponseWriter, r *http.Request) {

	err := render.Template(w, r, "search.page.tmpl", &models.TemplateData{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to update user"))
		return
	}
}
func (m *Repository) Exhibit(w http.ResponseWriter, r *http.Request) {

	err := render.Template(w, r, "exhibit.page.tmpl", &models.TemplateData{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to update user"))
		return
	}
}
