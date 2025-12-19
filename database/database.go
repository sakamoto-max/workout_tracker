package database

import (
	"context"
	"fmt"
	"log"
	"workout_tracker/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBConn *pgxpool.Pool

// const DataBaseURL string = "postgresql://postgres:root@localhost:5432/workout"


func makeDBPool() (*pgxpool.Pool, error) {
	
	DataBaseURL := fmt.Sprintf("%v://%v:%v@%v:%v/%v", config.Config.Db, config.Config.DbOwnerName, config.Config.DbPassword, config.Config.DbHost, config.Config.DbPort, config.Config.DbName)

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