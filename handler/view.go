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

type View struct {
	session           *service.SessionService
	userRepo          repository.UserRepository
	userInterestsRepo repository.InterestRepository
	friendRepo        repository.FriendRepository
	cityHelper        model.HelperInterface
	genderHelper      model.HelperInterface
	interestsHelper   model.HelperInterface
}

func NewView(session *service.SessionService,
	userRepo repository.UserRepository,
	userInterestsRepo repository.InterestRepository,
	friendRepo repository.FriendRepository,
	cityHelper model.HelperInterface,
	genderHelper model.HelperInterface,
	interestsHelper model.HelperInterface) *View {
	return &View{session: session,
		userRepo:          userRepo,
		userInterestsRepo: userInterestsRepo,
		friendRepo:        friendRepo,
		cityHelper:        cityHelper,
		genderHelper:      genderHelper,
		interestsHelper:   interestsHelper}
}

func (h *View) Handle(w http.ResponseWriter, r *http.Request) {

	tpl := service.NewTemplate()
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		tpl.AddVar("error", "error parsing viewID: "+err.Error())
		tpl.Render(w, service.TplNameError)
		return
	}
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		tpl.AddVar("error", "ошибка получения данных пользователяD: "+err.Error())
		tpl.Render(w, service.TplNameError)
		return
	}
	if user.IsEmpty() {
		tpl.AddVar("error", "такого пользователя не существует")
		tpl.Render(w, service.TplNameError)
		return
	}

	currentUserID := 0
	currentUserData, _ := h.session.UserSessions(r)
	currentUser, ok := currentUserData.(model.User)
	if ok {
		currentUserID = currentUser.GetID()
	}

	tpl.AddVar("ID", user.GetID())
	tpl.AddVar("Name", user.GetName())
	tpl.AddVar("Surname", user.GetSurname())
	tpl.AddVar("Age", user.GetAge())
	if user.GetGenderId() > 0 {
		g, haveValue := h.genderHelper.GetByID(user.GetGenderId())
		if !haveValue {
			tpl.AddVar("error", fmt.Sprintf("ошибка получения значения пола для id : %d", user.GetGenderId()))
		} else {
			tpl.AddVar("Gender", g)
		}
	}
	if user.GetCityId() > 0 {
		c, haveValue := h.cityHelper.GetByID(user.GetCityId())
		if !haveValue {
			tpl.AddVar("error", fmt.Sprintf("ошибка получения значения города для id : %d", user.GetGenderId()))
		} else {
			tpl.AddVar("City", c)
		}
	}

	interests, err := h.userInterestsRepo.GetByUserId(user.GetID())
	if err != nil {
		tpl.AddVar("error", "ошибка получения интересов пользователя: "+err.Error())
		tpl.Render(w, service.TplNameError)
	}
	intrs := make([]string, 0)
	for _, v := range interests {
		val, haveValue := h.interestsHelper.GetByID(v)
		if !haveValue {
			tpl.AddVar("error", fmt.Sprintf("ошибка получения значения интересоа пользователя %d: ", v))
		}
		intrs = append(intrs, val)
	}
	tpl.AddVar("Interests", intrs)

	friendList, err := h.friendRepo.GetFriends(currentUserID)
	if err != nil {
		tpl.AddVar("error", "ошибка получения друга: "+err.Error())
		tpl.Render(w, service.TplNameError)
		return
	}
	isFrended := false
	for friendID, _ := range friendList {
		if userID == friendID {
			isFrended = true
			break
		}
	}

	if isFrended {
		tpl.AddVar("link", "removefriend")
		tpl.AddVar("linkName", "Удалить из друзей")
	} else {
		tpl.AddVar("link", "addfriend")
		tpl.AddVar("linkName", "Добавить в друзья")
	}
	tpl.Render(w, service.TplNamePersonal)

	return
}
