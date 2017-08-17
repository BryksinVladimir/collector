package client

import (
	"testing"
	"time"

	"mobilda/model"

	"bitbucket.org/mobio/go-config"
	"bitbucket.org/mobio/go-logger"
	"github.com/stretchr/testify/suite"
)

type MobildaClientSuite struct {
	suite.Suite
	client *MobildaClient
	logger *logger.Logger
	config *config.Config
}

func (suite *MobildaClientSuite) SetupSuite() {
	suite.logger = logger.NewLogger()

	accounts := []*model.Account{}

	test_account := model.Account{
		Id:     1,
		Name:   "standard",
		Hash:   "2b24eb1a2286820356acf4cd5c507907",
		FeedId: 351,
		Url:    "http://s.marsfeeds.com",
	}

	accounts = append(accounts, &test_account)

	suite.client = NewMobildaClient(accounts, nil, time.Second*60, 5, suite.logger)
}

func TestSuite(t *testing.T) {
	suite.Run(t, &MobildaClientSuite{})
}
