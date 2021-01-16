package service

import (
	"encoding/gob"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/grdigger/otus-course/internal/model"
	"net/http"
)

type SessionService struct {
	isAuth      bool
	email       string
	Name        string
	Surname     string
	store       *sessions.CookieStore
	sessionData *sessions.Session
}

const (
	SessionName              = "go-sessions-otus-course"
	SessionAuthenticationKey = "u$/nV]4].aZ;/9>*"
	SessionEncryptionKey     = "89tf{'},PY3vdgQX"
)

func NewSession() (*SessionService, error) {
	gob.Register(model.User{})

	sessionStore := sessions.NewCookieStore([]byte(SessionAuthenticationKey))
	return &SessionService{store: sessionStore}, nil
}

// GetSession to return the session
func (s *SessionService) GetSession(r *http.Request) (*sessions.Session, error) {
	var err error
	s.sessionData, err = s.store.Get(r, "otus-course-session")
	if err != nil {
		return nil, err
	}
	return s.sessionData, nil
}

func (s *SessionService) UserSessions(r *http.Request) (interface{}, error) {
	session, err := s.GetSession(r)
	if err != nil {
		return nil, err
	}
	userData := session.Values["userData"]
	return userData, nil
}

func (s SessionService) IsUserLogged(r *http.Request) (bool, error) {
	userData, err := s.UserSessions(r)
	if err != nil {
		return false, err
	}
	if userData == nil {
		return false, nil
	}
	u, ok := userData.(model.User)
	if !ok {
		return false, nil
	}
	if u.IsEmpty() {
		return false, nil
	}
	return true, nil
}

func (s *SessionService) Save(user *model.User, r *http.Request, w http.ResponseWriter) (bool, error) {
	session, err := s.GetSession(r)
	if err != nil {
		return false, err
	}
	session.Values["userData"] = user
	err = session.Save(r, w)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	return true, nil
}
