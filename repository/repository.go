package repository

import (
	"context"
	"workout_tracker/database"
	"workout_tracker/models"

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

func CreateUserInDB(user models.User) error {
	
	_, err := database.DBConn.Exec(context.Background(), `
		INSERT INTO USERS(NAME, EMAIL, HASHED_PASSWORD, ROLE, CREATED_AT, UPDATED_AT)
		VALUES($1, $2, $3, $4, NOW(), NOW())
	`, user.Name, user.Email, user.Password, user.Role)
	if err != nil {
		return err
	}
	return nil
}

func GetUserFromDb(email string) (models.User, error) {
	var userDetailsFromDb models.User

	err := database.DBConn.QueryRow(context.Background(), `
		SELECT ID, NAME, EMAIL, ROLE, CREATED_AT, UPDATED_AT
		FROM USERS	
		WHERE EMAIL = $1
	`,email).Scan(&userDetailsFromDb.Id, &userDetailsFromDb.Name, &userDetailsFromDb.Email, &userDetailsFromDb.Role, &userDetailsFromDb.CreatedAt, &userDetailsFromDb.UpdatedAt)

	if err != nil {
		return userDetailsFromDb, err
	}

	userDetailsFromDb.Password = "confidential"

	return userDetailsFromDb, nil
}