package server

import (
	h "mobilda/server/handlers"
)

var (
	ah h.ApiHandlers
)

func (srv *AppServer) InitRouter() {
	srv.Router.Post("/run/:collector", ah.RunCollector())
}
