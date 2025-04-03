package main

import (
	"bytes"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/GeoNet/kit/cfg"
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

func TestImport(t *testing.T) {
	setup(t)
	defer teardown()

	b, err := os.ReadFile("testdata/test.pb")
	if err != nil {
		t.Fatal(err)
	}

	// This takes 150 seconds on my computer
	if err = processProto(bytes.NewBuffer(b)); err != nil {
		t.Error(err)
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
