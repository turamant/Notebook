package main

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	

	"askvart.ru/snippetbox/pkg/forms"
	"askvart.ru/snippetbox/pkg/models"
)


func (app *application) about(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/about"{
		app.notFound(w)
		return
	}

	app.render(w, r, "about.page.tmpl", nil)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Latest()
		if err != nil {
		app.serverError(w, err)
		return
	}
	
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
	}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}


func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil{
		app.clientError(w, http.StatusBadRequest)
		return

	}
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")
	
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
	return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

}

func (app *application) downloadHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Clean("./ui/static/myresume.pdf")
	http.ServeFile(w, r, path)
}