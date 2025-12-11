package repository

import (
	"context"
	"workout_tracker/database"

	"github.com/jackc/pgx/v5"
)

func GetAllExercisesFromDB() (pgx.Rows, error) {

	rows, err := database.DBConn.Query(context.Background(), `
		SELECT * FROM EXERCISES
	`)
	if err != nil {
		return rows, err
	}

	return rows, nil
}
