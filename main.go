package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"tutorial.sqlc.dev/app/api"
	db "tutorial.sqlc.dev/app/db/sqlc"
	"tutorial.sqlc.dev/app/db/util"
)

// const (
// 	dbDriver  = "postgres"
// 	dbSource  = "postgres://root:fHWFyt98gPR51h3NxjcroWoIscjt7QOb@dpg-cp649mmn7f5s73a6r8ag-a.oregon-postgres.render.com/simple_bank_7qc2"
// 	adrServer = "0.0.0.0:8080"
// )

func main() {

	var err error

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewsStore(conn)
	server := api.NewServer(store)

	err = server.Starts(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
