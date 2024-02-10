package types

import (
	"math/rand"
	"time"

	"github.com/uptrace/bun"
)

type TransferRequest struct {
	DestAccount int `json:"dest_account"`
	Amount      int `json:"amount"`
}

type CreateAccoutRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Account struct {
	bun.BaseModel `bun:"table:accounts,alias:a"`
	ID            int       `bun:",pk,autoincrement" json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Number        int64     `bun:",unique" json:"number"`
	Balance       int64     `json:"balance"`
	CreatedAt     time.Time `json:"createdAt"`
}

func NewAccount(fn, ln string) *Account {
	return &Account{
		FirstName: fn,
		LastName:  ln,
		Number:    int64(rand.Intn(10000)),
		CreatedAt: time.Now().UTC(),
	}
}
