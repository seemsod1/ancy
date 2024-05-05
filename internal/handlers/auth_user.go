package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	rend "github.com/go-chi/render"
	"github.com/seemsod1/ancy/internal/helpers"
	"github.com/seemsod1/ancy/internal/lib/api/response"
	"github.com/seemsod1/ancy/internal/models"
	"io"
	"net/http"
	"os"
	"strconv"
)

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_, ok := m.App.Session.Get(r.Context(), "user_id").(int)
	if !ok {
		//you are not logged in
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_ = m.App.Session.Destroy(r.Context())
}

func (m *Repository) CreateExhibit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("file too big"))
		return
	}

	title := r.Form.Get("title")
	typeID, err := strconv.Atoi(r.Form.Get("type"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("invalid type"))
		return

	}
	description := r.Form.Get("description")
	file, fileHeader, err := r.FormFile("file") // Отримуємо файл з форми
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("failed to get file"))
		return
	}
	defer file.Close()

	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("title is required"))
		return
	}
	finalTitle, err := helpers.GenerateHashedFileName(title)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to generate file name"))
		return

	}

	fileFormat := helpers.GetFileFormat(fileHeader.Filename)
	filePath := finalTitle + "." + fileFormat
	newFile, err := os.Create("storage/" + filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to save file"))
		return
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to save file"))
		return
	}

	authorID, _ := m.App.Session.Get(r.Context(), "user_id").(int)
	var statusID int
	if err = m.App.DB.Table("exhibit_statuses").Where("name = ?", "Pending").Pluck("id", &statusID).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to get status"))
		return
	}

	var tId int
	if err = m.App.DB.Table("exhibit_types").Where("id = ?", typeID).Pluck("id", &tId).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to get type"))
		return
	}

	exhibit := models.Exhibit{
		Title:       title,
		TypeID:      typeID,
		Description: description,
		AssetPath:   filePath,
		AuthorID:    authorID,
		StatusID:    statusID,
	}

	if err = m.App.DB.Create(&exhibit).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to create exhibit"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	rend.JSON(w, r, response.OK())

}

func (m *Repository) GetMyExhibits(w http.ResponseWriter, r *http.Request) {
	authorID, _ := m.App.Session.Get(r.Context(), "user_id").(int)
	var exhibits []models.Exhibit
	if err := m.App.DB.Where("author_id = ?", authorID).Preload("Type").Preload("Status").Find(&exhibits).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to get exhibits"))
		return
	}

	jsonData, err := json.Marshal(exhibits)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to marshal exhibits to JSON"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (m *Repository) DeleteExhibit(w http.ResponseWriter, r *http.Request) {
	exhibitID := chi.URLParam(r, "id")
	if exhibitID == "" {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("exhibit ID is required"))
		return
	}
	var eId int
	if _, err := fmt.Sscanf(exhibitID, "%d", &eId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("invalid exhibit ID"))
		return
	}

	authorID, _ := m.App.Session.Get(r.Context(), "user_id").(int)
	var exhibit models.Exhibit
	if err := m.App.DB.Where("id = ?", eId).First(&exhibit).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		rend.JSON(w, r, response.Error("exhibit not found"))
		return
	}
	userRole := m.GetLoggedInUserRole(r.Context())

	if userRole != "Admin" && exhibit.AuthorID != authorID {
		w.WriteHeader(http.StatusForbidden)
		rend.JSON(w, r, response.Error("forbidden"))
		return
	}

	if err := os.Remove("/storage" + exhibit.AssetPath); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to delete exhibit file"))
		return
	}

	if err := m.App.DB.Delete(&exhibit).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to delete exhibit"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (m *Repository) UpdatePhoto(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("file too big"))
		return
	}

	file, fileHeader, err := r.FormFile("file") // Отримуємо файл з форми
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("failed to get file"))
		return
	}
	defer file.Close()

	authorID, _ := m.App.Session.Get(r.Context(), "user_id").(int)
	var user models.User
	if err := m.App.DB.First(&user, authorID).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	finalTitle, err := helpers.GenerateHashedFileName(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to generate file name"))
		return

	}

	fileFormat := helpers.GetFileFormat(fileHeader.Filename)
	filePath := finalTitle + "." + fileFormat
	newFile, err := os.Create("storage/users/" + filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to save file"))
		return
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to save file"))
		return
	}
	// Delete old photo
	if user.ProfilePhotoPath != "default.png" {
		if err = os.Remove("storage/users/" + user.ProfilePhotoPath); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			rend.JSON(w, r, response.Error("failed to delete old photo"))
			return
		}
	}

	user.ProfilePhotoPath = filePath
	if err = m.App.DB.Save(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to update user"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
