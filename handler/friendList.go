package handler

import (
	"net/http"

	"github.com/grdigger/otus-course/internal/model"
	"github.com/grdigger/otus-course/internal/repository"
	"github.com/grdigger/otus-course/internal/service"
)

type FriendList struct {
	session    *service.SessionService
	friendRepo repository.FriendRepository
}

func NewFriendList(session *service.SessionService,
	friendRepo repository.FriendRepository,
) *FriendList {
	return &FriendList{session: session,
		friendRepo: friendRepo,
	}
}

func (h *FriendList) Handle(w http.ResponseWriter, r *http.Request) {
	tpl := service.NewTemplate()

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
	friendIDS, err := h.friendRepo.GetFriends(user.GetID())
	if err != nil {
		tpl.AddVar("error", "ошибка получения списка друзей: "+err.Error())
		tpl.Render(w, service.TplNameError)
		return
	}

	tpl.AddVar("friendList", friendIDS)
	tpl.Render(w, service.TplNameFriend)
	return
}
