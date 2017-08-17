package client

import (
	"github.com/stretchr/testify/assert"
)

func (suite MobildaClientSuite) TestApiReader_Offers() {
	t := suite.T()
	t.Parallel()
	reader := NewMobildaApiReader(suite.client, suite.logger)
	stop := make(chan bool)
	defer close(stop)
	res := reader.Offers(1, 1, 100, stop)
	count := 0
	for _ = range res {
		count++
		if count == 100 {
			stop <- true
		}
	}
	assert.Equal(suite.T(), count, 100)
}
