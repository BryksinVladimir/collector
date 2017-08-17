package model

type Account struct {
	tableName struct{} `sql:"mobilda.account"`
	Name      string   `mapstructure:"account_name"`
	Id        int      `mapstructure:"account_id"`
	Hash      string   `sql:"-"`
	FeedId    int      `mapstructure:"feed_id" sql:"-"`
	Url       string   `sql:"-"`
}
