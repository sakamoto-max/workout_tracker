package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBConn *pgxpool.Pool

const DataBaseURL string = "postgresql://postgres:root@localhost:5432/WORKOUT_TRACKER"

func makeDBPool() (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), DataBaseURL)
	if err != nil{
		return nil, err
	}

	return conn, nil
}

func InitDB() {
	conn, err := makeDBPool()
	if err != nil{
		log.Fatalf("error while connecting to the DB : %v", err)
	}else {
		DBConn = conn
	}
}