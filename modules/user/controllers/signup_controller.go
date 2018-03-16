package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-pg/pg"

	"devin/database"
	"devin/helpers"
	"devin/models"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

// Signup handle user registration in application.
// Method: POST
// Content-Type: josn/application
//
// برای ثبت نام در سامانه اطلاعات ایمیل- نام کاربری-رمز عبور الزامی است
// تاییدیه ایمیل ارسال می شود
func Signup(w http.ResponseWriter, r *http.Request) {

	// Check content type
	if !helpers.HasJSONRequest(r) {
		err := helpers.ErrorResponse{
			Message:   "Invalid content type.",
			ErrorCode: http.StatusUnsupportedMediaType,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	var user models.User
	e := json.NewDecoder(r.Body).Decode(&user)
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "Bad request body",
			ErrorCode: http.StatusBadRequest,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	// lowercase username and email
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)

	e, errorMessages := validateSignupInputs(user)
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   e.Error(),
			ErrorCode: http.StatusUnprocessableEntity,
			Errors:    errorMessages,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	db := database.NewPGInstance()
	defer db.Close()
	// Check for duplication of email
	is, _ := isUniqueValue(db, "email", user.Email, 0)
	if is == false {
		messages := make(map[string][]string)
		messages["email"] = []string{"This email is already registered."}
		err := helpers.ErrorResponse{
			Message:   "Invalid Email address.",
			ErrorCode: http.StatusUnprocessableEntity,
			Errors:    messages,
		}
		helpers.NewErrorResponse(w, &err)
		return

	}
	// Check for duplication of username
	is, _ = isUniqueValue(db, "username", user.Username, 0)
	if is == false {
		messages := make(map[string][]string)
		messages["username"] = []string{"This username is already registered."}
		err := helpers.ErrorResponse{
			Message:   "Invalid username.",
			ErrorCode: http.StatusUnprocessableEntity,
			Errors:    messages,
		}
		helpers.NewErrorResponse(w, &err)
		return

	}

	user.UserType = 1
	user.SetEncryptedPassword(user.Password)
	user.SetNewEmailVerificationToken()

	// Saving data
	_, e = db.Model(&user).Insert(&user)

	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "Internal server error.",
			ErrorCode: http.StatusInternalServerError,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	json.NewEncoder(w).Encode(&user)
}

// validateSignupInputs check data validatin on user registration
func validateSignupInputs(user models.User) (e error, errMessages map[string][]string) {
	// Validate inputs
	hasError := false
	errMessages = make(map[string][]string)

	isValidEmail := helpers.Validator{}.IsValidEmailFormat(user.Email)
	if isValidEmail == false {
		hasError = true
		errMessages["email"] = []string{"Invalid Email address"}
	}

	isValidUsername := helpers.Validator{}.IsValidUsernameFormat(user.Username)
	if isValidUsername == false {
		hasError = true
		errMessages["username"] = []string{"Invalid username"}
	}

	if len(user.Password) < 6 {
		hasError = true
		errMessages["password"] = []string{"Password length must be greater than 6 characters"}
	}

	if hasError {
		return errors.New("Validation failed"), errMessages
	}

	return nil, nil
}

// isUniqueValue check duplication of value in given column of users table.
// ignoredID use for ignore given ID of checking. Set ignoredID to 0 if you want to check all records.
func isUniqueValue(db *pg.DB, columnName string, value string, ignoredID uint64) (isUnique bool, e error) {
	qry := db.Model(&models.User{}).Where(columnName+"=?", value)
	if ignoredID != 0 {
		qry = qry.Where("id != ?", ignoredID)
	}

	cnt, e := qry.Count()
	if e != nil {
		return false, e
	}

	if cnt != 0 {
		return false, nil
	}
	return true, nil
}