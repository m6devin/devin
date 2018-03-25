package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"devin/middlewares"
	user_ctrl "devin/modules/user/controllers"
)

func LoadRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/signup", user_ctrl.Signup).Methods(http.MethodPost)
	r.HandleFunc("/signin", user_ctrl.Signin).Methods(http.MethodPost)

	secureArea := r.PathPrefix("/").Subrouter().StrictSlash(true)
	secureArea.Use(middlewares.Authenticate)
	secureArea.HandleFunc("/user/{id:[0-9]+}/update", user_ctrl.UpdateProfile).Methods(http.MethodPost)

	return r
}
