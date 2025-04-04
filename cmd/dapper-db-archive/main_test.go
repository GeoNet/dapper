package main

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/GeoNet/kit/cfg"
	_ "github.com/lib/pq"
)

func setTestEnvVariables(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_CONN_TIMEOUT", "5")
	t.Setenv("DB_USER", "dapper_w")
	t.Setenv("DB_PASSWD", "test")
	t.Setenv("DB_NAME", "dapper")
	t.Setenv("DB_SSLMODE", "disable")
	t.Setenv("DB_MAX_IDLE_CONNS", "30")
	t.Setenv("DB_MAX_OPEN_CONNS", "30")
	t.Setenv("DB_CONN_TIMEOUT", "5")
	t.Setenv("AWS_REGION", "ap-southeast-2")
}

// Note: Must ran dapper/etc/script/initdb-test.sh before running these tests
func TestSQL(t *testing.T) {
	setup(t)
	defer teardown()

	if err := checkQuery(sqlSelectArchive, 5, "test_db_archive"); err != nil {
		t.Error(err)
	}

	if err := checkQuery(sqlUpdateArchive, 2, "test_db_archive", "test_key1", time.Now().Add(-1*time.Hour), time.Now()); err != nil {
		t.Error(err)
	}

	if err := checkQuery(sqlDeleteArchive, 1, "test_db_archive"); err != nil {
		t.Error(err)
	}
}

func checkQuery(sql string, nexp int, args ...interface{}) error {
	res, err := db.Exec(sql, args...)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of rows affected: %v", err)
	}
	if n != int64(nexp) {
		return fmt.Errorf("expected %d affected, got %d", nexp, n)
	}

	return nil
}

func setup(t *testing.T) {
	setTestEnvVariables(t)
	var err error
	p, err := cfg.PostgresEnv()
	if err != nil {
		log.Fatalf("error reading DB config from the environment vars: %v", err)
	}

	db, err = sql.Open("postgres", p.Connection())
	if err != nil {
		log.Fatalf("error with DB config: %v", err)
	}
}

func teardown() {
	db.Close()
}
