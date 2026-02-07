package main

import (
	"context"
	"log"
	"simple_bank/api"
	db "simple_bank/db/sqlc"
	"simple_bank/util"

	_ "github.com/lib/pq"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server := api.NewServer(store)
	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
