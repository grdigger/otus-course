package inits

import (
	"database/sql"
	"github.com/grdigger/otus-course/internal/model"
	"github.com/grdigger/otus-course/internal/repository"
	"log"
)

func InitHelpers(db *sql.DB) (model.HelperInterface, model.HelperInterface, model.HelperInterface) {
	cityHelper, err := repository.NewHelperRepo(db).GelAll("city")
	if err != nil {
		log.Fatal("error getting data for helper: city")
	}
	genderHelper, err := repository.NewHelperRepo(db).GelAll("gender")
	if err != nil {
		log.Fatal("error getting data for helper: gender")
	}
	interestHelper, err := repository.NewHelperRepo(db).GelAll("interest")
	if err != nil {
		log.Fatal("error getting data for helper: interest")
	}
	return cityHelper, genderHelper, interestHelper
}
