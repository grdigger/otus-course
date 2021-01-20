package service

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type page struct {
	Error string
	Vars  map[string]interface{}
}

func (p *page) AddVar(key string, val interface{}) {
	p.Vars[key] = val
}

const TemplateDir = "templates"

func NewPage() *page {
	p := &page{}
	p.Vars = make(map[string]interface{})
	return p
}

const TplNameError = "error"
const TplNameLogin = "main"
const TplNameFriend = "friend"
const TplNamePersonal = "personal"
const TplNamePersonalEdit = "personal_edit"
const TplNameReg = "reg"
const TplNameSearch = "search"

type Template struct {
	page *page
}

func NewTemplate() *Template {
	t := &Template{}
	t.page = NewPage()
	return t
}

func (tpl *Template) AddVar(key string, val interface{}) {
	tpl.page.AddVar(key, val)
}

func (tpl *Template) Render(w http.ResponseWriter, tmpl string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.ParseFiles(
		fmt.Sprintf("%s/%s/%s.html", dir, TemplateDir, tmpl),
		fmt.Sprintf("%s/%s/%s.html", dir, TemplateDir, "header"),
		fmt.Sprintf("%s/%s/%s.html", dir, TemplateDir, "error"),
		fmt.Sprintf("%s/%s/%s.html", dir, TemplateDir, "personal"),
		fmt.Sprintf("%s/%s/%s.html", dir, TemplateDir, "personal_edit"),
		fmt.Sprintf("%s/%s/%s.html", dir, TemplateDir, "friend"),
		fmt.Sprintf("%s/%s/%s.html", dir, TemplateDir, "footer"))
	if err != nil {
		log.Fatal(err.Error())
	}
	err = t.Execute(w, tpl.page)
	if err != nil {
		log.Fatal(err.Error())
	}
}
