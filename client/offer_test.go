package client

import (
	"github.com/stretchr/testify/assert"
)

func (suite MobildaClientSuite) TestMobildaClient_Offers() {
	t := suite.T()
	t.Parallel()
	resp, er, err := suite.client.Offers(1, 100, 1)
	assert.Nil(t, er, "response.Error must be nil")
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, resp.Summary.Limit, uint32(100))
	assert.Equal(t, resp.Summary.CurrentPage, uint32(1))
	assert.Equal(t, resp.Summary.TotalRows > 0, true)
	assert.Equal(t, len(resp.Offers) > 0, true)
}

func (suite MobildaClientSuite) TestMobildaClient_OffersTotal() {
	t := suite.T()
	t.Parallel()
	total, er, err := suite.client.OffersTotal()
	assert.Nil(t, er, "response.Error must be nil")
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, total > 0, true)
}
