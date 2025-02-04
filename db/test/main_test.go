package db_test

import (
	"context"
	"log"
	"os"
	db "simplebank/db/sqlc"
	"simplebank/utils"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *db.Queries
var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	env, errEnv := utils.LoadConfig("../..")
	if errEnv != nil {
		log.Fatalf("unable to load env config: %v:", errEnv)
	}

	config, errConf := pgxpool.ParseConfig(env.DBSource)
	if errConf != nil {
		log.Fatalf("unable to parse config: %v:", errConf)
	}

	var err error

	testPool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("can't connect to db:", err)
	}
	defer testPool.Close()

	testQueries = db.New(testPool)

	os.Exit(m.Run())
}
