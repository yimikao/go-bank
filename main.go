package main

import (
	"database/sql"
	"gobank/api"
	db "gobank/db/sqlc"
	"gobank/util"
	"log"

	_ "github.com/lib/pq"
)

func main() {

	cfg, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open("postgres", cfg.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(cfg.ServerAddr)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
