package main

import (
	"api/api"
	midd "api/jwt"
	"api/logs"
	"api/store"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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
	router := mux.NewRouter()
	router.HandleFunc("/account", s.MakeHTTPHandler(s.HandleAccount))
	router.HandleFunc("/account/{id}", midd.JWTAuthentication(s.MakeHTTPHandler(s.HandleById), s.Store))
	router.HandleFunc("/transfer", s.MakeHTTPHandler(s.HandleTransfer))

	fmt.Println("the server is running on port", s.ListenAddress)
	if err := http.ListenAndServe(s.ListenAddress, router); err != nil {
		log.Fatal(err)
	}
}
