package cmdmain

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	kitoc "github.com/go-kit/kit/tracing/opencensus"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/mozgunovdm/example/pkg/oc"

	"github.com/mozgunovdm/example/internal/pkg/config"
	dbase "github.com/mozgunovdm/example/internal/pkg/database"
	empsvc "github.com/mozgunovdm/example/internal/pkg/implementation"
	"github.com/mozgunovdm/example/internal/pkg/middle"
	"github.com/mozgunovdm/example/internal/pkg/transport"
	httptransport "github.com/mozgunovdm/example/internal/pkg/transport/http"
)

func InitDatabase(cnf *config.Config, ctx context.Context, l log.Logger) (*sql.DB, error) {
	url := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		cnf.Database.Type,
		cnf.Database.User,
		cnf.Database.Pass,
		cnf.Database.Host,
		cnf.Database.Port,
		cnf.Database.Name)
	// Connect to the database
	db, err := sql.Open(cnf.Database.Type,
		url)
	if err != nil {
		return db, err
	}
	level.Info(l).Log("db-type", cnf.Database.Type, "addr", cnf.Database.Host, "user", cnf.Database.User, "port", cnf.Database.Port, "db-name", cnf.Database.Name)
	if err := db.PingContext(ctx); err != nil {
		return db, err
	}
	level.Info(l).Log("msg", "database connected")

	return db, nil
}

func Main() {
	//defer oc.Setup("order").Close()
	ctx := context.Background()
	//Create app logger
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"Service", "employe",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	if err := godotenv.Load(); err != nil {
		level.Error(logger).Log("No .env files", err)
		panic("No .env files")
	}
	cnf := config.NewConfig()

	//Create database connection
	var db *sql.DB
	{
		var err error
		db, err = InitDatabase(cnf, ctx, logger)
		if err != nil {
			level.Error(logger).Log("exit", err)
			panic(err)
		}
	}
	defer db.Close()

	var h http.Handler
	{
		repository, err := dbase.New(db, logger)
		if err != nil {
			level.Error(logger).Log("exit", err)
			panic(err)
		}

		// Create Go kit endpoints for the Employe Service
		// Then decorates with endpoint middlewares
		endpoints := transport.MakeEndpoints(middle.LoggingMiddleware(logger)(empsvc.NewService(repository, logger)))
		// Add endpoint level middlewares
		// Trace server side endpoints with open census
		endpoints = transport.Endpoints{
			Create:  oc.ServerEndpoint("Create")(endpoints.Create),
			GetByID: oc.ServerEndpoint("GetByID")(endpoints.GetByID),
			Status:  oc.ServerEndpoint("ChangeStatus")(endpoints.Status),
		}
		ocTracing := kitoc.HTTPServerTrace()
		serverOptions := []kithttp.ServerOption{ocTracing}
		h = httptransport.NewService(endpoints, serverOptions, logger)
	}

	//Create url addr
	httpAddr := fmt.Sprintf("%s:%s", cnf.HTTPServer.Host, cnf.HTTPServer.Port)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", httpAddr)
		server := &http.Server{
			Addr:    httpAddr,
			Handler: h,
		}
		errs <- server.ListenAndServe()
	}()

	level.Error(logger).Log("exit", <-errs)
}
