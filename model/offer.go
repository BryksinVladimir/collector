package model

import (
	"strconv"
	"time"
)

const (
	OfferStatusActive  = true
	OfferStatusStopped = false
)

type Offer struct {
	tableName        struct{} `sql:"mobilda.offer"`
	Id               uint64   `sql:"offer_id,pk"`
	AccountId        int      `sql:"account_id,pk"`
	PackageName      string   `sql:",notnull"`
	Title            string
	Description      string
	Domain           string `sql:",notnull"`
	PreviewUrl       string `sql:",notnull"`
	TrackingUrl      string
	BusinessModel    string
	Rate             string
	Currency         string
	Thumbnail        string
	Countries        []string `pg:",array"`
	Cities           []string `pg:",array"`
	Categories       []string `pg:",array"`
	Languages        []string `pg:",array"`
	BlackListSources []string `pg:",array"`
	MobileSupport    string
	AllowedDevices   []string `pg:",array"`
	MinOsVersion     []string `pg:",array"`
	AppPrice         string
	AppRating        string
	ContentRating    string
	Developer        string
	DeveloperWebsite string
	PromoVideo       string
	CapEnable        string
	CapAmount        string
	CapCurrentAmount string
	CapFrequency     string
	CappingField     string
	CappingTimeframe string
	IsActive         bool      `sql:",notnull"`
	StatusChangedAt  time.Time `hash:"-"`
	Hash             string    `hash:"-"`
}

func (this Offer) CacheId() string {
	return "moboffer:" + strconv.FormatUint(this.Id, 10) + "account" + strconv.Itoa(this.AccountId)
}
