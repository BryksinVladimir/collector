package mobilda

import (
	"mobilda/collectors/offers"
)

func (app *Application) addCollectors() {
	app.scheduler.AddTimeIntervalCollector("offers-collector", offers.NewOffersCollector(app.ctx))
}