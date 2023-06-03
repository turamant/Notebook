package main

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"askvart.ru/snippetbox/pkg/models"
)

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
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := "How about..."
	content := "How about...\nHow about Eddy,\nWhat can i do for you?\n\n- Askvart Victor"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

}

func (app *application) downloadHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Clean("./ui/static/myresume.pdf")
	http.ServeFile(w, r, path)
}