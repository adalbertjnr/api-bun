package store

import (
	"context"
	"fmt"

	"api/types"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
)

type Storager interface {
	CreateAccount(account *types.Account) error
	DeleteAccount(id int) error
	UpdateAccount(account *types.Account) error
	GetAccountById(id int) (*types.Account, error)
	GetAccounts() ([]*types.Account, error)
	GetAccountByNumber(id int) (*types.Account, error)
	Init() error
}

type PostgresStore struct {
	db  *bun.DB
	ctx context.Context
}

func NewPostgresStore(db *bun.DB) (*PostgresStore, error) {

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
	_, err := s.db.NewCreateTable().Model((*types.Account)(nil)).IfNotExists().Exec(s.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) CreateAccount(account *types.Account) error {
	_, err := s.db.NewInsert().
		Model(account).
		Exec(s.ctx)
	if err != nil {
		return err
	}
	return nil

}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.NewDelete().
		Model(&types.Account{}).
		Where("id = ?", id).
		Exec(s.ctx)
	if err != nil {
		return fmt.Errorf("error deleting user form the database")
	}
	return nil
}

func (s *PostgresStore) GetAccountByNumber(id int) (*types.Account, error) {
	var account types.Account
	err := s.db.NewSelect().
		Model(&account).
		Where("number = ?", id).
		Scan(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("account not found by number")
	}
	return &account, nil
}

func (s *PostgresStore) UpdateAccount(account *types.Account) error {
	return nil
}

func (s *PostgresStore) GetAccounts() ([]*types.Account, error) {
	var accounts []*types.Account
	err := s.db.NewSelect().Model(&accounts).Scan(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("not able to query the accounts")
	}
	return accounts, nil
}

func (s *PostgresStore) GetAccountById(id int) (*types.Account, error) {
	filteredAccountById := new(types.Account)
	if err := s.db.NewSelect().
		Model(filteredAccountById).
		Where("id = ?", id).
		Scan(s.ctx); err != nil {
		return nil, fmt.Errorf("account %d not found", id)
	}
	return filteredAccountById, nil
}
