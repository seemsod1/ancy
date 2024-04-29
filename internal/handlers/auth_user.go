package handlers

import (
	"net/http"
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
