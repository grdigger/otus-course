package handler

import (
	"github.com/grdigger/otus-course/internal/repository"
	"github.com/grdigger/otus-course/internal/service"
	"net/http"
)

type Login struct {
	session  *service.SessionService
	userRepo repository.UserRepository
}

func NewLogin(session *service.SessionService,
	userRepo repository.UserRepository) *Login {
	return &Login{session: session,
		userRepo: userRepo}
}

func (h *Login) Handle(w http.ResponseWriter, r *http.Request) {
	tpl := service.NewTemplate()
	err := r.ParseForm()
	if err != nil {
		tpl.AddVar("error", "error parsing form: "+err.Error())
		tpl.Render(w, service.TplNameError)
		return
	}
	user, err := h.userRepo.Auth(r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		tpl.Render(w, service.TplNameError)
		return
	}
	if user.IsEmpty() {
		tpl.AddVar("userNotFound", "Пользователь не найден")
		tpl.Render(w, service.TplNameLogin)
		return
	}
	result, err := h.session.Save(user, r, w)
	if err != nil {
		tpl.AddVar("error", "ошика  сохранения сессии"+err.Error())
		tpl.Render(w, service.TplNameLogin)
		return
	}
	if !result {
		tpl.AddVar("error", "ошика  сохранения сессии")
		tpl.Render(w, service.TplNameLogin)
		return
	}
	http.Redirect(w, r, "/personal", http.StatusTemporaryRedirect)
	return
}
