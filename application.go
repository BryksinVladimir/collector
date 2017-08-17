package mobilda

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"mobilda/client"
	acc "mobilda/collectors/accounts"
	"mobilda/consts"
	"mobilda/errors"
	"mobilda/model"
	"mobilda/server"

	"bitbucket.org/mobio/go-cache"
	"bitbucket.org/mobio/go-config"
	"bitbucket.org/mobio/go-dbmanager"
	"bitbucket.org/mobio/go-logger"
	"bitbucket.org/mobio/go-logger/hooks/mail"
	"bitbucket.org/mobio/go-scheduler"
	"github.com/sirupsen/logrus"
	"gopkg.in/pg.v5"
)

type Application struct {
	ctx       context.Context
	env       string
	configDir string

	server    *server.AppServer
	config    *config.Config
	logger    *logger.Logger
	dbmanager *dbmanager.DbManager
	scheduler *scheduler.Scheduler
	cache     *cache.Cache
	mobClient *client.MobildaClient
	accounts  []*model.Account

	quit chan os.Signal

	ShutdownFunc func(ctx context.Context)
}

func NewApplication(configDir, env string) (*Application, error) {
	app := &Application{
		env:       env,
		configDir: configDir,
		quit:      make(chan os.Signal, 1),
	}

	err := app.init()

	return app, err
}

func (app *Application) init() error {
	// init quit channel
	signal.Notify(app.quit, os.Interrupt, os.Kill)

	//Init logger
	if err := app.initLogger(); err != nil {
		return err
	}

	//Init config
	if err := app.initConfig(); err != nil {
		return err
	}

	//Init mail hook
	if err := app.initMail(); err != nil {
		return err
	}

	//Init pg connection
	app.logger.Info("Connecting to PostgreSQL...")
	if err := app.initPostgresDb(); err != nil {
		return err
	}

	app.logger.WithFields(logrus.Fields{
		"host": app.dbmanager.Options().Addr,
		"db":   app.dbmanager.Options().Database,
		"user": app.dbmanager.Options().User,
	}).Infof("Connected to PostgreSQL on %s", app.dbmanager.Options().Addr)

	//Init cache
	if err := app.initCache(); err != nil {
		return err
	}

	//Init scheduler
	if err := app.initScheduler(); err != nil {
		return err
	}

	//Init mobilda client
	if err := app.initClient(); err != nil {
		return err
	}

	if err := app.initContext(); err != nil {
		return err
	}

	if err := app.updateAccounts(); err != nil {
		return err
	}

	// add collectors
	app.logger.Info("Adding collectors to scheduler...")
	app.addCollectors()

	if err := app.initAppServer(); err != nil {
		return err
	}

	return nil
}

func (app *Application) initLogger() error {
	app.logger = logger.NewLogger()
	return nil
}

func (app *Application) initConfig() error {
	app.config = config.LoadConfig(app.configDir, app.env, "yaml", app.logger)
	app.logger.SetLevel(app.config.GetString(consts.Log_Level_Key))

	return nil
}

func (app *Application) initPostgresDb() error {
	//Init DbManager component
	opts := &dbmanager.Options{
		PgOpts: pg.Options{
			Addr:     app.config.GetString("postgres.addr"),
			User:     app.config.GetString("postgres.user"),
			Password: app.config.GetString("postgres.password"),
			Database: app.config.GetString("postgres.databse"),
		},
	}
	opts.LogQuery = app.config.GetBool("postgres.log_query")
	app.dbmanager = dbmanager.NewDbManager(opts, app.logger)

	_, err := app.dbmanager.Exec("SELECT 1")
	if err != nil {
		return err
	}

	return nil
}

func (app *Application) initCache() error {
	app.cache = cache.NewCache()
	return nil
}

func (app *Application) initScheduler() error {
	app.scheduler = scheduler.NewScheduler(app.logger)
	return nil
}

func (app *Application) initClient() error {
	app.config.UnmarshalKey(consts.Mobilda_Key, &app.accounts)

	if len(app.accounts) == 0 {
		return errors.ErrCantGetAccountsFromConfig
	}

	app.mobClient = client.NewMobildaClient(app.accounts, nil, time.Second*60, 0, app.logger)

	return nil
}

func (app *Application) updateAccounts() error {
	accountCollector := acc.NewAccountsCollector(app.ctx)
	accountCollector.Run()

	return nil
}

func (app *Application) initContext() error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, consts.Logger_Component_Key, app.logger)
	ctx = context.WithValue(ctx, consts.DbManager_Component_Key, app.dbmanager)
	ctx = context.WithValue(ctx, consts.Cache_Component_Key, app.cache)
	ctx = context.WithValue(ctx, consts.Config_Component_Key, app.config)
	ctx = context.WithValue(ctx, consts.Scheduler_Component_Key, app.scheduler)
	ctx = context.WithValue(ctx, consts.MobildaClient_Component_Key, app.mobClient)
	ctx = context.WithValue(ctx, consts.Accounts_Key, app.accounts)
	app.ctx = ctx

	return nil
}

func (app *Application) initMail() error {
	hook, err := logrus_mail.NewMailHook(
		app.config.GetString("mailer.appname"),
		app.config.GetString("mailer.host"),
		app.config.GetInt("mailer.port"),
		app.config.GetString("mailer.from"),
		app.config.GetString("mailer.to"),
	)
	if err != nil {
		return err
	}

	app.logger.Hooks.Add(hook)

	return nil
}

func (app *Application) initAppServer() error {
	as, err := server.NewAppServer(app.config.GetString(consts.ConfMobildaHost), app.logger)
	// Middleware to put application components into Request.Context
	as.Router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(app.ctx))
		})
	})

	//Init router
	as.InitRouter()
	app.server = as
	return err
}

func (app *Application) Env() string {
	return app.env
}

func (app *Application) ConfigDir() string {
	return app.configDir
}

func (app *Application) Run() {
	stop := make(chan struct{})

	//Start scheduler
	app.scheduler.Start()

	go app.server.Run(stop)

	<-app.quit
	stop <- struct{}{}
	<-stop
	app.shutdown()
}

func (app *Application) shutdown() {
	app.dbmanager.Close()
	app.logger.Info("PostgreSQL connection closed...")
}
