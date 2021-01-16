package inits

import (
	"github.com/grdigger/otus-course/internal/service"
	"log"
)

func InitSession() *service.SessionService {
	sess, err := service.NewSession()
	if err != nil {
		log.Fatal("session init error: " + err.Error())
	}
	return sess
}
