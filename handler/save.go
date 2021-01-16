package handler

import (
	"net/http"
	"strconv"

	"github.com/grdigger/otus-course/internal/model"
	"github.com/grdigger/otus-course/internal/repository"
	"github.com/grdigger/otus-course/internal/service"
	l "github.com/sirupsen/logrus"
)

type Save struct {
	session           *service.SessionService
	userRepo          repository.UserRepository
	userInterestsRepo repository.InterestRepository
	cityHelper        model.HelperInterface
	genderHelper      model.HelperInterface
	interestsHelper   model.HelperInterface
}

func NewSave(session *service.SessionService,
	userRepo repository.UserRepository,
	userInterestsRepo repository.InterestRepository,
	cityHelper model.HelperInterface,
	genderHelper model.HelperInterface,
	interestsHelper model.HelperInterface) *Save {
	return &Save{session: session,
		userRepo:          userRepo,
		userInterestsRepo: userInterestsRepo,
		cityHelper:        cityHelper,
		genderHelper:      genderHelper,
		interestsHelper:   interestsHelper}
}

func (h *Save) Handle(w http.ResponseWriter, r *http.Request) {
	tpl := service.NewTemplate()

	err := r.ParseForm()
	if err != nil {
		l.Errorf("error parsing form" + err.Error())
		tpl.AddVar("error", "error parsing form: "+err.Error())
		tpl.Render(w, service.TplNameError)
		return
	}
	userData, err := h.session.UserSessions(r)
	if err != nil {
		l.Errorf("ошика  чтения сессии: %s ", err.Error())
		tpl.AddVar("error", "server error: "+err.Error())
		tpl.Render(w, service.TplNameLogin)
		return
	}
	user, ok := userData.(model.User)
	if !ok {
		l.Errorf("ошика  преообразования сессии: %w ", userData)
		tpl.AddVar("error", "ошибка чтения сессии")
		tpl.Render(w, service.TplNameError)
		return
	}
	if user.IsEmpty() {
		tpl.AddVar("userNotFound", "Пользователь не найден")
		tpl.Render(w, service.TplNameLogin)
		return
	}
	err = r.ParseForm()
	if err != nil {
		l.Errorf("error parsing form" + err.Error())
		tpl.AddVar("error", "error parsing form: "+err.Error())
		tpl.Render(w, service.TplNameError)
		return
	}

	interests := r.Form["interests"]
	name := r.FormValue("name")
	surname := r.FormValue("surname")
	age := r.FormValue("age")
	gender := r.FormValue("gender")
	city := r.FormValue("city")

	ageVal, err := strconv.Atoi(age)
	if err != nil {
		tpl.AddVar("error", "поле Возраст должно быть цифрой")
		tpl.Render(w, service.TplNameError)
		return
	}

	genderVal, err := strconv.Atoi(gender)
	if err != nil {
		tpl.AddVar("error", "поле Пол не корректно")
		tpl.Render(w, service.TplNameError)
		return
	}
	_, ok = h.genderHelper.GetAll()[gender]
	if !ok {
		tpl.AddVar("error", "поле Пол не корректно")
		tpl.Render(w, service.TplNameError)
		return
	}

	cityVal, err := strconv.Atoi(city)
	if err != nil {
		tpl.AddVar("error", "поле Город должно быть цифрой")
		tpl.Render(w, service.TplNameError)
		return
	}
	_, ok = h.cityHelper.GetAll()[city]
	if !ok {
		tpl.AddVar("error", "поле Город не корректно")
		tpl.Render(w, service.TplNameError)
		return
	}

	_, err = h.userRepo.UpdateById(user.GetID(), name, surname, ageVal, genderVal, cityVal)
	if err != nil {
		tpl.AddVar("error", "ошибка сохранения пользователя"+err.Error())
		tpl.Render(w, service.TplNameError)
	}
	newUser := new(model.User)
	newUser.ID = user.GetID()
	newUser.Name = &name
	newUser.Surname = &surname
	newUser.Age = &ageVal
	newUser.GenderId = &genderVal
	newUser.CityId = &cityVal

	result, err := h.session.Save(newUser, r, w)
	if err != nil {
		l.Errorf("ошибка сохранения сессии" + err.Error())
		tpl.AddVar("error", "ошибка сохранения сессии"+err.Error())
		tpl.Render(w, service.TplNameError)
		return
	}
	if !result {
		tpl.AddVar("error", "сессия не записана")
		tpl.Render(w, service.TplNameError)
		return
	}
	_, err = h.userInterestsRepo.UpdateInterestsById(user.GetID(), interests)
	if err != nil {
		l.Errorf("ошибка сохранения интересов пользователя" + err.Error())
		tpl.AddVar("error", "ошибка сохранения интересов пользователя"+err.Error())
		tpl.Render(w, service.TplNameError)
	}
	http.Redirect(w, r, "/personal", http.StatusTemporaryRedirect)
	return
}
