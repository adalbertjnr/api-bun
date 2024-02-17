package main

import (
	"api/api"
	midd "api/jwt"
	"api/logs"
	"api/store"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var (
		storeSvc store.Storager
		err      error
	)
	db, err := store.LoadDBWithBunClient()
	if err != nil {
		panic(err.Error())
	}

	storeSvc, err = store.NewPostgresStore(db)
	if err != nil {
		log.Fatal(err)
	}
	if err := storeSvc.Init(); err != nil {
		log.Fatal(err)
	}
	storeSvc = logs.NewLogMiddleware(storeSvc)

	s := api.NewAPIServer(":3000", storeSvc)

	router := http.NewServeMux()
	router.HandleFunc("/login", s.MakeHTTPHandler(s.HandleLogin))
	router.HandleFunc("/account", s.MakeHTTPHandler(s.HandleAccount))
	router.HandleFunc("/account/{id}", midd.JWTAuthentication(s.MakeHTTPHandler(s.HandleById), s.Store))
	router.HandleFunc("/transfer", s.MakeHTTPHandler(s.HandleTransfer))

	fmt.Println("the server is running on port", s.ListenAddress)
	if err := http.ListenAndServe(s.ListenAddress, router); err != nil {
		log.Fatal(err)
	}
}
