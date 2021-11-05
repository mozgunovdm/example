package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
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

func (repo *repository) CreateEmploye(ctx context.Context, employe employe.EmployeDB) (string, error) {

	//  Insert Employe into the "employe" table.

	var id, sql string
	var err error
	if employe.HeadID == "" {
		sql = `SELECT "p_employe_add" ($1,$2,$3)`
		err = repo.db.QueryRowContext(ctx, sql, employe.Name, employe.Job, employe.EmployedAt).
			Scan(&id)
	} else {
		sql = `SELECT "p_employe_add" ($1,$2,$3,$4)`
		err = repo.db.QueryRowContext(ctx, sql, employe.Name, employe.Job, employe.EmployedAt, employe.HeadID).
			Scan(&id)
	}

	if err != nil {
		level.Error(repo.logger).Log("err", err.Error())
		return "0", err
	}
	return id, nil
}

// GetEmployeByID query the employe by given id
func (repo *repository) GetEmployeByID(ctx context.Context, id string) (employe.Employe, error) {
	var employeById = employe.Employe{}
	if err := repo.db.QueryRowContext(ctx,
		"SELECT id, name, job, employed_at FROM employe WHERE (id = $1)",
		id).
		Scan(
			&employeById.ID, &employeById.Name, &employeById.Job, &employeById.EmployedAt,
		); err != nil {
		level.Error(repo.logger).Log("err", err.Error())
		return employeById, err
	}
	if err := repo.getEmployesByID(ctx, &employeById); err != nil {
		level.Error(repo.logger).Log("err", err.Error())
		return employeById, err
	}
	return employeById, nil
}

func (repo *repository) getEmployesByID(ctx context.Context, head *employe.Employe) error {
	rows, err := repo.db.QueryContext(ctx,
		"SELECT id, name, job, employed_at FROM employe WHERE (head_id = $1)",
		head.ID)
	if err != nil {
		level.Error(repo.logger).Log("err", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var employeById = employe.Employe{}
		if err := rows.Scan(&employeById.ID, &employeById.Name, &employeById.Job, &employeById.EmployedAt); err != nil {
			level.Error(repo.logger).Log("err", err.Error())
			return err
		}
		if err := repo.getEmployesByID(ctx, &employeById); err != nil {
			level.Error(repo.logger).Log("err", err.Error())
			return err
		}
		head.Employes = append(head.Employes, employeById)
	}
	return nil
}

// Close implements DB.Close
func (repo *repository) Close() error {
	return repo.db.Close()
}
