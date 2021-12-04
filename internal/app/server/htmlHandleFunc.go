package server

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

func (s *server) handleAlbumHTML() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		files,_:= ioutil.ReadDir("private/static/img/")
		t, err := template.ParseFiles("html/template/album.html","html/template/header.html","html/template/footer.html")
		if err != nil {
			fmt.Fprintf(os.Stdout, err.Error())
			return
		}
		err = t.ExecuteTemplate(w,"album", files)
	}
}
func (s *server) handleHomeHTML() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("html/template/index.html","html/template/header.html","html/template/footer.html")
		if err != nil {
			fmt.Fprintf(os.Stdout, err.Error())
			return
		}
		err = t.ExecuteTemplate(w,"index", struct {
		}{})
		if err != nil {
			fmt.Fprintf(os.Stdout, err.Error())
			return
		}
	}
}
func (s *server) handleRegisterHTML() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("html/template/register.html","html/template/header.html","html/template/footer.html")
		if err != nil {
			fmt.Fprintf(os.Stdout, err.Error())
			return
		}
		err = t.ExecuteTemplate(w,"register", struct {}{})
	}
}