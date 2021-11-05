package database

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/go-kit/log"

	"github.com/mozgunovdm/example/employe"

	"strconv"
)

func TestCreateEmploye(t *testing.T) {

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)

	var db *sql.DB
	{
		var err error
		db, err = sql.Open("postgres",
			"postgres://testuser:12345678@localhost:5433/test_db?sslmode=disable")
		if err != nil {
			t.Fatalf("Db no connection " + err.Error())
		}
	}
	defer db.Close()

	repository, err := New(db, logger)
	if err != nil {
		t.Fatalf("Could not create repository")
	}

	emp := employe.EmployeDB{
		Name:       "testName",
		Job:        "testJob",
		EmployedAt: "0000-00-00",
	}

	ctx := context.Background()

	var id string
	id, err = repository.CreateEmploye(ctx, emp)
	if err != nil {
		t.Errorf(err.Error())
	}
	id_int, errC := strconv.Atoi(id)
	if errC != nil {
		t.Errorf(errC.Error())
	}
	if id_int <= 0 {
		t.Errorf("Error returned id %v", id_int)
	}
}

func TestGetEmployeByID(t *testing.T) {

	logger := log.NewLogfmtLogger(os.Stderr)

	var db *sql.DB
	{
		var err error
		db, err = sql.Open("postgres",
			"postgres://testuser:12345678@localhost:5433/test_db?sslmode=disable")
		if err != nil {
			t.Fatalf("Db no connection")
		}
	}
	defer db.Close()

	repository, err := New(db, logger)
	if err != nil {
		t.Fatalf("Could not create repository")
	}

	emp := employe.Employe{
		ID:         "1",
		Name:       "testName",
		Job:        "testJob",
		EmployedAt: "0000-00-00",
	}

	ctx := context.Background()
	sql := `SELECT "p_employe_add" ($1,$2,$3,$4)`
	err = db.QueryRowContext(ctx, sql, emp.Name, emp.Job, emp.EmployedAt, emp.ID).
		Scan(&emp.ID)
	if err != nil {
		t.Errorf(err.Error())
	}

	employeResult, errResult := repository.GetEmployeByID(ctx, emp.ID)
	if errResult != nil {
		t.Errorf(errResult.Error())
	}

	if (emp.ID != employeResult.ID) ||
		(emp.Name != employeResult.Name) ||
		(emp.EmployedAt != employeResult.EmployedAt) ||
		(emp.Job != employeResult.Job) {
		t.Errorf("structs don't equil")
	}
}
