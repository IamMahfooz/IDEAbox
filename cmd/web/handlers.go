package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/pkg/forms"
	"snippetbox/pkg/models"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
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
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.Maxlength("title", 100)
	form.PermittedValues("expire", "365", "7", "1")
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}
	byy := app.session.Get(r, "authenticatedUserName").(string)
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"), byy)
	if err != nil {
		app.serverError(w, err)
	}
	app.session.Put(r, "flash", "Snippet successfullly created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.Maxlength("name", 225)
	form.Maxlength("email", 225)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Adress is already in use")
			app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.session.Put(r, "flash", "Your signup was successful. Please logn in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusInternalServerError)
		return
	}
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.session.Put(r, "authenticatedUserID", id)
	app.session.Put(r, "authenticatedUserName", form.Get("email"))
	if app.session.Get(r, "rlink") != nil {
		rlink := app.session.Get(r, "rlink").(string)
		app.session.Remove(r, "rlink")
		http.Redirect(w, r, rlink, http.StatusSeeOther)

	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Remove(r, "authenticatedUserName")
	app.session.Put(r, "flash", "You'have been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//func (app *application,m *U) forget(w http.ResponseWriter, r *http.Request) {
//	idd := app.session.Get(r, "authenticatedUserID").(string)
//	hashed_pass, err := bcrypt.GenerateFromPassword([]byte(password), 12).(string)
//	if err != nil {
//		return err
//	}
//	stmt := `UPDATE users SET hashed_password = hashed_pass WHERE id =idd;`
//	_, err = m.DB.Exec(stmt)
//}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
func (app *application) about(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "about.page.tmpl", &templateData{})
}
func (app *application) faceism(w http.ResponseWriter, r *http.Request) {
	fmt.Println(app.session.Get(r, "authenticatedUserName"))
	app.render(w, r, "upload.page.tmpl", &templateData{Count: app.count})

}
