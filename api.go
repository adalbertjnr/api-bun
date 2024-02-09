package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddress string
	store         Storager
}

func NewAPIServer(listenAddress string, store Storager) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		store:         store,
	}
}

// router
func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandler(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandler(s.handleGetAccountById))

	fmt.Println("running on port:", s.listenAddress)
	if err := http.ListenAndServe(s.listenAddress, router); err != nil {
		log.Fatal(err)
	}
}

// handle account ep
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGetAccount(w, r)
	case http.MethodPost:
		return s.handleCreateAccount(w, r)
	case http.MethodDelete:
		return s.handleDeleteAccount(w, r)
	case http.MethodPut:
		return s.handleTransfer(w, r)
	default:
		return fmt.Errorf("method not allowed")
	}
}

// account eps
func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	idVars := mux.Vars(r)["id"]
	idVarsInt, err := strconv.Atoi(idVars)
	if err != nil {
		return err
	}
	acc, err := s.store.GetAccountById(idVarsInt)
	if err != nil {
		return err
	}
	return WriteJSONResponse(w, http.StatusOK, acc)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSONResponse(w, http.StatusOK, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	accountParams := new(CreateAccoutRequest)
	if err := json.NewDecoder(r.Body).Decode(accountParams); err != nil {
		return err
	}

	newAccount := NewAccount(accountParams.FirstName, accountParams.LastName)
	if err := s.store.CreateAccount(newAccount); err != nil {
		return err
	}
	return WriteJSONResponse(w, http.StatusOK, newAccount)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// helpers
func WriteJSONResponse(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandler(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			WriteJSONResponse(w, http.StatusBadGateway, ApiError{Error: err.Error()})
		}
	}
}
