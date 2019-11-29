package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.Handler
}

type Routes []Route

var routes = Routes{
	Route{
		"SendEmail",
		"POST",
		"/send",
		SendHandler{},
	},
	Route{
		"ReceiveEmail",
		"POST",
		"/receive",
		ReceiveHandler{},
	},
	Route{
		"Healthcheck",
		"Get",
		"/healthcheck",
		HealthcheckHandler{},
	},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		log.Println(route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.Handler)
	}

	return router
}