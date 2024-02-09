package main

import (
	"log"
)

func main() {
	var (
		storeSvc Storager
		err      error
	)
	storeSvc, err = NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	if err := storeSvc.Init(); err != nil {
		log.Fatal(err)
	}
	storeSvc = NewLogMiddleware(storeSvc)
	s := NewAPIServer(":3000", storeSvc)
	s.Run()
}
