package main

import (
	"database/sql"
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

	res, err := db.Exec(sqlInsert, "test_ingest", "test_key3", "field1", time.Now(), "1.1")
	if err != nil {
		t.Error(err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		t.Error("failed to get number of rows affected:", err)
	}

	if n != 1 {
		t.Error("expected 1 affected got", n)
	}
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
