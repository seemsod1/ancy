package handlers

import (
	rend "github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/seemsod1/ancy/internal/lib/api/response"
	"github.com/seemsod1/ancy/internal/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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
			if err = m.App.DB.Table("users").Where("email = ?", req.Login).Take(&user).Error; err != nil {
				w.WriteHeader(http.StatusNotFound)
				rend.JSON(w, r, response.Error("user not found"))
				return
			}
		} else {
			if err = m.App.DB.Table("users").Where("username = ?", req.Login).Take(&user).Error; err != nil {
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

	rend.JSON(w, r, response.OK())

}

func (m *Repository) SingUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpForm

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

	var user models.User
	user.Username = req.Username
	user.Email = req.Email
	user.Password = req.Password

	err = m.App.DB.Create(&user).Error
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		rend.JSON(w, r, response.Error("user with this email already exists or failed to create user"))
		return
	}

	rend.JSON(w, r, response.OK())
}
