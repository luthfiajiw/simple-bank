package main

import (
	"context"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	address  = "0.0.0.0:8080"
)

func main() {
	config, errConf := pgxpool.ParseConfig(dbSource)
	if errConf != nil {
		log.Fatalf("unable to parse config: %v:", errConf)
	}

	var err error

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("can't connect to db:", err)
	}

	store := db.NewStore(pool)
	server := api.NewServer(store)

	err = server.Start(address)
	if err != nil {
		log.Fatal("can't start server:", err)
	}
}
