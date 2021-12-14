package dbase

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/mozgunovdm/example/internal/pkg/employe"
)

var (
	ErrRepository = errors.New("unable to handle request")
)

func subDate(s string) string {
	if len(s) == 0 || len(s) < 10 {
		return s
	}
	return s[0:10]
}

type repository struct {
	db     *sql.DB
	logger log.Logger
}

// New returns a concrete repository backed by PostgreSQL
func New(db *sql.DB, l log.Logger) (employe.Repository, error) {
	// return  repository
	return &repository{
		db:     db,
		logger: log.With(l, "rep", "database"),
	}, nil
}

// CreateEmploye query add new employe
func (repo *repository) CreateEmploye(ctx context.Context, emp employe.EmployeDB) (string, error) {

	//  Insert Employe into the "employe" table.
	id := "nil"

	//Check head id exist
	if len(strings.TrimSpace(emp.HeadID)) > 0 {
		sql := `SELECT id FROM employe WHERE (id = $1)`
		err := repo.db.QueryRowContext(ctx, sql, emp.HeadID).Scan(&id)
		if err != nil {
			level.Error(repo.logger).Log("err", err.Error())
			return id, err
		}
		//level.Error(repo.logger).Log("id", id, "err", err.Error())
	}

	// Insert new Employe
	sql := `INSERT INTO employe (name, job, employed_at) values ($1,$2,$3) returning id`
	err := repo.db.QueryRowContext(ctx, sql, emp.Name, emp.Job, emp.EmployedAt).Scan(&id)
	if err != nil {
		level.Error(repo.logger).Log("err", err.Error())
		return id, err
	}

	//Insert head of employe if needed
	if len(strings.TrimSpace(emp.HeadID)) > 0 {
		sql := `INSERT INTO relation_employes (stuff_id, head_id) values ($1, $2)`
		err := repo.db.QueryRowContext(ctx, sql, id, emp.HeadID).Err()
		if err != nil {
			level.Error(repo.logger).Log("err", err.Error())
			return id, err
		}
	}

	return id, nil
}

// GetEmployeByID query the employe by given id
func (repo *repository) GetEmployeByID(ctx context.Context, id string) (employe.Employe, error) {
	var emp = employe.Employe{}
	if err := repo.db.QueryRowContext(ctx,
		"SELECT id, name, job, employed_at FROM employe WHERE (id = $1)",
		id).
		Scan(
			&emp.ID, &emp.Name, &emp.Job, &emp.EmployedAt,
		); err != nil {
		level.Error(repo.logger).Log("err", err.Error())
		return emp, err
	}
	emp.EmployedAt = subDate(emp.EmployedAt)
	hLimit := 0
	if err := repo.getEmployesByID(ctx, &emp, hLimit); err != nil {
		level.Error(repo.logger).Log("err", err.Error())
		return emp, err
	}
	return emp, nil
}

//Find sub employe and add to array
func (repo *repository) getEmployesByID(ctx context.Context, h *employe.Employe, n int) error {
	//check limit of hierarchy
	if n++; n > 4 {
		return nil
	}

	rows, err := repo.db.QueryContext(ctx,
		"SELECT stuff_id FROM relation_employes WHERE (head_id = $1)",
		h.ID)
	if err != nil {
		level.Error(repo.logger).Log("err", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			level.Error(repo.logger).Log("err", err.Error())
			return err
		}

		var emp = employe.Employe{}
		if err := repo.db.QueryRowContext(ctx,
			"SELECT id, name, job, employed_at FROM employe WHERE (id = $1)",
			id).
			Scan(
				&emp.ID, &emp.Name, &emp.Job, &emp.EmployedAt,
			); err != nil {
			level.Error(repo.logger).Log("err", err.Error())
			return err
		}
		emp.EmployedAt = subDate(emp.EmployedAt)
		if err := repo.getEmployesByID(ctx, &emp, n); err != nil {
			level.Error(repo.logger).Log("err", err.Error())
			return err
		}
		h.Employes = append(h.Employes, emp)
	}
	return nil
}

// Close implements DB.Close
func (r *repository) Close() error {
	return r.db.Close()
}
