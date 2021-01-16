package handler

import (
	"github.com/grdigger/otus-course/internal/service"
	"net/http"
)

type Index struct {
	session *service.SessionService
}

func NewIndex(session *service.SessionService) *Index {
	return &Index{session: session}
}

func (h *Index) Handle(w http.ResponseWriter, r *http.Request) {
	tpl := service.NewTemplate()
	isAuth, err := h.session.IsUserLogged(r)
	if err != nil {
		tpl.AddVar("error", "server error: "+err.Error())
		tpl.Render(w, service.TplNameError)
		return
	}

	if !isAuth {
		tpl.Render(w, service.TplNameLogin)
	} else {
		http.Redirect(w, r, "/personal", http.StatusTemporaryRedirect)
		return
	}
}
