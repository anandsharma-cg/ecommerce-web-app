package authentication

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ecommerce/internal/core/session"
	"github.com/ecommerce/internal/services/user"

	"github.com/gorilla/mux"
)

const (
	authBasePath = "auth"
	apiVersion   = "prod"
	apiBasePath  = "api"
)

// SetupRoutes :
func SetupAuthRoutes(r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/%s/login", apiBasePath, authBasePath), loginHandler)
	r.HandleFunc(fmt.Sprintf("/%s/%s/register", apiBasePath, authBasePath), registerHandler)

	// -------------------------PROD----------------------
	r.HandleFunc(fmt.Sprintf("/%s/%s/login", apiVersion, authBasePath), loginProdHandler)
	r.HandleFunc(fmt.Sprintf("/%s/%s/register", apiVersion, authBasePath), registerProdHandler)
	r.HandleFunc(fmt.Sprintf("/%s/%s/logout", apiVersion, authBasePath), logoutProdHandler)
}

func registerProdHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the template file (adjust path if necessary)
	tmpl, err := template.ParseFiles("template/register.html")
	if err != nil {
		log.Println("Template parsing error:", err)
		http.Error(w, "Error loading register page", http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case http.MethodGet:
		// Execute the template, sending data if needed (or nil if not)
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println("Template execution error:", err)
			http.Error(w, "Error rendering register page", http.StatusInternalServerError)
		}
	case http.MethodPost:

		sess, err := session.GetSessionFromContext(r)
		if sess == nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = r.ParseForm()
		if err != nil {
			err := errors.New("Error parsing form data")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// extracting data of form values
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		// email and pass empty validation
		if email == "" || password == "" {
			err := errors.New("email and password are required")
			log.Println(err)
			tmpl.Execute(w, map[string]string{"Error": err.Error()})
			return
		}

		// pass and confirm pass validation
		if password != confirmPassword {
			err := errors.New("password and confirm password is not same")
			log.Println(err)
			tmpl.Execute(w, map[string]string{"Error": err.Error()})
			return
		}

		//register user
		newUser := user.User{Email: email, Password: password}
		res, err := registerUserService(newUser)

		if err != nil {
			log.Println(res, ": ", err)
			tmpl.Execute(w, map[string]string{"Error": err.Error()})
			return
		}

		//storing user in session
		user, err, res := user.GetUserByEmailService(email)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), res)
			return
		}

		userObj := session.User{
			UserID: user.UserID, Email: user.Email,
			Password: user.Password, IsAdmin: user.IsAdmin}

		sess.Values["user"] = &userObj
		sess.Values["userId"] = user.UserID
		err = sess.Save(r, w)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// If register is successful, redirect
		// Redirect to dashboard page on successful login
		log.Println("redirect-to-dashboard")
		http.Redirect(w, r, "/prod/users/dashboard", http.StatusSeeOther) // 302 Found
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func loginProdHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the template file (adjust path if necessary)
	tmpl, err := template.ParseFiles("template/login.html")
	if err != nil {
		log.Println("Template parsing error:", err)
		http.Error(w, "Error loading login page", http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Execute the template, sending data if needed (or nil if not)
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println("Template execution error:", err)
			http.Error(w, "Error rendering login page", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		sess, err := session.GetSessionFromContext(r)
		if sess == nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		// extracting data of form values
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Simple validation
		if email == "" || password == "" {
			err := errors.New("email and password are required")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//login user
		existingUser := user.User{Email: email, Password: password}
		res, err := loginUserService(existingUser)

		if err != nil {
			log.Println(err)
			tmpl.Execute(w, map[string]string{"Error": err.Error()})
			return
		}

		//storing user in session
		user, err, res := user.GetUserByEmailService(email)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), res)
			return
		}

		userObj := session.User{
			UserID: user.UserID, Email: user.Email,
			Password: user.Password, IsAdmin: user.IsAdmin}

		sess.Values["user"] = &userObj
		sess.Values["userId"] = user.UserID
		err = sess.Save(r, w)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// If login is successful, redirect
		// Redirect to dashboard page on successful login
		log.Println("redirect-to-dashboard")
		http.Redirect(w, r, "/prod/users/dashboard", http.StatusSeeOther) // 303
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func logoutProdHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the template file (adjust path if necessary)
	tmpl, err := template.ParseFiles("template/logout.html")
	if err != nil {
		log.Println("Template parsing error:", err)
		http.Error(w, "Error loading logout page", http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Execute the template, sending data if needed (or nil if not)
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println("Template execution error:", err)
			http.Error(w, "Error rendering logout page", http.StatusInternalServerError)
		}
	case http.MethodPost:
		session, err := session.GetSessionFromContext(r)
		if session == nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Clear session values
		session.Values = nil

		//delete session before logout
		session.Options.MaxAge = -1

		// Save the session to apply the deletion
		err = session.Save(r, w)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if the session was deleted
		deletedSession, err := session.Store().Get(r, "session-name")

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(deletedSession.Values) == 0 {
			log.Println("Session successfully deleted.")
		} else {
			log.Fatalln("Failed to delete session.")
		}

		// If logout is successful, redirect
		// Redirect to home page on successful logout
		log.Println("redirect-to-homepage")
		http.Redirect(w, r, "/prod/auth/logout", http.StatusSeeOther) // 303
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// add a new user to the list
		var newUser user.User
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &newUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if newUser.UserID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//register user
		res, err := registerUserService(newUser)

		if err == nil {
			w.WriteHeader(res)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// add a new product to the list
		var existingUser user.User
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &existingUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if existingUser.UserID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//login user
		res, err := loginUserService(existingUser)

		if err != nil {
			w.WriteHeader(res)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}