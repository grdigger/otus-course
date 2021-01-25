package handler

import (
	"net/http"
	"strconv"

	"github.com/grdigger/otus-course/internal/model"
	"github.com/grdigger/otus-course/internal/repository"
	"github.com/grdigger/otus-course/internal/service"
	l "github.com/sirupsen/logrus"
)

type Reg struct {
	session           *service.SessionService
	userRepo          repository.UserRepository
	userInterestsRepo repository.InterestRepository
	cityHelper        model.HelperInterface
	genderHelper      model.HelperInterface
	interestsHelper   model.HelperInterface
}

func NewReg(session *service.SessionService,
	userRepo repository.UserRepository,
	userInterestsRepo repository.InterestRepository,
	cityHelper model.HelperInterface,
	genderHelper model.HelperInterface,
	interestsHelper model.HelperInterface,
) *Reg {
	return &Reg{session: session,
		userInterestsRepo: userInterestsRepo,
		userRepo:          userRepo,
		cityHelper:        cityHelper,
		genderHelper:      genderHelper,
		interestsHelper:   interestsHelper,
	}
}

func (h *Reg) Handle(w http.ResponseWriter, r *http.Request) {

	tpl := service.NewTemplate()

	if r.Method == http.MethodGet {
		tpl.AddVar("GenderList", h.genderHelper.GetAll())
		tpl.AddVar("InterestsList", h.interestsHelper.GetAll())
		tpl.AddVar("CityList", h.cityHelper.GetAll())
		tpl.Render(w, service.TplNameReg)
		return
	}
	if r.Method == http.MethodPost {
		r.ParseForm()
		interests := r.Form["interests"]
		email := r.FormValue("email")
		password := r.FormValue("password")
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
		_, ok := h.genderHelper.GetAll()[gender]
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

		userID, err := h.userRepo.Add(email, password, name, surname, ageVal, genderVal, cityVal)
		if err != nil {
			tpl.AddVar("error", "ошибка сохранения пользователя"+err.Error())
			tpl.Render(w, service.TplNameError)
		}

		newUser := new(model.User)
		newUser.ID = userID
		newUser.Name = &name
		newUser.Surname = &surname
		newUser.Age = &ageVal
		newUser.GenderId = &genderVal
		newUser.CityId = &cityVal

		result, err := h.session.Save(newUser, r, w)
		if err != nil {
			l.Errorf("error saving session: " + err.Error())
			tpl.AddVar("error", "ошибка сохранения сессии"+err.Error())
			tpl.Render(w, service.TplNameError)
			return
		}
		if !result {
			l.Errorf("error writing session %w", newUser)
			tpl.AddVar("error", "сессия не записана")
			tpl.Render(w, service.TplNameError)
			return
		}
		_, err = h.userInterestsRepo.UpdateInterestsById(userID, interests)
		if err != nil {
			l.Errorf("error saving user interests" + err.Error())
			tpl.AddVar("error", "ошибка сохранения интересов пользователя"+err.Error())
			tpl.Render(w, service.TplNameError)
		}
		http.Redirect(w, r, "/personal", http.StatusTemporaryRedirect)
	}
	http.Redirect(w, r, "/reg", http.StatusTemporaryRedirect)
	return
}
