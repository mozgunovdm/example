package middle

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/mozgunovdm/example/internal/pkg/employe"
)

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next employe.Service) employe.Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   employe.Service
	logger log.Logger
}

func (mw loggingMiddleware) Create(ctx context.Context, employe employe.EmployeDB) (id string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Create", "Employe name", employe.Name, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Create(ctx, employe)
}

func (mw loggingMiddleware) GetByID(ctx context.Context, id string) (employe employe.Employe, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetByID", "Employe ID", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetByID(ctx, id)
}

func (mw loggingMiddleware) Status(ctx context.Context) (s string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Status", "State", s, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Status(ctx)
}
