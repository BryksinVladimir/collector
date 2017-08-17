package client

import (
	"reflect"
	"strconv"
	"time"

	"mobilda/model"

	"bitbucket.org/mobio/go-logger"
	"github.com/sirupsen/logrus"
)

const (
	OffersMaxLimit = 500
)

type MobildaApiReader struct {
	client *MobildaClient

	log *logger.Logger
}

func NewMobildaApiReader(c *MobildaClient, l *logger.Logger) *MobildaApiReader {
	return &MobildaApiReader{
		client: c,
		log:    l,
	}
}

func (mar *MobildaApiReader) Offers(accountId int, page, limit uint32, stop <-chan bool) <-chan model.Offer {
	if limit > OffersMaxLimit {
		limit = OffersMaxLimit
	}
	results := make(chan model.Offer)
	go func() {
		defer close(results)

		l, p := limit, page
		retries := 5
		for {

			select {
			case <-stop:
				return
			default:
			}
			offers, er, err := mar.client.Offers(accountId, l, p)
			if er != nil || err != nil {
				retries--
				time.Sleep(time.Millisecond * 300)
				if retries == 0 {
					return
				}
				continue
			}
			retries = 5
		Loop:
			for _, offer := range offers.Offers {
				select {
				case <-stop:
					return
				default:
					offer_id, err := strconv.ParseUint(offer.Attributes.ID, 10, 64)
					if err != nil {
						mar.log.WithFields(logrus.Fields{
							"collector": "mobilda-offers-collector",
						}).Warnf("Mobilda Offer [ID: %s] has invalid ID", offer.Attributes.ID)
						continue Loop
					}
					results <- model.Offer{
						Id:            offer_id,
						PackageName:   offer.Attributes.PackageName,
						Title:         offer.Attributes.Title,
						Description:   offer.Attributes.Description,
						Domain:        offer.Attributes.Domain,
						PreviewUrl:    offer.Attributes.PreviewURL,
						TrackingUrl:   offer.Attributes.TrackingURL,
						BusinessModel: offer.Attributes.BusinessModel,
						Rate: func() string {
							if reflect.ValueOf(offer.Attributes.Rate).Kind() != reflect.String {
								return strconv.FormatFloat(offer.Attributes.Rate.(float64), 'E', -1, 64)
							}
							return offer.Attributes.Rate.(string)
						}(),
						Currency:         offer.Attributes.Currency,
						Thumbnail:        offer.Attributes.Thumbnail,
						Countries:        offer.Targeting.Countries,
						Cities:           offer.Targeting.Cities,
						Categories:       offer.Targeting.Categories,
						Languages:        offer.Targeting.Languages,
						BlackListSources: offer.Targeting.BlackListSources,
						MobileSupport:    offer.MobileAttributes.MobileSupport,
						AllowedDevices:   offer.MobileAttributes.AllowedDevices,
						MinOsVersion:     offer.MobileAttributes.MinOsVersion,
						AppPrice:         offer.MobileAttributes.AppPrice,
						AppRating:        offer.MobileAttributes.AppRating,
						ContentRating:    offer.MobileAttributes.ContentRating,
						Developer:        offer.MobileAttributes.Developer,
						DeveloperWebsite: offer.MobileAttributes.DeveloperWebsite,
						PromoVideo:       offer.MobileAttributes.PromoVideo,
						CapEnable:        offer.Capping.CapEnable,
						CapAmount:        offer.Capping.CapAmount,
						CapCurrentAmount: offer.Capping.CapCurrentAmount,
						CapFrequency:     offer.Capping.CapFrequency,
						CappingField:     offer.Capping.CappingField,
						CappingTimeframe: offer.Capping.CappingTimeframe,
						IsActive:         model.OfferStatusActive,
						StatusChangedAt:  time.Now(),
					}
				}
			}

			if offers.Summary.CurrentPage < offers.Summary.TotalPages {
				p = offers.Summary.CurrentPage + 1
			} else {
				return
			}
		}

	}()

	return results
}
