package authentication

import (
	"log"
	"net/http"

	"github.com/ecommerce/user"
)

func registerUserService(newUser user.User) (int, error) {
	_, err := user.RegisterUser(newUser)
	if err != nil {
		log.Print(err)
		return http.StatusBadRequest, nil
	}
	return http.StatusOK, err
}

func loginUserService(existingUser user.User) (int, error) {
	_, err := user.LoginUser(existingUser)
	if err != nil {
		log.Print(err)
		return http.StatusBadRequest, nil
	}
	return http.StatusOK, err
}