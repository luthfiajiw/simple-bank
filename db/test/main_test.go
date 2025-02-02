package db_test

import (
	"context"
	"log"
	"os"
	db "simplebank/db/sqlc"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *db.Queries
var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	config, errConf := pgxpool.ParseConfig(dbSource)
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
