package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thanhphuocnguyen/go-simple-bank/api"
	db "github.com/thanhphuocnguyen/go-simple-bank/db/sqlc"
	"github.com/thanhphuocnguyen/go-simple-bank/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	dbSource := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create server: %v", err)
	}

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
