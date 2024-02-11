package types

import (
	"errors"
	"math/rand"
	"time"

	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrAccessDenied     = errors.New("access denied")
	ErrMethodNotAllowed = errors.New("method not allowed")
)

type LoginParams struct {
	Number   int    `json:"number"`
	Password string `json:"password"`
}

type TransferRequest struct {
	DestAccount int `json:"dest_account"`
	Amount      int `json:"amount"`
}

type CreateAccoutRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type Account struct {
	bun.BaseModel     `bun:"table:accounts,alias:a"`
	ID                int       `bun:",pk,autoincrement" json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Number            int64     `bun:",unique" json:"number"`
	EncryptedPassword string    `json:"-"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"createdAt"`
}

func NewAccount(fn, ln, password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		FirstName:         fn,
		LastName:          ln,
		EncryptedPassword: string(encpw),
		Number:            int64(rand.Intn(10000)),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
