package cmdmain

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/mozgunovdm/example/internal/pkg/config"
	dbase "github.com/mozgunovdm/example/internal/pkg/database"
	empsvc "github.com/mozgunovdm/example/internal/pkg/implementation"
	"github.com/mozgunovdm/example/internal/pkg/transport"
	httptransport "github.com/mozgunovdm/example/internal/pkg/transport/http"
)

func InitDatabase(c *config.Config, l log.Logger) (*sql.DB, error) {
	url := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		c.Database.Type,
		c.Database.User,
		c.Database.Pass,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name)
	// Connect to the database
	db, err := sql.Open(c.Database.Type,
		url)
	if err != nil {
		level.Error(l).Log("exit", err)
		return db, err
	}
	level.Info(l).Log("msg", "database connected")
	return db, nil
}

func InitHttpServer(c *config.Config, l log.Logger) (*http.Server, error) {

	srv := &http.Server{}
	//Create database connection
	var db *sql.DB
	{
		var err error
		db, err = InitDatabase(c, l)
		if err != nil {
			level.Error(l).Log("exit", err)
			return srv, err
		}
	}

	var h http.Handler
	{
		repository, err := dbase.New(db, l)
		if err != nil {
			level.Error(l).Log("exit", err)
			return srv, err
		}
		endpoints := transport.MakeEndpoints(empsvc.NewService(repository, l))
		h = httptransport.NewService(endpoints, l)
	}

	//Create url addr
	httpAddr := fmt.Sprintf("%s:%s", c.HTTPServer.Host, c.HTTPServer.Port)

	level.Info(l).Log("transport", "HTTP", "addr", httpAddr)
	srv = &http.Server{
		Addr:    httpAddr,
		Handler: h,
	}
	return srv, nil
}

func Main() {
	//Create app logger
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "employe",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	if err := godotenv.Load(".env"); err != nil {
		level.Error(logger).Log("No .env files", err)
		return
	}
	config := config.NewConfig()

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		server, err := InitHttpServer(config, logger)
		if err == nil {
			errs <- server.ListenAndServe()
		} else {
			errs <- fmt.Errorf("%s", err)
		}
	}()

	level.Error(logger).Log("exit", <-errs)
}
