package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/grdigger/otus-course/internal/model"
	"github.com/grdigger/otus-course/internal/repository"
	"github.com/grdigger/otus-course/internal/service"
	"net/http"
	"strconv"
)

type FriendAdd struct {
	session    *service.SessionService
	friendRepo repository.FriendRepository
}

func NewFriendAdd(session *service.SessionService,
	friendRepo repository.FriendRepository,
) *FriendAdd {
	return &FriendAdd{session: session,
		friendRepo: friendRepo,
	}
}

func (h *FriendAdd) Handle(w http.ResponseWriter, r *http.Request) {
	tpl := service.NewTemplate()
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		tpl.AddVar("error", "error parsing ID: "+err.Error())
		tpl.Render(w, service.TplNameError)
		return
	}
	userData, err := h.session.UserSessions(r)
	if err != nil {
		tpl.AddVar("error", "server error: "+err.Error())
		tpl.Render(w, service.TplNameLogin)
		return
	}
	user, ok := userData.(model.User)
	if !ok {
		tpl.AddVar("error", "нужно сначала залогинится")
		tpl.Render(w, service.TplNameError)
		return
	}
	err = h.friendRepo.Add(user.ID, userID)
	if err != nil {
		tpl.AddVar("error", "ошибка добавления друга: "+err.Error())
		tpl.Render(w, service.TplNameError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/view/%d", userID), http.StatusTemporaryRedirect)
	return
}
