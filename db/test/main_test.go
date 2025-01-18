package db_test

import (
	"context"
	"log"
	"os"
	db "simplebank/db/sqlc"
	"testing"

	"github.com/jackc/pgx/v5"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *db.Queries

func TestMain(m *testing.M)  {
	conn, err := pgx.Connect(context.Background(), dbSource)
	if err != nil {
		log.Fatal("can't connect to db:", err)	
	}
	defer conn.Close(context.Background())

	testQueries = db.New(conn)

	os.Exit(m.Run())
}