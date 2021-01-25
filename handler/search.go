package handler

import (
	"github.com/grdigger/otus-course/internal/repository"
	"github.com/grdigger/otus-course/internal/service"
	l "github.com/sirupsen/logrus"
	"net/http"
)

type Search struct {
	session    *service.SessionService
	friendRepo repository.FriendRepository
}

func NewSearch(session *service.SessionService,
	friendRepo repository.FriendRepository,
) *Search {
	return &Search{session: session, friendRepo: friendRepo}
}

func (h *Search) Handle(w http.ResponseWriter, r *http.Request) {
	tpl := service.NewTemplate()
	us := service.NewUserSession(tpl, h.session, w, r)
	user, err := us.User()
	if err != nil {
		tpl.AddVar("error", err.Error())
		tpl.Render(w, service.TplNameLogin)
		return
	}
	if user.IsEmpty() {
		tpl.Render(w, service.TplNameLogin)
		return
	}

	if r.Method == http.MethodGet {
		tpl.Render(w, service.TplNameSearch)
		return
	}
	if r.Method == http.MethodPost {
		r.ParseForm()
		name := r.FormValue("name")
		surname := r.FormValue("surname")
		users, err := h.friendRepo.Find(name, surname)
		if err != nil {
			l.Errorf("error getting users data: %s ", err.Error())
			tpl.AddVar("error", err.Error())
			tpl.Render(w, service.TplNameError)
			return
		}
		tpl.AddVar("friendList", users)
		tpl.Render(w, service.TplNameFriend)
		return
	}
	tpl.Render(w, service.TplNameFriend)
	return
}
