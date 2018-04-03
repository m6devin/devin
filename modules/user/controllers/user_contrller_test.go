package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"devin/database"
	"devin/middlewares"
	"devin/models"
)

func getValidUser(id int, isRoot bool) (user models.User, claim models.Claim, tokenString string) {
	db := database.NewPGInstance()
	defer db.Close()
	db.Exec(`delete from users where id=?; insert into users (id, username, email, is_root_user) values (?, ?, ?, ?)`, id, id, fmt.Sprintf("mgh%v", id), fmt.Sprintf("m6devin%v@gmail.com", id), isRoot)
	db.Model(&user).Where("id=?", id).First()
	claim = user.GenerateNewTokenClaim()
	tokenString, _ = user.GenerateNewTokenString(claim)

	return user, claim, tokenString
}

func deleteTestUser(id int) {
	db := database.NewPGInstance()
	defer db.Close()
	db.Exec(`delete from users where id=?;`, id)
}

func TestUpdateProfile(t *testing.T) {
	user1, _, tokenString := getValidUser(1, true)
	_, _, tokenStringStandardUser := getValidUser(2, false)
	defer deleteTestUser(1)
	defer deleteTestUser(2)

	route := mux.NewRouter()
	path := "/user/{id}/update"
	route.Handle(path, http.HandlerFunc(UpdateProfile))
	route.Use(middlewares.Authenticate)

	t.Run("bad content-type", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), nil)
		req.Header.Add("Authorization", tokenString)
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()
		if res.StatusCode != http.StatusUnsupportedMediaType {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid content type") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("no id variable passed to the mux router", func(t *testing.T) {
		route := mux.NewRouter()
		path := "/user/update"
		route.Handle(path, http.HandlerFunc(UpdateProfile))
		req, _ := http.NewRequest(http.MethodPost, path, nil)

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid User ID") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("invalid user_id data type", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "-1", 1), nil)

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid User ID. Just integer values accepted") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("not exists user in DB", func(t *testing.T) {
		deleteTestUser(0)
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "0", 1), nil)

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusInternalServerError {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Error on loading user data") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Authorization token error", func(t *testing.T) {
		route := mux.NewRouter()
		route.Handle(path, http.HandlerFunc(UpdateProfile))
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), nil)

		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusUnauthorized {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Auhtentication failed") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Access denied", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), nil)

		req.Header.Add("Authorization", tokenStringStandardUser)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusForbidden {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "This action is not allowed for you") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Empty Request body", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), nil)

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusInternalServerError {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Request body cant be empty") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader("Bad Request body"))

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusInternalServerError {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid request body") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("OK", func(t *testing.T) {
		user1.FirstName = "Updated first name"
		bts, _ := json.Marshal(&user1)
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), bytes.NewReader(bts))

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusOK {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ = ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Updated first name") {
			t.Fatal("Invalid response")
		}
	})

}
