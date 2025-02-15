package main

import (
	"context"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	env, errEnv := utils.LoadConfig(".")
	if errEnv != nil {
		log.Fatalf("unable to load env config: %v:", errEnv)
	}

	config, errConf := pgxpool.ParseConfig(env.DBSource)
	if errConf != nil {
		log.Fatalf("unable to parse config: %v:", errConf)
	}

	var err error

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("can't connect to db:", err)
	}

	store := db.NewStore(pool)
	server, err := api.NewServer(env, store)
	if err != nil {
		log.Fatal("can't create server:", err)
	}

	err = server.Start(env.ServerAddress)
	if err != nil {
		log.Fatal("can't start server:", err)
	}
}
