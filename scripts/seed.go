package main

import (
	"api/store"
	"api/types"
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func main() {
	db, err := store.LoadDBWithBunClient()
	if err != nil {
		panic(err.Error())
	}

	s, err := store.NewPostgresStore(db)
	if err != nil {
		panic(err.Error())
	}
	dropTable(s, db)
	if err := s.CreateAccountTable(); err != nil {
		panic(err.Error())
	}
	seedAccounts(s, db)
}

func seedAccount(st store.Storager, db *bun.DB, fn, ln, pw string) *types.Account {
	acc, err := types.NewAccount(fn, ln, pw)
	if err != nil {
		panic(err.Error())
	}
	if err := st.CreateAccount(acc); err != nil {
		panic(err.Error())
	}
	fmt.Println("new account -> ", acc.FirstName, acc.LastName, acc.Number)
	return acc
}

func dropTable(s store.Storager, db *bun.DB) {
	_, err := db.NewDropTable().
		Model(&types.Account{}).
		IfExists().
		Exec(context.TODO())
	if err != nil {
		panic(err.Error())
	}
}

func seedAccounts(s store.Storager, db *bun.DB) {
	seedAccount(s, db, "foo", "ff", "pass123")
	seedAccount(s, db, "bar", "bb", "pass123")
}
