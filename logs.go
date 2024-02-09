package main

import (
	"time"

	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Storager
	log  *logrus.Logger
}

func NewLogMiddleware(next Storager) *LogMiddleware {
	log := logrus.New()
	log.Formatter = new(logrus.JSONFormatter)
	return &LogMiddleware{
		next: next,
		log:  log,
	}
}

func (l *LogMiddleware) CreateAccount(account *Account) (err error) {
	defer func() {
		l.log.WithFields(logrus.Fields{
			"ID":        account.ID,
			"FirstName": account.FirstName,
			"LastName":  account.LastName,
			"Number":    account.Number,
			"CreatedAt": account.CreatedAt,
			"err":       err,
		}).Info("newly created account")
	}()
	return l.next.CreateAccount(account)
}

func (l *LogMiddleware) DeleteAccount(id int) (err error) {
	defer func() {
		l.log.WithFields(logrus.Fields{
			"ID":  id,
			"err": err,
		}).Infof("deleted account with id %d", id)
	}()
	return l.next.DeleteAccount(id)
}

func (l *LogMiddleware) UpdateAccount(account *Account) (err error) {
	defer func() {
		logrus.WithFields(logrus.Fields{
			"ID":        account.ID,
			"FirstName": account.FirstName,
			"LastName":  account.LastName,
			"Number":    account.Number,
			"CreatedAt": account.CreatedAt,
			"err":       err,
		}).Info("updated account")
	}()
	return l.next.CreateAccount(account)
}

func (l *LogMiddleware) GetAccountById(id int) (account *Account, err error) {
	var (
		fn        string
		ln        string
		number    int64
		createdAt time.Time
	)
	if account != nil {
		fn = account.FirstName
		ln = account.LastName
		number = account.Number
		createdAt = account.CreatedAt
	}
	defer func() {
		l.log.WithFields(logrus.Fields{
			"ID":        id,
			"FirstName": fn,
			"LastName":  ln,
			"Number":    number,
			"CreatedAt": createdAt,
			"err":       err,
		}).Infof("filtering account by id %d", id)
	}()
	return l.next.GetAccountById(id)
}

func (l *LogMiddleware) Init() (err error) {
	logrus.WithFields(logrus.Fields{
		"err": err,
	}).Info("initializing db setup")
	return l.next.Init()
}

func (l *LogMiddleware) GetAccounts() (accounts []*Account, err error) {
	defer func() {
		logrus.WithFields(logrus.Fields{
			"err": err,
		})
	}()
	return l.next.GetAccounts()
}
