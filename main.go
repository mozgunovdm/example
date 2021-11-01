package main

import (
	"database/sql"
	"flag"
	"fmt"

	//	"github.com/shijuvar/gokit-examples/services/order/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	//	kithttp "github.com/go-kit/kit/transport/http"
	_ "github.com/lib/pq"
	//"github.com/shijuvar/gokit-examples/pkg/oc"

	"github.com/mozgunovdm/example/employe"
	"github.com/mozgunovdm/example/employe/database"
	employesvc "github.com/mozgunovdm/example/employe/implementation"
	"github.com/mozgunovdm/example/employe/transport"
	httptransport "github.com/mozgunovdm/example/employe/transport/http"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

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

	var db *sql.DB
	{
		var err error
		// Connect to the "ordersdb" database
		db, err = sql.Open("postgres",
			"postgres://postgres:060701@localhost:5432/employedb?sslmode=disable")
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
		level.Info(logger).Log("msg", "database connected")
	}

	// Create Employe Service
	var svc employe.Service
	{
		repository, err := database.New(db, logger)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
		svc = employesvc.NewService(repository, logger)
	}

	var h http.Handler
	{
		endpoints := transport.MakeEndpoints(svc)
		h = httptransport.NewService(endpoints, logger)
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
		server := &http.Server{
			Addr:    *httpAddr,
			Handler: h,
		}
		errs <- server.ListenAndServe()
	}()

	level.Error(logger).Log("exit", <-errs)
}
