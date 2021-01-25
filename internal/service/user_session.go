package service

import (
	"fmt"
	"github.com/grdigger/otus-course/internal/model"
	"net/http"
)

type UserSession struct {
	tpl     *Template
	session *SessionService
	r       *http.Request
	w       http.ResponseWriter
}

func NewUserSession(tpl *Template, session *SessionService, w http.ResponseWriter, r *http.Request) *UserSession {
	return &UserSession{tpl: tpl, session: session, r: r, w: w}
}

func (us *UserSession) User() (*model.User, error) {
	userData, err := us.session.UserSessions(us.r)
	if err != nil {
		return nil, err
	}
	user, ok := userData.(model.User)
	if !ok {
		return nil, fmt.Errorf("error reading session %+v", userData)
	}
	return &user, nil
}
