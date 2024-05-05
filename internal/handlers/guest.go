package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	rend "github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/seemsod1/ancy/internal/helpers"
	"github.com/seemsod1/ancy/internal/lib/api/response"
	"github.com/seemsod1/ancy/internal/models"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type LoginForm struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type SignUpForm struct {
	Username string `json:"username" validate:"required,min=3,max=255,excludesall=!@#$%^&*()_+-="`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	_, ok := m.App.Session.Get(r.Context(), "user_id").(int)
	if ok {
		w.WriteHeader(http.StatusUnauthorized)
		rend.JSON(w, r, response.Error("already logged in"))
		return
	}

	_ = m.App.Session.RenewToken(r.Context())

	var req LoginForm

	if err := rend.DecodeJSON(r.Body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("invalid request"))
		return
	}

	var user models.User
	if err := validator.New().Var(req.Login, "required"); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("invalid request"))
		return
	} else {
		if err = validator.New().Var(req.Login, "contains=@"); err == nil {
			if err = m.App.DB.Table("users").Preload("Role").Where("email = ?", req.Login).Take(&user).Error; err != nil {
				w.WriteHeader(http.StatusNotFound)
				rend.JSON(w, r, response.Error("user not found"))
				return
			}
		} else {
			if err = m.App.DB.Table("users").Preload("Role").Where("username = ?", req.Login).Take(&user).Error; err != nil {
				w.WriteHeader(http.StatusNotFound)
				rend.JSON(w, r, response.Error("user not found"))
				return
			}
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		rend.JSON(w, r, response.Error("invalid password"))
		return
	}
	m.App.Session.Put(r.Context(), "user_id", user.ID)
	m.App.Session.Put(r.Context(), "user_role", user.RoleID)

	if user.Role.Name == "Admin" {
		m.App.Session.Put(r.Context(), "is_admin", true)
	} else {
		m.App.Session.Put(r.Context(), "is_admin", false)

	}

	rend.JSON(w, r, response.OK())

}

func (m *Repository) SignUp(w http.ResponseWriter, r *http.Request) {
	_, ok := m.App.Session.Get(r.Context(), "user_id").(int)
	if ok {
		w.WriteHeader(http.StatusUnauthorized)
		rend.JSON(w, r, response.Error("already logged in"))
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("file too big"))
		return
	}

	var req SignUpForm
	req.Username = r.FormValue("username")
	req.Email = r.FormValue("email")
	req.Password = r.FormValue("password")

	// Validate form inputs
	if err := validator.New().Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("invalid request"))
		return
	}

	// Process profile photo
	var profilePhotoPath string
	profilePhoto, fileHeader, err := r.FormFile("profile_photo")
	if err != nil {
		profilePhotoPath = "default.png"
	} else {
		defer profilePhoto.Close()

		finalTitle, err := helpers.GenerateHashedFileName(req.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			rend.JSON(w, r, response.Error("failed to generate file name"))
			return
		}
		fileFormat := helpers.GetFileFormat(fileHeader.Filename)
		profilePhotoPath = finalTitle + "." + fileFormat

		newFile, err := os.Create("storage/users/" + profilePhotoPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			rend.JSON(w, r, response.Error("failed to save file"))
			return
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, profilePhoto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			rend.JSON(w, r, response.Error("failed to save file"))
			return
		}
	}

	// Hash password
	pass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to hash password"))
		return
	}
	req.Password = string(pass)

	var user models.User
	user.Username = req.Username
	user.Email = req.Email
	user.Password = req.Password
	user.ProfilePhotoPath = profilePhotoPath

	err = m.App.DB.Create(&user).Error
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		rend.JSON(w, r, response.Error("user with this email or username already exists or failed to create user"))
		return
	}

	rend.JSON(w, r, response.OK())
}

func (m *Repository) GetExhibit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("exhibit ID is required"))
		return
	}
	var eId int
	if _, err := fmt.Sscanf(id, "%d", &eId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("invalid exhibit ID"))
		return
	}
	var exhibit models.Exhibit
	if err := m.App.DB.Table("exhibits").Preload("Author").Preload("Type").Preload("Status").Where("id = ?", eId).Take(&exhibit).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		rend.JSON(w, r, response.Error("exhibit not found"))
		return
	}
	role := m.GetLoggedInUserRole(r.Context())
	if (role == "User" || role == "") && (exhibit.Status.Name == "Pending" || exhibit.Status.Name == "Rejected") && exhibit.AuthorID != m.GetLoggedInUserID(r.Context()) {
		w.WriteHeader(http.StatusUnauthorized)
		rend.JSON(w, r, response.Error("unauthorized"))
		return
	}
	exhibit.Author.Password = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exhibit)
}

func (m *Repository) GetAllExhibits(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	username := r.URL.Query().Get("username")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			rend.JSON(w, r, response.Error("invalid start date format"))
			return
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			rend.JSON(w, r, response.Error("invalid end date format"))
			return
		}
	}

	dbQuery := m.App.DB.Preload("Type").Preload("Status").Joins("JOIN exhibit_statuses ON exhibits.status_id = exhibit_statuses.id")
	dbQuery = dbQuery.Joins("JOIN users ON exhibits.author_id = users.id")

	if status != "" {
		dbQuery = dbQuery.Where("exhibit_statuses.name = ?", status)
	}

	if username != "" {
		usernameLower := strings.ToLower(username)
		dbQuery = dbQuery.Where("LOWER(users.username) LIKE ?", "%"+usernameLower+"%")
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		dbQuery = dbQuery.Where("exhibits.created_at BETWEEN ? AND ?", startDate, endDate)
	} else if !startDate.IsZero() {
		dbQuery = dbQuery.Where("exhibits.created_at >= ?", startDate)
	} else if !endDate.IsZero() {
		dbQuery = dbQuery.Where("exhibits.created_at <= ?", endDate)
	}

	userRole := m.GetLoggedInUserRole(r.Context())
	if userRole == "Admin" && status != "" {
		dbQuery = dbQuery.Where("exhibit_statuses.name = ?", status)
	} else if userRole != "Admin" {
		dbQuery = dbQuery.Where("exhibit_statuses.name = ?", "Approved")
	}

	var exhibits []models.Exhibit
	if err = dbQuery.Preload("Author").Find(&exhibits).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rend.JSON(w, r, response.Error("failed to get exhibits"))
		return
	}
	for i := range exhibits {
		exhibits[i].Author.Password = ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exhibits)
}

func (m Repository) GetUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		rend.JSON(w, r, response.Error("username is required"))
		return
	}

	var user models.User
	if err := m.App.DB.Where("username = ?", username).Take(&user).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		rend.JSON(w, r, response.Error("user not found"))
		return
	}
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
