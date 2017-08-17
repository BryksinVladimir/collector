package offers

import (
	"context"
	"encoding/hex"
	"sync"
	"time"

	"mobilda/client"
	"mobilda/consts"
	"mobilda/model"

	"bitbucket.org/mobio/go-cache"
	"bitbucket.org/mobio/go-collector"
	"bitbucket.org/mobio/go-config"
	"bitbucket.org/mobio/go-dbmanager"
	"bitbucket.org/mobio/go-logger"
	"github.com/cnf/structhash"
)

var (
	_    collector.ICollector = (*OffersCollector)(nil)
	lock sync.RWMutex
)

type OffersCollector struct {
	*collector.BaseCollector
	ctx context.Context

	log    *logger.Logger
	config *config.Config
	client *client.MobildaClient
	db     *dbmanager.DbManager
	cache  *cache.Cache
	acs    []*model.Account

	init     sync.Once
	interval uint64

	statsLock sync.RWMutex
	isRunning bool
}

func NewOffersCollector(ctx context.Context) *OffersCollector {
	return &OffersCollector{
		ctx:           ctx,
		BaseCollector: collector.NewBaseCollector(),
		log:           logger.FromContext(ctx, consts.Logger_Component_Key),
		config:        config.FromContext(ctx, consts.Config_Component_Key),
		client:        client.FromContext(ctx, consts.MobildaClient_Component_Key),
		db:            dbmanager.FromContext(ctx, consts.DbManager_Component_Key),
		cache:         cache.FromContext(ctx, consts.Cache_Component_Key),
		acs:           ctx.Value(consts.Accounts_Key).([]*model.Account),
	}
}

func (this *OffersCollector) collectorInit() {
	defer func() func() {
		start := time.Now()
		return func() {
			this.log.Infof("Init Mobilda Offers collector... Time elapsed: %s", time.Since(start))
		}
	}()()

	// init hash cache
	this.initCache()

	this.BaseCollector.UpdateStats(this.config.GetString("collector.offers_interval"), 0, time.Time{})
}

func (this *OffersCollector) Run() {
	lock.RLock()
	if this.isRunning {
		this.log.Warnf("Mobilda Offers collector already running. Exit...")
		lock.RUnlock()
		return
	}
	lock.RUnlock()

	lock.Lock()
	this.isRunning = true
	lock.Unlock()
	defer func() {
		lock.Lock()
		this.isRunning = false
		lock.Unlock()
	}()
	defer this.UpdateStats()()

	this.init.Do(this.collectorInit)

	var wg sync.WaitGroup

	for _, acc := range this.acs {
		wg.Add(1)
		go this.collect(acc, &wg)
	}
	wg.Wait()
}

func (this *OffersCollector) collect(acc *model.Account, wg *sync.WaitGroup) error {
	reader := client.NewMobildaApiReader(this.client, this.log)

	stop := make(chan bool)
	defer close(stop)

	forInsert := []model.Offer{}
	loaded := []interface{}{}
	for item := range reader.Offers(acc.Id, 1, client.OffersMaxLimit, stop) {
		item.AccountId = acc.Id
		// check hash cache
		hash := hex.EncodeToString(structhash.Sha1(item, 1))

		loaded = append(loaded, item.Id)

		if _, ok := this.cache.Get(item.CacheId()); !ok {
			item.Hash = hash
			forInsert = append(forInsert, item)
			if len(forInsert)%client.OffersMaxLimit == 0 {
				this.bulkInsert(forInsert)
				forInsert = []model.Offer{}
			}
		} else {
			if h, _ := this.cache.Get(item.CacheId()); h != hash {
				item.Hash = hash
				err := this.db.Update(&item)
				if err != nil {
					this.log.Error(err)
				} else {
					this.cache.Set(item.CacheId(), hash, cache.NoExpiration)
				}
			}
		}
	}

	if len(forInsert) > 0 {
		this.bulkInsert(forInsert)
	}

	this.setStoppedStatus(loaded, acc.Id)

	wg.Done()

	return nil
}

func (this *OffersCollector) setStoppedStatus(loaded []interface{}, accountId int) {
	suspended := []model.Offer{}

	query := this.db.
		Model(&suspended).
		Where("is_active = ?", model.OfferStatusActive).
		Where("account_id = ?", accountId)

	if len(loaded) > 0 {
		query.WhereIn("offer_id NOT IN (?)", loaded...)
	}

	if err := query.Select(); err != nil {
		this.log.WithField("collector", "mobilda-offers-collector").Error(err)
	}

	now := time.Now()
	for _, offer := range suspended {
		offer.IsActive = model.OfferStatusStopped
		offer.StatusChangedAt = now
		oldHash := offer.Hash
		offer.Hash = hex.EncodeToString(structhash.Sha1(offer, 1))
		err := this.db.Update(&offer)
		if err != nil {
			this.log.Error(err)
		} else {
			this.cache.Set(offer.CacheId(), offer.Hash, cache.NoExpiration)
			this.cache.Delete(oldHash)
		}
	}
}

func (this *OffersCollector) bulkInsert(data []model.Offer) {
	if _, err := this.db.Model(&data).Insert(); err != nil {
		this.log.WithField("collector", "mobilda-offers-collector").Error(err)
	} else {
		for _, item := range data {
			this.cache.Set(item.CacheId(), item.Hash, cache.NoExpiration)
		}
	}
}

func (this *OffersCollector) initCache() {
	offers := []model.Offer{}
	if err := this.db.Model(&offers).
		Column("offer_id", "account_id", "hash").
		Select(); err != nil {
		this.log.Error(err)
		return
	}

	for _, item := range offers {
		this.cache.Set(item.CacheId(), item.Hash, cache.NoExpiration)
	}
}

func (this *OffersCollector) UpdateStats() func() {
	start := time.Now()
	this.BaseCollector.UpdateStats(this.TimeInterval(), 1, start)

	return func() {
		this.BaseCollector.StatsSetLastDuration(time.Since(start))
		this.log.Infof("Load Mobilda Offers... Time elapsed: %s", this.BaseCollector.Stats().LastRunDuration)
	}
}

// time interval (seconds)
func (this *OffersCollector) TimeInterval() string {
	if this.NewTimeInterval != "" {
		return this.NewTimeInterval
	}
	return this.config.GetString("collector.offers_interval")
}

func (this *OffersCollector) IsRunImmediate() bool {
	return this.config.GetBool("collector.offers_run_immediate")
}
