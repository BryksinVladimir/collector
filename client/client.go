package client

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"mobilda/client/response"
	"mobilda/model"

	"bitbucket.org/mobio/go-logger"
	"github.com/beefsack/go-rate"
	"github.com/sirupsen/logrus"
)

const (
	DEFAULT_RATE_LIMIT = 5 //max 5 requests per second
)

type IClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type MobildaClient struct {
	log        *logger.Logger
	httpClient IClient
	timeout    time.Duration
	accounts   []*model.Account

	rateLimiter *rate.RateLimiter
}

//NewMobildaClient - creates and returns new mobilda api client
func NewMobildaClient(acc []*model.Account, client IClient, timeout time.Duration, ratelimit int, l *logger.Logger) *MobildaClient {
	var c IClient

	if client != nil {
		c = client
	} else {
		c = &http.Client{
			Timeout: timeout,
		}
	}

	// Set rate limit
	limit := 0
	if ratelimit == 0 {
		limit = DEFAULT_RATE_LIMIT
	} else {
		limit = ratelimit
	}
	rl := rate.New(limit, time.Second)

	return &MobildaClient{
		log:         l,
		httpClient:  c,
		accounts:    acc,
		rateLimiter: rl,
	}
}

func (client *MobildaClient) request(r *http.Request, responseData interface{}) (*response.Error, error) {
	client.rateLimiter.Wait()

	resp, err := client.httpClient.Do(r)
	if err != nil {
		client.log.WithField("url", r.URL).Error(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	result := &response.Error{}
	if err := json.Unmarshal(body, result); err != nil {
		client.log.WithFields(logrus.Fields{"url": r.URL, "response": string(body)}).Error(err)
		return nil, err
	} else if result.Error != 0 {
		client.log.WithFields(logrus.Fields{
			"url": r.URL, "error": result.Error, "error_message": result.ErrorMessage,
		}).Error(result.Error)
		return result, nil
	}

	if err := json.Unmarshal(body, responseData); err != nil {
		client.log.WithFields(logrus.Fields{"url": r.URL, "response": len(body)}).Error(err)
		return nil, err
	}

	return nil, nil
}

func (client *MobildaClient) getAccountIndexById(accountId int) (accIndex int) {
	for i, acc := range client.accounts {
		if acc.Id == accountId {
			accIndex = i
			return accIndex
		}
	}
	return accIndex
}

func FromContext(ctx context.Context, key string) *MobildaClient {
	return ctx.Value(key).(*MobildaClient)
}
