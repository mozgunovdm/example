package dbase

import (
	"context"
	"database/sql"
	"os"
	"strconv"
	"testing"

	"github.com/go-kit/log"
	_ "github.com/lib/pq"
	"github.com/mozgunovdm/example/internal/pkg/employe"
)

func TestCreateEmploye(t *testing.T) {

	logger := log.NewLogfmtLogger(os.Stderr)

	var db *sql.DB
	{
		var err error
		db, err = sql.Open("postgres",
			"postgres://user:12345678@localhost:5433/employedb?sslmode=disable")
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
		EmployedAt: "2011-01-01",
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
			"postgres://user:12345678@localhost:5433/employedb?sslmode=disable")
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
		EmployedAt: "2020-05-05",
	}

	var id string
	ctx := context.Background()
	id, err = repository.CreateEmploye(ctx, emp)
	if err != nil {
		t.Fatalf("Could not create employe")
	}

	emp2 := employe.EmployeDB{
		Name:       "testName2",
		Job:        "testJob2",
		EmployedAt: "2000-01-01",
		HeadID:     id,
	}
	_, err = repository.CreateEmploye(ctx, emp2)
	if err != nil {
		t.Fatalf("Could not create sub employe %v", emp2)
	}

	empRes, err := repository.GetEmployeByID(ctx, id)
	if err != nil {
		t.Errorf(err.Error())
	}

	if (id != empRes.ID) ||
		(emp.Name != empRes.Name) ||
		(emp.EmployedAt != empRes.EmployedAt) ||
		(emp.Job != empRes.Job) ||
		(len(empRes.Employes) == 0) {
		t.Errorf("structs don't equil %v and %v", emp, empRes)
	}
}
