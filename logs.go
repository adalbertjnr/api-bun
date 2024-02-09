package main

import "github.com/sirupsen/logrus"

type LogMiddleware struct {
	next Storager
}

func NewLogMiddleware(next Storager) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) CreateAccount(account *Account) (err error) {
	defer func() {
		logrus.WithFields(logrus.Fields{
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
		logrus.WithFields(logrus.Fields{
			"ID":  id,
			"err": err,
		}).Infof("deleted id %d", id)
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
	defer func() {
		logrus.WithFields(logrus.Fields{
			"ID":        account.ID,
			"FirstName": account.FirstName,
			"LastName":  account.LastName,
			"Number":    account.Number,
			"CreatedAt": account.CreatedAt,
			"err":       err,
		}).Infof("searching account by id %d", id)
	}()
	return l.next.GetAccountById(id)
}

func (l *LogMiddleware) Init() (err error) {
	logrus.WithFields(logrus.Fields{
		"err": err,
	}).Info("initializing db setup")
	return l.next.Init()
}
