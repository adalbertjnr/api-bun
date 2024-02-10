package api

import (
	"api/internal/util"
	"api/store"
	"api/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	ListenAddress string
	store         store.Storager
}

func NewAPIServer(listenAddress string, store store.Storager) *APIServer {
	return &APIServer{
		ListenAddress: listenAddress,
		store:         store,
	}
}

// Handle account ep
func (s *APIServer) HandleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.HandleGetAccount(w, r)
	case http.MethodPost:
		return s.HandleCreateAccount(w, r)
	default:
		return fmt.Errorf("method not allowed")
	}
}

// account get && delete by id
func (s *APIServer) HandleById(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.HandleGetAccountById(w, r)
	case http.MethodDelete:
		return s.HandleDeleteAccountById(w, r)
	default:
		return fmt.Errorf("method not allowed")
	}
}

func (s *APIServer) HandleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	idVars := mux.Vars(r)["id"]
	idVarsInt, err := strconv.Atoi(idVars)
	if err != nil {
		return err
	}
	acc, err := s.store.GetAccountById(idVarsInt)
	if err != nil {
		return err
	}
	return util.WriteJSONResponse(w, http.StatusOK, acc)
}

func (s *APIServer) HandleDeleteAccountById(w http.ResponseWriter, r *http.Request) error {
	idVarsInt, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}
	if err := s.store.DeleteAccount(idVarsInt); err != nil {
		return err
	}
	return util.WriteJSONResponse(w, http.StatusAccepted, map[string]int{"deleted user id": idVarsInt})
}

func (s *APIServer) HandleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return util.WriteJSONResponse(w, http.StatusOK, accounts)
}

func (s *APIServer) HandleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	accountParams := new(types.CreateAccoutRequest)
	if err := json.NewDecoder(r.Body).Decode(accountParams); err != nil {
		return err
	}
	defer r.Body.Close()

	newAccount := types.NewAccount(accountParams.FirstName, accountParams.LastName)
	if err := s.store.CreateAccount(newAccount); err != nil {
		return err
	}
	return util.WriteJSONResponse(w, http.StatusOK, newAccount)
}

func (s *APIServer) HandleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(types.TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()
	return util.WriteJSONResponse(w, http.StatusOK, transferReq)
}
