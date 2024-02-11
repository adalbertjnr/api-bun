package api

import (
	"api/internal/util"
	midd "api/jwt"
	"api/store"
	"api/types"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	ListenAddress string
	Store         store.Storager
}

func NewAPIServer(listenAddress string, Store store.Storager) *APIServer {
	return &APIServer{
		ListenAddress: listenAddress,
		Store:         Store,
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

func (s *APIServer) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return util.WriteJSONResponse(w, http.StatusBadRequest, util.NewError(types.ErrMethodNotAllowed))
	}
	req := new(types.LoginParams)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	acc, err := s.Store.GetAccountByNumber(req.Number)
	if err != nil {
		return err
	}
	if !types.BcryptValidator(acc.EncryptedPassword, req.Password) {
		return util.WriteJSONResponse(w, http.StatusBadRequest, util.NewError(errors.New("not allowed")))
	}
	token, err := midd.NewJWTToken(acc)
	if err != nil {
		return err
	}
	resp := types.LoginResponse{
		Account: *acc,
		Token:   token,
	}
	return util.WriteJSONResponse(w, http.StatusAccepted, resp)
}
func (s *APIServer) HandleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	idVars := mux.Vars(r)["id"]
	idVarsInt, err := strconv.Atoi(idVars)
	if err != nil {
		return err
	}
	acc, err := s.Store.GetAccountById(idVarsInt)
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
	if err := s.Store.DeleteAccount(idVarsInt); err != nil {
		return err
	}
	return util.WriteJSONResponse(w, http.StatusAccepted, map[string]int{"deleted user id": idVarsInt})
}

func (s *APIServer) HandleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.Store.GetAccounts()
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

	newAccount, err := types.NewAccount(accountParams.FirstName, accountParams.LastName, accountParams.Password)
	if err != nil {
		return err
	}
	if err := s.Store.CreateAccount(newAccount); err != nil {
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

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func (s *APIServer) MakeHTTPHandler(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			util.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	}
}
