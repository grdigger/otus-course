package service

import (
	"github.com/grdigger/otus-course/internal/model"
	"github.com/grdigger/otus-course/internal/repository"
)

type UserTemplate struct {
	user              *model.User
	tpl               Template
	userRepo          repository.UserRepository
	userInterestsRepo repository.InterestRepository
	cityHelper        model.HelperInterface
	genderHelper      model.HelperInterface
	interestsHelper   model.HelperInterface
}

func (ut *UserTemplate) Fill() {

}
