package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type Storager interface {
	CreateAccount(account *Account) error
	DeleteAccount(id int) error
	UpdateAccount(account *Account) error
	GetAccountById(id int) (*Account, error)
	GetAccounts() ([]*Account, error)
	Init() error
}

type PostgresStore struct {
	db  *bun.DB
	ctx context.Context
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	sqldb, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := sqldb.Ping(); err != nil {
		return nil, err
	}
	db := bun.NewDB(sqldb, pgdialect.New())
	return &PostgresStore{
		db:  db,
		ctx: context.TODO(),
	}, nil
}

func (s *PostgresStore) Init() error {
	if err := s.CreateAccountTable(); err != nil {
		return fmt.Errorf("failed to create account table %w", err)
	}
	return nil
}

func (s *PostgresStore) CreateAccountTable() error {
	_, err := s.db.NewCreateTable().Model((*Account)(nil)).IfNotExists().Exec(s.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) CreateAccount(account *Account) error {
	_, err := s.db.NewInsert().
		Model(account).
		Exec(s.ctx)
	if err != nil {
		return err
	}
	return nil

}

func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}

func (s *PostgresStore) UpdateAccount(account *Account) error {
	return nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	var accounts []*Account
	err := s.db.NewSelect().Model(&accounts).Scan(s.ctx)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	filteredAccountById := new(Account)
	if err := s.db.NewSelect().
		Model(filteredAccountById).
		Where("id = ?", id).
		Scan(s.ctx); err != nil {
		return nil, err
	}
	return filteredAccountById, nil
}
