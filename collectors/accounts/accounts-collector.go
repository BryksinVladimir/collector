package accounts

import (
	"context"
	"time"

	"mobilda/consts"
	"mobilda/model"

	"bitbucket.org/mobio/go-dbmanager"
	"bitbucket.org/mobio/go-logger"
)

type AccountsCollector struct {
	ctx context.Context

	log *logger.Logger
	db  *dbmanager.DbManager
	acs []*model.Account

	startedAt time.Time
}

func NewAccountsCollector(ctx context.Context) *AccountsCollector {
	return &AccountsCollector{
		ctx: ctx,
		log: logger.FromContext(ctx, consts.Logger_Component_Key),
		db:  dbmanager.FromContext(ctx, consts.DbManager_Component_Key),
		acs: ctx.Value(consts.Accounts_Key).([]*model.Account),
	}
}

func (this *AccountsCollector) collectorInit() {
	defer func() func() {
		start := time.Now()
		return func() {
			this.log.Infof("Init Mobilda Accounts collector... Time elapsed: %s", time.Since(start))
		}
	}()()
}

func (this *AccountsCollector) Run() {
	this.collectorInit()

	this.startedAt = time.Now()

	forInsert := []model.Account{}
	for _, account := range this.acs {
		forInsert = append(forInsert, *account)
	}

	this.bulkInsert(forInsert)

	this.log.Infof("Mobilda Accounts updated... Time elapsed: %s", time.Since(this.startedAt))

}

func (this *AccountsCollector) bulkInsert(data []model.Account) {
	if _, err := this.db.Model(&data).OnConflict("(name) DO NOTHING").Insert(); err != nil {
		this.log.WithField("collector", "mobilda-account-collector").Error(err)
	}
}
