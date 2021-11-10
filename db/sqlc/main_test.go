package db

import (
	"database/sql"
	"gobank/util"
	"log"
	"os"
	"testing"
	// _ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	cfg, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(cfg.DBDriver, cfg.DBSource)

	if err != nil {
		log.Fatalf("can't connect to db: %s", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
