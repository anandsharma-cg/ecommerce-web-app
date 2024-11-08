package authentication

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ecommerce/user"
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
}

func registerProdHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the template file (adjust path if necessary)
	tmpl, err := template.ParseFiles("template/register.html")
	if err != nil {
		http.Error(w, "Error loading register page", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		// Execute the template, sending data if needed (or nil if not)
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error rendering register page", http.StatusInternalServerError)
			log.Println("Template execution error:", err)
		}
	case http.MethodPost:

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		// extracting data of form values
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		// email and pass empty validation
		if email == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.Execute(w, map[string]string{"Error": errors.New("email and password are required").Error()})
			return
		}

		// pass and confirm pass validation
		if password != confirmPassword {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.Execute(w, map[string]string{"Error": errors.New("password and confirm password is not same").Error()})
			return
		}

		//register user
		newUser := user.User{Email: email, Password: password}
		res, err := registerUserService(newUser)

		if err != nil {
			w.WriteHeader(res)
			tmpl.Execute(w, map[string]string{"Error": err.Error()})
			return
		}

		// If register is successful, redirect
		// Redirect to dashboard page on successful login
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
		http.Error(w, "Error loading login page", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Execute the template, sending data if needed (or nil if not)
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error rendering login page", http.StatusInternalServerError)
			log.Println("Template execution error:", err)
		}
	case http.MethodPost:
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		// extracting data of form values
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Simple validation
		if email == "" || password == "" {
			http.Error(w, "email and password are required", http.StatusBadRequest)
			return
		}

		//login user
		existingUser := user.User{Email: email, Password: password}
		res, err := loginUserService(existingUser)

		if err != nil {
			w.WriteHeader(res)
			tmpl.Execute(w, map[string]string{"Error": err.Error()})
			return
		}

		// If login is successful, redirect
		// Redirect to dashboard page on successful login
		http.Redirect(w, r, "/prod/users/dashboard", http.StatusSeeOther) // 303
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
