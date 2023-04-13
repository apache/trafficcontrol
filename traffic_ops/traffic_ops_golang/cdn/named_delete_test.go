package cdn

import (
	"testing"

	"github.com/jmoiron/sqlx"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestDeleteCDNByName(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM cdn").WithArgs("cdn1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = deleteCDNByName(db.MustBegin().Tx, "cdn1")
	if err != nil {
		t.Fatalf("no error expected while deleting CDN by name, but got: %v", err)
	}
}

func TestCDNUsed(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	cols := []string{"?column?"}
	rows := sqlmock.NewRows(cols)
	rows.AddRow(5)

	mock.ExpectBegin()
	mock.ExpectQuery("WITH cdn_id as").WithArgs("cdn1").WillReturnRows(rows)
	mock.ExpectCommit()

	unused, err := cdnUnused(db.MustBegin().Tx, "cdn1")
	if err != nil {
		t.Fatalf("no error expected in call to cdnUnused, but got: %v", err)
	}
	if unused {
		t.Errorf("expected CDN to be used, but is unused")
	}
}
