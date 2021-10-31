package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/mozgunovdm/example/employe"
)

var (
	ErrRepository = errors.New("unable to handle request")
)

type repository struct {
	db     *sql.DB
	logger log.Logger
}

// New returns a concrete repository backed by PostgreSQL
func New(db *sql.DB, logger log.Logger) (employe.Repository, error) {
	// return  repository
	return &repository{
		db:     db,
		logger: log.With(logger, "rep", "database"),
	}, nil
}

// CreateOrder inserts a new order and its order items into db
func (repo *repository) CreateOrder(ctx context.Context, employe employe.Employe) error {

	repo.logger.Log("function", "Called CreateOrder")
	// // Run a transaction to sync the query model.
	// err := crdb.ExecuteTx(ctx, repo.db, nil, func(tx *sql.Tx) error {
	// 	return createOrder(tx, employe)
	// })
	// if err != nil {
	// 	return err
	// }
	return nil
}

// func createEmploye(tx *sql.Tx, employe employe.Employe) error {

// 	// Insert employe into the "employe" table.
// 	sql := `
// 			INSERT INTO employe (name, job, employed_at)
// 			VALUES ($1,$2,$3,$4,$5)`
// 	_, err := tx.Exec(sql, employe.Name, employe.Job, employe.EmployedAt)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// GetOrderByID query the order by given id
func (repo *repository) GetEmployeByID(ctx context.Context, id string) (employe.Employe, error) {
	var employeById = employe.Employe{}
	if err := repo.db.QueryRowContext(ctx,
		"SELECT id, name, job, employed_at FROM employe WHERE id = $1",
		id).
		Scan(
			&employeById.ID, &employeById.Name, &employeById.Job, &employeById.EmployedAt,
		); err != nil {
		level.Error(repo.logger).Log("err", err.Error())
		return employeById, err
	}
	// ToDo: Query employe for head
	return employeById, nil
}

// Close implements DB.Close
func (repo *repository) Close() error {
	return repo.db.Close()
}
