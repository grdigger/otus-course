package main

import (
	"github.com/gorilla/mux"
	"github.com/grdigger/otus-course/handler"
	"github.com/grdigger/otus-course/inits"
	"github.com/grdigger/otus-course/internal/repository"
	"log"
	"net/http"
)

func main() {
	db := inits.NewDB()
	defer db.Close()

	sess := inits.InitSession()
	userRepository := repository.NewUser(db)
	userInterestsRepository := repository.NewInterest(db)
	friendRepository := repository.NewFriend(db)
	cityHelper, genderHelper, interestHelper := inits.InitHelpers(db)

	r := mux.NewRouter()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/static/"))))

	indexHandler := handler.NewIndex(sess)
	loginHandler := handler.NewLogin(sess, userRepository)
	logoutHandler := handler.NewLogout(sess, userRepository)
	personalHandler := handler.NewPersonal(sess, userRepository, userInterestsRepository, cityHelper, genderHelper, interestHelper)
	editHandler := handler.NewEdit(sess, userRepository, userInterestsRepository, cityHelper, genderHelper, interestHelper)
	saveHandler := handler.NewSave(sess, userRepository, userInterestsRepository, cityHelper, genderHelper, interestHelper)
	viewHandler := handler.NewView(sess, userRepository, userInterestsRepository, friendRepository, cityHelper, genderHelper, interestHelper)
	friendAddHandler := handler.NewFriendAdd(sess, friendRepository)
	friendRemoveHandler := handler.NewFriendRemove(sess, friendRepository)
	friendHandler := handler.NewFriendList(sess, friendRepository)
	regHandler := handler.NewReg(sess, userRepository, userInterestsRepository, cityHelper, genderHelper, interestHelper)
	searchHandler := handler.NewSearch(sess, friendRepository)

	r.HandleFunc("/", indexHandler.Handle).Methods("GET")
	r.HandleFunc("/login", loginHandler.Handle).Methods("POST")
	r.HandleFunc("/logout", logoutHandler.Handle).Methods("GET")
	r.HandleFunc("/personal", personalHandler.Handle).Methods("POST", "GET")
	r.HandleFunc("/edit", editHandler.Handle).Methods("GET")
	r.HandleFunc("/save", saveHandler.Handle).Methods("POST")
	r.HandleFunc("/view/{id}", viewHandler.Handle).Methods("GET")
	r.HandleFunc("/addfriend/{id}", friendAddHandler.Handle).Methods("GET")
	r.HandleFunc("/removefriend/{id}", friendRemoveHandler.Handle).Methods("GET")
	r.HandleFunc("/friends", friendHandler.Handle).Methods("GET")
	r.HandleFunc("/reg", regHandler.Handle).Methods("GET", "POST")
	r.HandleFunc("/search", searchHandler.Handle).Methods("GET", "POST")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", r))
}
