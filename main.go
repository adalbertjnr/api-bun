package main

import (
	"api/api"
	"api/internal/util"
	"api/jwt"
	"api/logs"
	"api/store"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandler(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			util.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	}
}

func main() {
	var (
		storeSvc store.Storager
		err      error
	)
	storeSvc, err = store.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	if err := storeSvc.Init(); err != nil {
		log.Fatal(err)
	}
	storeSvc = logs.NewLogMiddleware(storeSvc)

	s := api.NewAPIServer(":3000", storeSvc)
	fmt.Println("running on port:", s.ListenAddress)
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandler(s.HandleAccount))
	router.HandleFunc("/account/{id}", jwt.JWTAuthentication(makeHTTPHandler(s.HandleById)))
	router.HandleFunc("/transfer", makeHTTPHandler(s.HandleTransfer))

	if err := http.ListenAndServe(s.ListenAddress, router); err != nil {
		log.Fatal(err)
	}
}
