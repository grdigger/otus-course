package handler

import (
	"fmt"
	"net/http"

	"github.com/grdigger/otus-course/internal/repository"
	"github.com/grdigger/otus-course/internal/service"
	l "github.com/sirupsen/logrus"
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
	err = fmt.Errorf("error from auth")
	if err != nil {
		l.Errorf("auth error: %s ", err.Error())
		tpl.AddVar("error", err.Error())
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
		l.Errorf("ошика  сохранения сессии: %s ", err.Error())
		tpl.AddVar("error", "ошика  сохранения сессии"+err.Error())
		tpl.Render(w, service.TplNameLogin)
		return
	}
	if !result {
		l.Errorf("ошика  сохранения сессии: %s ", err.Error())
		tpl.AddVar("error", "ошика  сохранения сессии")
		tpl.Render(w, service.TplNameLogin)
		return
	}
	http.Redirect(w, r, "/personal", http.StatusTemporaryRedirect)
	return
}
