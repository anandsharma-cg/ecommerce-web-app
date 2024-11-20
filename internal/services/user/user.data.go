package user

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/ecommerce/database"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	// Convert the string password to a byte slice
	passwordBytes := []byte(password)

	// Generate the hashed password
	hashed, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Convert the hashed password back to a string
	return string(hashed), nil
}

func getUserByEmail(email string) (*User, error) {
	row := database.DbConn.QueryRow(`SELECT 
	userId, 	
	email,
	password,
	isAdmin	
	FROM users
	WHERE email = ?`, email)

	user := &User{}
	err := row.Scan(
		&user.UserID,
		&user.Email,
		&user.Password,
		&user.IsAdmin)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return user, nil
}

func getUser(userID int) (*User, error) {
	row := database.DbConn.QueryRow(`SELECT 
	userId, 	
	email,
	password,
	isAdmin	
	FROM users
	WHERE userId = ?`, userID)

	user := &User{}
	err := row.Scan(
		&user.UserID,
		&user.Email,
		&user.Password,
		&user.IsAdmin)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return user, nil
}

func removeUser(userID int) error {
	_, err := database.DbConn.Exec(`DELETE FROM users where userId = ?`, userID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func getAllUsers() ([]User, error) {
	results, err := database.DbConn.Query(`SELECT 
	userId, 	 
	email,
	password,
	isAdmin	
	FROM users`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	users := make([]User, 0)
	for results.Next() {
		var user User
		results.Scan(&user.UserID,
			&user.Email,
			&user.Password,
			&user.IsAdmin)

		users = append(users, user)
	}
	return users, nil
}

func updateEmail(user User) error {
	_, err := database.DbConn.Exec(`UPDATE users SET 		 				 
		email=?		
		WHERE userId=?`,
		user.Email,
		user.UserID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func updatePassword(user User) error {
	hashedPass, err := hashPassword(user.Password)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	result, err := database.DbConn.Exec(`UPDATE users SET 		 				 
		password=?		
		WHERE userId=?`,
		hashedPass,
		user.UserID)
	if err != nil {
		result.RowsAffected()
		log.Println(err.Error())
		return err
	}
	return nil
}

func updateUser(user User) error {
	err_e := updateEmail(user)
	// err_p := updatePassword(user)

	if err_e != nil {
		log.Println(err_e.Error())
		return err_e
	}

	// if err_p != nil {
	// 	log.Println(err_p.Error())
	// 	return err_p
	// }

	return nil
}

func RegisterUser(user User) (int, error) {
	return addUser(user)
}
func LoginUser(user User) (int, error) {
	existingUser, err := getUserByEmail(user.Email)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Check if user was found (i.e., existingUser is not nil)
	if existingUser == nil {
		return http.StatusNotFound, errors.New("user not found")
	}

	//compare existing-hashed-pass and request-pass
	isCredMisMatchError := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if isCredMisMatchError != nil {
		return http.StatusUnauthorized, errors.New("incorrect password")
	}
	return http.StatusOK, nil
}
func addUser(user User) (int, error) {
	hashedPass, err := hashPassword(user.Password)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	result, err := database.DbConn.Exec(`INSERT INTO users  
	(email,
	password) VALUES (?, ?)`,
		user.Email,
		hashedPass)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}