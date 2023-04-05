package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/zaid13/simplebank/api"
	db "github.com/zaid13/simplebank/db/sqlc"
	"github.com/zaid13/simplebank/db/util"
	"log"
)

func main() {

	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to DB", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create log ")
	}

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server", err)

	}

}
