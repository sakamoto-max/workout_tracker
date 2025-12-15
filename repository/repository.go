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
	`, email).Scan(&userDetailsFromDb.Id, &userDetailsFromDb.Name, &userDetailsFromDb.Email, &userDetailsFromDb.Role, &userDetailsFromDb.CreatedAt, &userDetailsFromDb.UpdatedAt)

	if err != nil {
		return userDetailsFromDb, err
	}

	userDetailsFromDb.Password = "confidential"

	return userDetailsFromDb, nil
}

func GetHashedPassFromDB(email string) (string, error) {
	var hashedpass string
	err := database.DBConn.QueryRow(context.Background(), `
		SELECT HASHED_PASSWORD FROM USERS
		WHERE EMAIL = $1
	`, email).Scan(&hashedpass)
	if err != nil {
		return "", err
	}

	return hashedpass, nil
}

// user -> create plan -> plan_name -> will give-> plan_id
// exercise_tracker -> insert all the exercises
// get all the ids
// insert plan_id, all the exercise tracker ids in plan

func CreateAPlanInWorkOutPlans(user_id int, planName string) (int, error) {
	var planId int

	err := database.DBConn.QueryRow(context.Background(), `
		INSERT INTO WORKOUTPLANS(USER_ID, PLAN_NAME)
		VALUES($1, $2)
	`, user_id, planName).Scan(&planId)

	if err != nil {
		return planId, err
	}
	return planId, nil
}

func InsertExercisesIntoTracker(exerciseName []string) ([]int, error) {

	var exerciseTracerIds []int

	for _, v := range exerciseName {

		var exerciseTracerId int

		err := database.DBConn.QueryRow(context.Background(), `
			INSERT INTO EXERCISE_TRACKER(EXERCISE_NAME)
			VALUES($1)
			RETURNING ID
		`, v).Scan(&exerciseTracerId)

		if err != nil {
			return exerciseTracerIds, err
		}

		exerciseTracerIds = append(exerciseTracerIds, exerciseTracerId)
	}

	return exerciseTracerIds, nil
}

func InsertIntoPlan(planId int, exerciserTrackerIds []int) error {

	for _, v := range exerciserTrackerIds {
		_, err := database.DBConn.Exec(context.Background(), `
			INSERT INTO PLAN(PLAN_ID, EXERCISE_TRACKER_ID)
			VALUES($1, $2)
		`, planId, v)

		if err != nil {
			return err
		}
	}
	return nil
}

func EeAaO(userId int, planName string, exerciseNames []string) error {

	var planId int

	var exerciseTrackerIds []int

	trnx, err := database.DBConn.Begin(context.Background())
	if err != nil {
		return err
	}

	err = trnx.QueryRow(context.Background(), `
		INSERT INTO WORKOUTPLANS(USER_ID, PLAN_NAME)
		VALUES($1, $2)
		RETURNING ID
	`, userId, planName).Scan(&planId)
	if err != nil {
		return err
	}

	for _, v := range exerciseNames {

		var exerciseTracerId int

		err := trnx.QueryRow(context.Background(), `
			INSERT INTO EXERCISE_TRACKER(EXERCISE_NAME)
			VALUES($1)
			RETURNING ID
		`, v).Scan(&exerciseTracerId)

		if err != nil {
			return err
		}

		exerciseTrackerIds = append(exerciseTrackerIds, exerciseTracerId)
	}

	for _, v := range exerciseTrackerIds {
		_, err := trnx.Exec(context.Background(), `
			INSERT INTO PLAN(PLAN_ID, EXERCISE_TRACKER_ID)
			VALUES($1, $2)
		`, planId, v)

		if err != nil {
			return err
		}
	}

	err = trnx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func GetUserIdFromDB(email string) (int, error) {

	var userId int
	err := database.DBConn.QueryRow(context.Background(), `
		SELECT ID FROM USERS
		WHERE EMAIL = $1
	`, email).Scan(&userId)
	if err != nil {
		return userId, err
	}

	return userId, nil
}
func GetUserRoleFromDB(email string) (string, error) {

	var userRole string
	err := database.DBConn.QueryRow(context.Background(), `
		SELECT ROLE FROM USERS
		WHERE EMAIL = $1
	`, email).Scan(&userRole)
	if err != nil {
		return userRole, err
	}

	return userRole, nil
}

func InsertANewExerciseInDB(exerciseName string, exercisetype string, bodyPart string) (models.Exercise, error) {

	var exercise models.Exercise
	err := database.DBConn.QueryRow(context.Background(), `
		INSERT INTO EXERCISES(EXERCISE_NAME, TYPE, BODY_PART)
		VALUES($1, $2, $3)
		RETURNING ID, EXERCISE_NAME, TYPE, BODY_PART
	`, exerciseName, exercisetype, bodyPart).Scan(&exercise.Id, &exercise.ExerciseName, &exercise.Type, &exercise.BodyPart)

	if err != nil {
		return exercise, err
	}

	return exercise, nil

}

func DeleteExerciseFromDb(exerciseName string, exerciseId int) error {

	if exerciseId == 0 {
		_, err := database.DBConn.Exec(context.Background(), `
			DELETE FROM EXERCISES
			WHERE EXERCISE_NAME = $1		
		`,exerciseName)
		
		if err != nil {
			return err
		}
	}else {

		trnx, err := database.DBConn.Begin(context.Background())
		if err != nil{
			return err
		}
	
		_, err = trnx.Exec(context.Background(), `
			DELETE FROM EXERCISES 
			WHERE EXERCISE_NAME = $1
		`, exerciseName)
	
		if err != nil {
			return err
		}
	
		_,err = trnx.Exec(context.Background(), `
			DELETE FROM WORKOUT_TRACKER
			WHERE EXERCISE_NAME = $1		
		`, exerciseName)
		if err != nil {
			return err
		}

		_, err = trnx.Exec(context.Background(), `
			DELETE FROM PLAN 
			WHERE EXERCISE_TRACKER_ID = $1		
		`, exerciseId)

		if err != nil{
			return err
		}

		err = trnx.Commit(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}

func GetExerciseIdFromTrackerInDB(exerciseName string) (int, error){

	var exerciseId int

	err := database.DBConn.QueryRow(context.Background(), `
		SELECT ID FROM EXERCISE_TRACKER
		WHERE EXERCISE_NAME = $1
	`, exerciseName).Scan(&exerciseId)

	if err != nil{
		if err == pgx.ErrNoRows{
			return 0, pgx.ErrNoRows
		}else {
			return exerciseId, err
		}
	}

	return exerciseId, nil
}

// func DeleteUserFromDb(userId int) error {
// 	trnx, err := database.DBConn.Begin(context.Background())
// 	if err != nil {
// 		return err
// 	}

// 	trnx.Exec(context.Background(), `

	
// 	`)
// 	return nil

// 	// delete user from users
// 	// delete user details from workout_tracker, plan, workoutplan
// }





func GetAllUserPlansFromDB(userId int) (pgx.Rows, error) {

	rows, err := database.DBConn.Query(context.Background(), `
		SELECT PLAN_NAME FROM WORKOUTPLANS
		WHERE USER_ID = $1
	`, userId)

	if err != nil {
		return nil, err
	}

	return rows, nil

}

func GetAllUserExercisesByPlanNameFromDB(userId int, planName string) (pgx.Rows, error) {

	rows, err := database.DBConn.Query(context.Background(), `
		SELECT EXERCISE_NAME FROM WORKOUTPLANS
		INNER JOIN PLAN
		ON WORKOUTPLANS.ID = PLAN.PLAN_ID
		INNER JOIN EXERCISE_TRACKER
		ON PLAN.EXERCISE_TRACKER_ID = EXERCISE_TRACKER.ID
		WHERE USER_ID = $1 AND PLAN_NAME = $2;
	`, userId, planName)
	if err != nil {
		return nil, err
	}

	return rows, nil
}