package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
)
var tpl *template.Template
var userData = make(map[string]User)

type PageData struct {
	EmailInvalid string
	PassInvalid  string
}
type User struct {
	Name     string
	Email    string
	Password string
}
type Session struct {
	SessionId string
}

var RandumNumber = rand.Intn(200)

func main() {
	tpl, _ = template.ParseGlob("template/*.html")
	http.HandleFunc("/", indexHandle)
	http.HandleFunc("/loginpost", postmethod1)
	http.HandleFunc("/signuppost", signupmethod)
	http.HandleFunc("/home", homefunc)
	http.HandleFunc("/login", loginfunc)
	http.HandleFunc("/signup", handlefuncSignup)
	http.HandleFunc("/logout", logoutfunc)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8080", nil)
}
func indexHandle(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("logincookie")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)

}
func homefunc(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("logincookie")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	tpl.ExecuteTemplate(w, "index.html", nil)

}
func loginfunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	cookie, err := r.Cookie("logincookie")
	if err == nil && cookie.Value == strconv.Itoa(RandumNumber) {

		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return

	}

	tpl.ExecuteTemplate(w, "login.html", nil)
}

func postmethod1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	cookie, err := r.Cookie("logincookie")
	if err != nil {
		fmt.Println(err)

	} else if cookie.Value != "" {
		http.Redirect(w, r, "/loginpost", http.StatusSeeOther)
		return
	}
	// fmt.Println(email, password)
	email := r.FormValue("emailLogin")
	password := r.FormValue("passwordLogin")

	user, ok := userData[email]

	if email == "" {
		var n = PageData{EmailInvalid: "Email is Invalid"}
		tpl.ExecuteTemplate(w, "login.html", n)
		fmt.Println("EmailEmpty")
		return
	} else if password == "" {
		var n = PageData{PassInvalid: "Password Cannot be Empty"}
		tpl.ExecuteTemplate(w, "login.html", n)
		fmt.Println("PasswordEmpty")
		return
	}
	if ok && password != user.Password {
		var n = PageData{PassInvalid: "Invalid Credentials"}
		tpl.ExecuteTemplate(w, "login.html", n)
		fmt.Println("PasswordEmpty")
		return

	}
	if ok && password == user.Password {

		CookieForLogin := &http.Cookie{}
		CookieForLogin.Name = "logincookie"
		
		CookieForLogin.Value = strconv.Itoa(RandumNumber)
		CookieForLogin.MaxAge = 300
		CookieForLogin.Path = "/"
		http.SetCookie(w, CookieForLogin)
		fmt.Println(CookieForLogin)
		tpl.ExecuteTemplate(w, "index.html", CookieForLogin.Value)

	} else {
		// tpl.ExecuteTemplate(w, "login.html", "Invalid Credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}
func handlefuncSignup(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "signup.html", nil)
}
func signupmethod(w http.ResponseWriter, r *http.Request) {
	Name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" {
		tpl.ExecuteTemplate(w, "signup.html", "EmailInvalid")
		return
	}
	if password == "" {
		tpl.ExecuteTemplate(w, "signup.html", "PasswordInvalid")
		return
	}
	userData[email] = User{Name: Name,
		Password: password,
		Email:    email,
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)

	fmt.Printf("%+v", userData)
}
func logoutfunc(w http.ResponseWriter, r *http.Request) {
	cookielogout := http.Cookie{}
	cookielogout.Name = "logincookie"
	cookielogout.Value = ""
	cookielogout.MaxAge = -1
	cookielogout.Path = "/"
	http.SetCookie(w, &cookielogout)

	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	http.Redirect(w, r, "/login", http.StatusSeeOther)

	cookie, err := r.Cookie("logincookie")
	if err != nil {
		fmt.Println(err)

	} else if cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}
