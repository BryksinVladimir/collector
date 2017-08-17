package client

import (
	"strconv"

	"mobilda/client/request"
	"mobilda/client/response"
	"mobilda/errors"

	"github.com/dghubble/sling"
)

const DEFAULT_API_FORMAT = "json"

type MobildaOfferResponse struct {
	Summary response.Pagination `json:"summary"`
	Offers  []struct {
		Attributes struct {
			BusinessModel      string        `json:"business_model"`
			Currency           string        `json:"currency"`
			Description        string        `json:"description"`
			Domain             string        `json:"domain"`
			ID                 string        `json:"id"`
			OfferType          interface{}   `json:"offer_type"`
			PackageName        string        `json:"package_name"`
			ParametersRequired []interface{} `json:"parameters_required"`
			PreviewURL         string        `json:"preview_url"`
			Rate               interface{}   `json:"rate"`
			Status             string        `json:"status"`
			Thumbnail          string        `json:"thumbnail"`
			Title              string        `json:"title"`
			TrackingURL        string        `json:"tracking_url"`
		} `json:"attributes"`
		Capping struct {
			CapAmount        string `json:"cap_amount"`
			CapCurrentAmount string `json:"cap_current_amount"`
			CapEnable        string `json:"cap_enable"`
			CapFrequency     string `json:"cap_frequency"`
			CappingField     string `json:"capping_field"`
			CappingTimeframe string `json:"capping_timeframe"`
		} `json:"capping"`
		MobileAttributes struct {
			MinOsVersion     []string `json:"MinOs_version"`
			AllowedDevices   []string `json:"allowed_devices"`
			AppPrice         string   `json:"app_price"`
			AppRating        string   `json:"app_rating"`
			ContentRating    string   `json:"content_rating"`
			Developer        string   `json:"developer"`
			DeveloperWebsite string   `json:"developer_website"`
			MobileSupport    string   `json:"mobile_support"`
			PromoVideo       string   `json:"promo_video"`
		} `json:"mobile_attributes"`
		Targeting struct {
			BlackListSources []string `json:"black_list_sources"`
			Categories       []string `json:"categories"`
			Cities           []string `json:"cities"`
			Countries        []string `json:"countries"`
			Languages        []string `json:"languages"`
		} `json:"targeting"`
	} `json:"products"`
}

func (client *MobildaClient) Offers(accountId int, limit, page uint32) (*MobildaOfferResponse, *response.Error, error) {
	pageLimitParams := request.PageLimit{Limit: limit, Page: page}
	accIndex := client.getAccountIndexById(accountId)

	apiParams := request.ApiParams{
		ApiHash:   client.accounts[accIndex].Hash,
		ApiFeedId: strconv.Itoa(client.accounts[accIndex].FeedId),
		Format:    DEFAULT_API_FORMAT,
	}

	if !pageLimitParams.IsValid() {
		return nil, nil, errors.ErrLimitPage
	}

	if !apiParams.IsValid() {
		return nil, nil, errors.ErrApiParams
	}

	url := client.makeUrl(accIndex, "/xml/cpa_feeds/feed.php")
	req, err := sling.New().Get(url).QueryStruct(pageLimitParams).QueryStruct(apiParams).Request()
	if err != nil {
		client.log.Error(err)
		return nil, nil, err
	}

	buffer := &MobildaOfferResponse{}
	er, err := client.request(req, buffer)
	return buffer, er, err
}

func (client *MobildaClient) OffersTotal() (uint32, *response.Error, error) {
	resp, er, err := client.Offers(1, 1, 1)
	return resp.Summary.TotalRows, er, err
}
