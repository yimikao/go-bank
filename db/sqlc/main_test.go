package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbURI    = "postgres://yinka:supadiski@localhost:5432/go_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbURI)

	if err != nil {
		log.Fatalf("can't connect to db: %s", err)
	}
	testQueries = New(conn)
	os.Exit(m.Run())
}
