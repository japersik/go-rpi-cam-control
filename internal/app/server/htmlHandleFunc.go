package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

func (s *server) handleAlbumHTML() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil || id < 1 {
			id = 1
			http.Redirect(w, r, strconv.Itoa(id), http.StatusTemporaryRedirect)
		}
		fmt.Println(id)
		//files, _ := ioutil.ReadDir("private/static/img/")
		files, _ := s.service.GetPhotoNames(12, id)
		fmt.Println(files)
		t, err := template.ParseFiles("html/template/album.html", "html/template/header.html", "html/template/footer.html")
		if err != nil {
			s.logger.Error(err.Error())
			return
		}
		err = t.ExecuteTemplate(w, "album", files)
	}
}
func (s *server) handleHomeHTML() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("html/template/index.html", "html/template/header.html", "html/template/footer.html")
		if err != nil {
			s.logger.Error(err.Error())
			return
		}
		err = t.ExecuteTemplate(w, "index", struct {
		}{})
		if err != nil {
			s.logger.Error(err.Error())
			return
		}
	}
}
func (s *server) handleRegisterHTML() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("html/template/register.html", "html/template/header.html", "html/template/footer.html")
		if err != nil {
			s.logger.Error(err.Error())
			return
		}
		err = t.ExecuteTemplate(w, "register", struct{}{})
	}
}
func (s *server) handleLoginHTML() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("html/template/login.html", "html/template/header.html", "html/template/footer.html")
		if err != nil {
			fmt.Fprintf(os.Stdout, err.Error())
			return
		}
		err = t.ExecuteTemplate(w, "login", struct{}{})
	}
}
