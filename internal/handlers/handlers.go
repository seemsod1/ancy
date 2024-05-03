package handlers

import (
	"context"
	"github.com/seemsod1/ancy/internal/config"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) GetLoggedInUserRole(ctx context.Context) string {
	roleId, _ := m.App.Session.Get(ctx, "user_role").(int)
	var role string
	if err := m.App.DB.Table("user_roles").Where("id = ?", roleId).Pluck("name", &role).Error; err != nil {
		return ""
	}
	return role
}

func (m *Repository) GetLoggedInUserID(ctx context.Context) int {
	userId, _ := m.App.Session.Get(ctx, "user_id").(int)
	return userId
}
