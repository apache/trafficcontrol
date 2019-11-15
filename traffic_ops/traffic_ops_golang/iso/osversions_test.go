package iso

import (
	"context"
	"path"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestOsversionsCfgPath(t *testing.T) {
	cases := []struct {
		name           string
		parameterValue string
		expectedPath   string
	}{
		{
			"default",
			"", // No db parameter entry
			path.Join(cfgDefaultDir, cfgFilename),
		},
		{
			"override",
			"/this/is/not/the/default/",
			path.Join("/this/is/not/the/default", cfgFilename),
		},
		{
			"override-cwd",
			".", // CWD
			cfgFilename,
		},
	}

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() err: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dbCtx, cancel := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
			defer cancel()

			// Setup mock DB to return rows for SELECT query on parameter table.
			// If parameterValue is empty, no rows will be returned.
			mock.ExpectBegin()
			cols := []string{"value"}
			rows := sqlmock.NewRows(cols)
			if tc.parameterValue != "" {
				rows = rows.AddRow(tc.parameterValue)
			}
			mock.ExpectQuery("SELECT").WillReturnRows(rows)
			mock.ExpectCommit()

			tx, err := db.BeginTxx(dbCtx, nil)
			if err != nil {
				t.Fatalf("BeginTxx() err: %v", err)
			}
			defer tx.Commit()

			got, err := osversionCfgPath(tx)

			if err != nil {
				t.Fatalf("osversionCfgPath() err: %v, expected: nil", err)
			}

			if got != tc.expectedPath {
				t.Fatalf("osversionCfgPath(): %q, expected: %q", got, tc.expectedPath)
			}
			t.Logf("osversionCfgPath(): %q", got)
		})
	}
}
