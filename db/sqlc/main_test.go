package db

import (
	"database/sql"
	"gobank/util"
	"log"
	"os"
	"testing"
	// _ "github.com/lib/pq"
)

var testQueries *Queries //A test Queries object(Test Database wrapping either sql.DB/Tx)
var testDB *sql.DB

// main entry point of all unit tests inside 1 specific package
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
