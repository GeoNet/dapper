package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/GeoNet/kit/cfg"
	wt "github.com/GeoNet/kit/weft/wefttest"
	"github.com/lib/pq"
)

var (
	// temporary disable TestRoutes
	testServer *httptest.Server //nolint: unused
)

type dbTest struct {
	id   string
	sql  string
	exp  int
	args []interface{}
}

func setTestEnvVariables(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_CONN_TIMEOUT", "5")
	t.Setenv("DB_USER", "dapper_r")
	t.Setenv("DB_PASSWD", "test")
	t.Setenv("DB_NAME", "dapper")
	t.Setenv("DB_SSLMODE", "disable")
	t.Setenv("DB_MAX_IDLE_CONNS", "30")
	t.Setenv("DB_MAX_OPEN_CONNS", "30")
	t.Setenv("DB_CONN_TIMEOUT", "5")
	t.Setenv("AWS_REGION", "ap-southeast-2")
	t.Setenv("DOMAINS", "test_api")
	t.Setenv("DOMAIN_BUCKETS", "tf-dev-geonet-cache")
	t.Setenv("DOMAIN_PREFIXES", "dapper/fmp-p2")
}

// Note: Must ran dapper/etc/script/initdb-test.sh before running these tests
func TestMetaDB(t *testing.T) {
	setup(t)
	defer teardown()

	tests := []dbTest{
		{id: wt.L(), sql: sqlReverseFindKey, exp: 1, args: []interface{}{"test_api", "hostname", "rfap5g-soundstage", time.Now().UTC()}},
		{id: wt.L(), sql: sqlListKeys, exp: 3, args: []interface{}{"test_api", time.Now().UTC()}},
		{id: wt.L(), sql: sqlMetaFields, exp: 11, args: []interface{}{"test_api", time.Now().UTC()}},
		{id: wt.L(), sql: sqlTag, exp: 2, args: []interface{}{"test_api", "5G", time.Now().UTC()}},
		{id: wt.L(), sql: sqlILike, exp: 6, args: []interface{}{"test_api", "RFAP5G-soundstage", time.Now().UTC()}},
		{id: wt.L(), sql: sqlAny, exp: 6, args: []interface{}{"test_api", pq.Array([]string{"wance-avalonlab", "rfap5g-soundstage"}), time.Now().UTC()}},
		{id: wt.L(), sql: sqlGeomILike, exp: 1, args: []interface{}{"test_api", "RFAP5G-soundstage", time.Now().UTC()}},
		{id: wt.L(), sql: sqlGeomAny, exp: 1, args: []interface{}{"test_api", pq.Array([]string{"wance-avalonlab", "rfap5g-soundstage"}), time.Now().UTC()}},
		{id: wt.L(), sql: sqlRelILike, exp: 1, args: []interface{}{"test_api", "RF2soundstage-towai", time.Now().UTC()}},
		{id: wt.L(), sql: sqlRelAny, exp: 1, args: []interface{}{"test_api", pq.Array([]string{"rf2soundstage-towai", "rfap5g-soundstage"}), time.Now().UTC()}},
	}

	for _, test := range tests {
		if err := checkQuery(test.sql, test.exp, test.args...); err != nil {
			t.Error(test.id, err)
		}
	}
}

func TestDataDB(t *testing.T) {
	setup(t)
	defer teardown()

	tests := []dbTest{
		{id: wt.L(), sql: sqlDataSpan, exp: 1, args: []interface{}{"test_api", "rfap5g-soundstage", "temperature", time.Now().Add(-24 * time.Hour).UTC(), time.Now().UTC()}},
		{id: wt.L(), sql: sqlFields, exp: 2, args: []interface{}{"test_api", "rfap5g-soundstage"}},
		{id: wt.L(), sql: sqlDataLatest, exp: 1, args: []interface{}{"test_api", "rfap5g-soundstage", "temperature", 2}},
		{id: wt.L(), sql: sqlCacheLatest, exp: 2, args: []interface{}{}},
	}

	for _, test := range tests {
		if err := checkQuery(test.sql, test.exp, test.args...); err != nil {
			t.Error(test.id, err)
		}
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

// TODO: fix this test - this test requires a live S3 bucket
// func TestRoutes(t *testing.T) {
// 	setup(t)
// 	defer func() {
// 		testServer.Close()
// 		teardown()
// 	}()

// 	var err error

// 	s3Client, err = s3.New()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	initHandlers()
// 	initVars()
// 	if err = cacheLatest(); err != nil {
// 		log.Printf("error caching latest tables: %v", err)
// 	}

// 	routes := wt.Requests{
// 		{ID: wt.L(), Accept: CONTENT_TYPE_JSON, Content: CONTENT_TYPE_JSON, URL: "/meta/test_api/entries"},
// 		{ID: wt.L(), Accept: CONTENT_TYPE_JSON, Content: CONTENT_TYPE_JSON, URL: "/meta/test_api/entries?aggregate=locality"},
// 		{ID: wt.L(), Accept: CONTENT_TYPE_JSON, Content: CONTENT_TYPE_JSON, URL: "/meta/test_api/entries?key=rfap5g-soundstage"},
// 		{ID: wt.L(), Accept: CONTENT_TYPE_JSON, Content: CONTENT_TYPE_JSON, URL: "/meta/test_api/entries?tags=5g"},
// 		{ID: wt.L(), Accept: CONTENT_TYPE_JSON, Content: CONTENT_TYPE_JSON, URL: "/meta/test_api/entries?query=hostname=rfap5g-soundstage"},
// 		{ID: wt.L(), Accept: CONTENT_TYPE_GEOJSON, Content: CONTENT_TYPE_GEOJSON, URL: "/meta/test_api/entries?query=hostname=rfap5g-soundstage"},
// 		{ID: wt.L(), Accept: CONTENT_TYPE_JSON, Content: CONTENT_TYPE_JSON, URL: "/meta/test_api/list"},
// 		{ID: wt.L(), Accept: CONTENT_TYPE_PROTOBUF, Content: CONTENT_TYPE_PROTOBUF, URL: "/meta/test_api/entries?query=hostname=rfap5g-soundstage"},
// 		{ID: wt.L(), Accept: CONTENT_TYPE_JSON, Content: CONTENT_TYPE_JSON, URL: "/data/test_api?key=all"},
// 		{ID: wt.L(), Accept: CONTENT_TYPE_JSON, Content: CONTENT_TYPE_JSON, URL: "/data/test_api?key=rfap5g-soundstage&latest=2"},
// 		// Note the test below might fail if test data were inserted 24 hours before
// 		{ID: wt.L(), Accept: CONTENT_TYPE_JSON, Content: CONTENT_TYPE_JSON, URL: "/data/test_api?key=rfap5g-soundstage&fields=temperature&starttime=" + time.Now().UTC().Add(-24*time.Hour).Format(time.RFC3339) + "&endtime=" + time.Now().UTC().Format(time.RFC3339)},
// 	}

// 	testServer = httptest.NewServer(inbound(mux))

// 	for _, r := range routes {
// 		if _, err := r.Do(testServer.URL); err != nil {
// 			t.Error(err)
// 		}
// 	}
// }

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
