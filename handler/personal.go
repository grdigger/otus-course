package handler

import (
	"fmt"
	"github.com/grdigger/otus-course/internal/model"
	"github.com/grdigger/otus-course/internal/repository"
	"github.com/grdigger/otus-course/internal/service"
	l "github.com/sirupsen/logrus"
	"net/http"
)

type Personal struct {
	session           *service.SessionService
	userRepo          repository.UserRepository
	userInterestsRepo repository.InterestRepository
	cityHelper        model.HelperInterface
	genderHelper      model.HelperInterface
	interestsHelper   model.HelperInterface
}

func NewPersonal(session *service.SessionService,
	userRepo repository.UserRepository,
	userInterestsRepo repository.InterestRepository,
	cityHelper model.HelperInterface,
	genderHelper model.HelperInterface,
	interestsHelper model.HelperInterface) *Personal {
	return &Personal{session: session,
		userRepo:          userRepo,
		userInterestsRepo: userInterestsRepo,
		cityHelper:        cityHelper,
		genderHelper:      genderHelper,
		interestsHelper:   interestsHelper}
}

func (h *Personal) Handle(w http.ResponseWriter, r *http.Request) {
	tpl := service.NewTemplate()
	us := service.NewUserSession(tpl, h.session, w, r)
	user, err := us.User()
	if err != nil {
		l.Errorf("ошика  получения данных сессии: %s ", err.Error())
		tpl.Render(w, service.TplNameLogin)
		return
	}

	tpl.AddVar("ID", user.GetID())
	tpl.AddVar("Name", user.GetName())
	tpl.AddVar("Surname", user.GetSurname())
	tpl.AddVar("Age", user.GetAge())
	if user.GetGenderId() > 0 {
		g, haveValue := h.genderHelper.GetByID(user.GetGenderId())
		if !haveValue {
			l.Errorf("ошибка получения значения пола для id : %d", user.GetGenderId())
			tpl.AddVar("error", fmt.Sprintf("ошибка получения значения пола для id : %d", user.GetGenderId()))
		} else {
			tpl.AddVar("Gender", g)
		}
	}
	if user.GetCityId() > 0 {
		c, haveValue := h.cityHelper.GetByID(user.GetCityId())
		if !haveValue {
			l.Errorf("ошибка получения значения города для id : %d", user.GetGenderId())
			tpl.AddVar("error", fmt.Sprintf("ошибка получения значения города для id : %d", user.GetGenderId()))
		} else {
			tpl.AddVar("City", c)
		}
	}

	interests, err := h.userInterestsRepo.GetByUserId(user.GetID())
	if err != nil {
		l.Errorf("ошибка получения интересов пользователя: " + err.Error())
		tpl.AddVar("error", "ошибка получения интересов пользователя: "+err.Error())
		tpl.Render(w, service.TplNameLogin)
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
	tpl.AddVar("IsOwner", true)

	tpl.Render(w, service.TplNamePersonal)
	return
}
