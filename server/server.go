package server

import (
	"context"
	"net/http"
	"time"

	"bitbucket.org/mobio/go-logger"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

type AppServer struct {
	mobildaHost string

	Router *chi.Mux
	logger *logger.Logger
}

func NewAppServer(mobildaHost string, logger *logger.Logger) (*AppServer, error) {
	srv := &AppServer{
		mobildaHost: mobildaHost,
		logger:      logger,
	}

	if err := srv.init(); err != nil {
		return nil, err
	}

	return srv, nil
}

func (srv *AppServer) init() error {
	srv.Router = chi.NewRouter()
	srv.Router.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.CloseNotify,
		middleware.Timeout(60*time.Second),
		middleware.Heartbeat("/ping"),
	)

	return nil
}

func (srv *AppServer) Run(stop chan struct{}) {
	server := http.Server{Addr: srv.mobildaHost, Handler: srv.Router}

	go func() {
		srv.logger.Infof("Listening on http://%s", srv.mobildaHost)
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				srv.logger.Fatal(err)
			}
		}
	}()

	<-stop

	srv.logger.Info("Shutting down...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	server.Shutdown(ctx)
	srv.logger.Info("Server gracefully stopped...")
	stop <- struct{}{}

}
