package handler

import (
	"github.com/grdigger/otus-course/internal/model"
	"github.com/grdigger/otus-course/internal/repository"
	"github.com/grdigger/otus-course/internal/service"
	"net/http"
)

type Logout struct {
	session  *service.SessionService
	userRepo repository.UserRepository
}

func NewLogout(session *service.SessionService,
	userRepo repository.UserRepository) *Logout {
	return &Logout{session: session,
		userRepo: userRepo}
}

func (h *Logout) Handle(w http.ResponseWriter, r *http.Request) {
	tpl := service.NewTemplate()

	result, err := h.session.Save(new(model.User), r, w)
	if err != nil {
		tpl.AddVar("error", "ошика  сохранения сессии"+err.Error())
		tpl.Render(w, service.TplNameError)
		return
	}
	if !result {
		tpl.AddVar("error", "ошика  сохранения сессии")
		tpl.Render(w, service.TplNameError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return
}
