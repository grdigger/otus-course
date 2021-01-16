package handler

import (
	"github.com/grdigger/otus-course/internal/model"
	"github.com/grdigger/otus-course/internal/repository"
	"github.com/grdigger/otus-course/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type Edit struct {
	session           *service.SessionService
	userRepo          repository.UserRepository
	userInterestsRepo repository.InterestRepository
	cityHelper        model.HelperInterface
	genderHelper      model.HelperInterface
	interestsHelper   model.HelperInterface
}

func NewEdit(session *service.SessionService,
	userRepo repository.UserRepository,
	userInterestsRepo repository.InterestRepository,
	cityHelper model.HelperInterface,
	genderHelper model.HelperInterface,
	interestsHelper model.HelperInterface) *Edit {
	return &Edit{session: session,
		userRepo:          userRepo,
		userInterestsRepo: userInterestsRepo,
		cityHelper:        cityHelper,
		genderHelper:      genderHelper,
		interestsHelper:   interestsHelper}
}

func (h *Edit) Handle(w http.ResponseWriter, r *http.Request) {
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

	tpl.AddVar("ID", user.GetID())
	tpl.AddVar("Name", user.GetName())
	tpl.AddVar("Surname", user.GetSurname())
	tpl.AddVar("Age", user.GetAge())
	tpl.AddVar("GenderId", user.GetGenderId())
	tpl.AddVar("CityId", user.GetCityId())

	cityName, _ := h.cityHelper.GetByID(user.GetCityId())
	genderName, _ := h.genderHelper.GetByID(user.GetGenderId())

	tpl.AddVar("CityName", cityName)
	tpl.AddVar("GenderName", genderName)

	interests, err := h.userInterestsRepo.GetByUserId(user.GetID())
	if err != nil {
		tpl.AddVar("error", "ошибка получения интересов пользователя: "+err.Error())
		tpl.Render(w, service.TplNameLogin)
	}
	intrs := make([]string, 0)
	for _, v := range interests {
		intrs = append(intrs, strconv.Itoa(v))
	}
	tpl.AddVar("Interests", strings.Join(intrs, ","))
	tpl.AddVar("GenderList", h.genderHelper.GetAll())
	tpl.AddVar("InterestsList", h.interestsHelper.GetAll())
	tpl.AddVar("CityList", h.cityHelper.GetAll())

	tpl.Render(w, service.TplNamePersonalEdit)
	return
}
