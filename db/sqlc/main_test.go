package db

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/thanhphuocnguyen/go-simple-bank/utils"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	config, err := utils.LoadConfig("../..")
	if err != nil {
		panic(err)
	}

	dbSource := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	testDB, err = pgxpool.New(context.Background(), dbSource)
	if err != nil {
		panic(err)
	}
	defer testDB.Close()
	testQueries = New(testDB)
	fmt.Println("Connected to database")

	os.Exit(m.Run())
}
