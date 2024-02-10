package logs

import (
	"api/store"
	"api/types"

	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next store.Storager
	log  *logrus.Logger
}

func NewLogMiddleware(next store.Storager) *LogMiddleware {
	log := logrus.New()
	log.Formatter = new(logrus.JSONFormatter)
	return &LogMiddleware{
		next: next,
		log:  log,
	}
}

func (l *LogMiddleware) CreateAccount(account *types.Account) (err error) {
	defer func() {
		l.log.WithFields(logrus.Fields{
			"ID":        account.ID,
			"FirstName": account.FirstName,
			"LastName":  account.LastName,
			"Number":    account.Number,
			"CreatedAt": account.CreatedAt,
			"err":       err,
		}).Info("new account created")
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

func (l *LogMiddleware) UpdateAccount(account *types.Account) (err error) {
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

func (l *LogMiddleware) GetAccountById(id int) (account *types.Account, err error) {

	defer func() {
		if account != nil {
			l.log.WithFields(logrus.Fields{
				"ID":        id,
				"FirstName": account.FirstName,
				"LastName":  account.LastName,
				"Number":    account.Number,
				"CreatedAt": account.CreatedAt,
				"err":       err,
			}).Infof("filtering account by id %d", id)
		} else {
			l.log.WithField("ID", id).Infof("id %d not found in database", id)
		}
	}()
	return l.next.GetAccountById(id)
}

func (l *LogMiddleware) Init() (err error) {
	logrus.WithFields(logrus.Fields{
		"err": err,
	}).Info("initializing db setup")
	return l.next.Init()
}

func (l *LogMiddleware) GetAccounts() (accounts []*types.Account, err error) {
	defer func() {
		logrus.WithFields(logrus.Fields{
			"err": err,
		})
	}()
	return l.next.GetAccounts()
}
