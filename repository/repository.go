package repository

import (
	"context"
	"workout_tracker/customerrors"
	"workout_tracker/database"
	"workout_tracker/models"

	"github.com/jackc/pgx/v5"
)

// user funtions
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

func GetAllUsersFromDB() (pgx.Rows, error) {
	rows, err := database.DBConn.Query(context.Background(), `
		SELECT ID, NAME, EMAIL, ROLE, CREATED_AT, UPDATED_AT
		FROM USERS	
		WHERE ROLE = $1
	`, "user")
	if err != nil {
		return nil, err
	}

	return rows, nil
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

func GetUserFromDBbyId(userId int) (models.User, error) {

	var oD models.User

	err := database.DBConn.QueryRow(context.Background(), `
		SELECT * FROM USERS
		WHERE ID = $1	
	`, userId).Scan(&oD.Id, &oD.Name, &oD.Email, &oD.Password, &oD.Role, &oD.CreatedAt, &oD.UpdatedAt)
	if err != nil{
		return oD, err
	}

	return oD, nil
}

func UpdateUserDetailsInDB(newDetails models.User) (models.User, error) {

	var a models.User

	err := database.DBConn.QueryRow(context.Background(), `
		UPDATE USERS
		SET NAME = $1, EMAIL = $2, HASHED_PASSWORD = $3, UPDATED_AT = NOW()
		WHERE ID = $4
		RETURNING ID, NAME, EMAIL, ROLE, CREATED_AT, UPDATED_AT
	`, newDetails.Name, newDetails.Email, newDetails.Password, newDetails.Id).Scan(&a.Id, &a.Name, 
		&a.Email, &a.Role, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		return a, err
	}

	return a, nil
}


func DeleteUserByUserInDB(userId int, setRepIds []int, sessionIds []int, planExerciseIds []int, planIds []int) (error) {
	// _, err := database.DBConn.Exec(context.Background(), `
	// 	DELETE FROM USERS
	// 	WHERE ID = $1	
	// `, userId)

	// if err != nil {
	// 	return err
	// }

	// return nil

	trnx, err := database.DBConn.Begin(context.Background())
	if err != nil {
		return err
	}

	for _, v := range(setRepIds) {
		_, err := trnx.Exec(context.Background(), `
			DELETE FROM SETREPS
			WHERE ID = $1
		`, v)

		if err != nil {
			return err
		}
	}

	for _, v := range(sessionIds) {
		_, err := trnx.Exec(context.Background(), `
			DELETE FROM SESSION
			WHERE ID = $1
		`, v)

		if err != nil {
			return err
		}
	}

	for _, v := range(planExerciseIds) {
		_, err := trnx.Exec(context.Background(), `
			DELETE FROM PLAN_EXERCISES
			WHERE ID = $1
		`, v)

		if err != nil {
			return err
		}
	}

	for _, v := range(planIds) {
		_, err := trnx.Exec(context.Background(), `
			DELETE FROM PLANS
			WHERE ID = $1
		`, v)

		if err != nil {
			return err
		}
	}

	_, err = trnx.Exec(context.Background(), `
		DELETE FROM USERS
		WHERE ID = $1	
	`, userId)
	if err != nil {
		return err
	}


	err = trnx.Commit(context.Background())
		if err != nil {
		return err
	}

	return nil
}


// session function

func CreateANewSessionInDB(planId int, planName string) (models.Session, error) {

	var session models.Session

	trnx, err := database.DBConn.Begin(context.Background())
	if err != nil {
		return session, err
	}
	err = trnx.QueryRow(context.Background(), `
		INSERT INTO SESSION(PLAN_ID, PLAN_NAME, STARTED_AT, OPEN)
		VALUES($1, $2, NOW(), $3)
		RETURNING ID, STARTED_AT, OPEN
	`,planId, planName, true).Scan(&session.SessionId, &session.StartTime, &session.Open)
	if err != nil{
		return session, err
	}

	session.PlanName = planName

	err = trnx.Commit(context.Background())
	if err != nil {
		return session, err
	}

	return session, nil
}

func EndASessionInDB(planId int, planName string) (models.EndSession, error) {

	var endSession models.EndSession

	err := database.DBConn.QueryRow(context.Background(), `
		UPDATE SESSION
		SET OPEN = FALSE, ended_at = NOW()
		WHERE plan_id = $1 AND PLAN_NAME = $2
		RETURNING ID, PLAN_NAME, STARTED_AT, ENDED_AT
	`, planId, planName).Scan(&endSession.SessionId, &endSession.PlanName, &endSession.StartedAt, &endSession.EndedAt)
	if err != nil {
		return endSession, err
	}

	return endSession, nil
}

func AddSetAndRepsInDB(reps models.AddRepsWeights) (models.AddRepsWeights, error){

	var r models.AddRepsWeights

	err := database.DBConn.QueryRow(context.Background(), `
		INSERT INTO SETREPS(EXERCISE_NAME, SESSION_ID, REPS, WEIGHT, COMMENTS, CREATED_AT, SET_NUMBER)
		VALUES($1, $2, $3, $4, $5, NOW(), $6)
		RETURNING ID, EXERCISE_NAME, SESSION_ID, REPS, WEIGHT, COMMENTS, CREATED_AT, SET_NUMBER
	`, reps.ExerciseName, reps.SessionId, reps.RepCount, reps.Weight, reps.Comments, reps.SetNumber).Scan(
			&r.SetsAndRepsId, &r.ExerciseName, &r.SessionId, &r.RepCount, &r.Weight, &r.Comments, &r.CreatedAt, &r.SetNumber)

	if err != nil {
		return r, err
	}

	return r, nil
}

func GetSessionIdFromDB(userId int, planName string) (int, bool, error) {

	var sessionId int
	var open bool

	err := database.DBConn.QueryRow(context.Background(), `
		SELECT SESSION.ID, SESSION.OPEN FROM SESSION
		INNER JOIN PLANS
		ON PLANS.ID = SESSION.PLAN_ID
		WHERE USER_ID = $1 AND SESSION.PLAN_NAME = $2 AND OPEN = TRUE
	`, userId, planName).Scan(&sessionId, &open)

	if err != nil {
		return 0, false, err
	}

	return sessionId, open, nil
}

func GetSessionIdFromDBTwo(userId int, planName string) (int, bool, error) {

	var sessionId int
	var open bool

	err := database.DBConn.QueryRow(context.Background(), `
		SELECT SESSION.ID, SESSION.OPEN FROM SESSION
		INNER JOIN PLANS
		ON PLANS.ID = SESSION.PLAN_ID
		WHERE USER_ID = $1 AND SESSION.PLAN_NAME = $2
	`, userId, planName).Scan(&sessionId, &open)

	if err != nil {
		return 0, false, err
	}

	return sessionId, open, nil
}

func CheckIfSessionIsOpen(userId int, planName string) (error) {
	// check if user has a same plan name

	var open bool
	err := database.DBConn.QueryRow(context.Background(), `
		SELECT OPEN FROM SESSION
		INNER JOIN PLANS
		ON SESSION.PLAN_ID = PLANS.ID
		WHERE USER_ID = $1 AND OPEN IS TRUE
	`, userId).Scan(&open)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	}

	if open {
		return customerrors.ErrDuplicateSession		
	}else {
		return nil
	}
	// no rows -> create session
	// true -> dont create session
	// false -> create session
}

func GetAllUserSessionsByPlanNameFromDb(userId int, planName string) (pgx.Rows, error) {

	rows, err := database.DBConn.Query(context.Background(), `
		SELECT SESSION.id, SESSION.plan_name, started_at, ended_at, open FROM SESSION
		INNER JOIN PLANS
		ON SESSION.PLAN_ID = PLANS.ID
		WHERE USER_ID = $1 AND SESSION.PLAN_NAME = $2 AND OPEN IS FALSE
	`, userId, planName)

	if err != nil{
		return nil, err
	}

	return rows, err
}

func GetAllUserSessions(userId int) (pgx.Rows, error) {
	
	rows, err := database.DBConn.Query(context.Background(), `
		SELECT SESSION.id, SESSION.plan_name, started_at, ended_at, open FROM SESSION
		INNER JOIN PLANS
		ON SESSION.PLAN_ID = PLANS.ID
		WHERE USER_ID = $1
	`, userId)
	
	if err != nil{
		return nil, err
	}
	
	return rows, err
}

func GetAllSessionIdsOfUser(userId int) ([]int ,error) {
	var ids []int

	rows, err := database.DBConn.Query(context.Background(), `
		SELECT DISTINCT SESSION.ID FROM SESSION
		INNER JOIN PLANS
		ON SESSION.PLAN_ID = PLANS.ID
		WHERE USER_ID = $1
	`, userId)
	if err != nil {
		return ids, err
	}

	defer rows.Close()

	var id int

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil

}

func GetAllIdsFromSetReps(userId int) ([]int, error) {
		var ids []int

	rows, err := database.DBConn.Query(context.Background(), `
		SELECT SETREPS.ID FROM SESSION
		INNER JOIN PLANS
		ON SESSION.PLAN_ID = PLANS.ID
		INNER JOIN SETREPS
		ON SETREPS.SESSION_ID = SESSION.ID
		WHERE USER_ID = $1
	`, userId)

	if err != nil {
		return ids, err
	}

	defer rows.Close()

	var id int

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil

}

func DeleteFromSessionBySessionId(ids []int) error {
	trnx, err := database.DBConn.Begin(context.Background())
	if err != nil {
		return err
	}

	for _, v := range(ids) {
		_, err := trnx.Exec(context.Background(), `
			DELETE FROM SESSION
			WHERE ID = $1
		`, v)

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

func DeleteFromSetRepsById(ids []int) (error) {

	trnx, err := database.DBConn.Begin(context.Background())
	if err != nil {
		return err
	}

	for _, v := range(ids) {
		_, err := trnx.Exec(context.Background(), `
			DELETE FROM SETREPS
			WHERE ID = $1
		`, v)

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

func GetLastSetNumber(userId int, exerciseName string, sessionId int) (int, error)  {
	var setNumber int

	err := database.DBConn.QueryRow(context.Background(), `
		SELECT SET_NUMBER FROM SETREPS
		INNER JOIN SESSION
		ON SETREPS.SESSION_ID = SESSION.ID
		INNER JOIN PLANS
		ON SESSION.PLAN_ID = PLANS.ID
		WHERE USER_ID = $1 AND EXERCISE_NAME = $2 AND SESSION_ID = $3 AND SESSION.OPEN IS TRUE
		ORDER BY SET_NUMBER DESC
		LIMIT 1
	`,userId, exerciseName, sessionId).Scan(&setNumber)

	if err != nil {
		return setNumber, err
	}

	return setNumber, nil
}

// get all the exercises performed in a session

func GetAllExercisesBySession(userId int, planName string, sessionId int) ([]string , error) {

	var exerciseNames []string

	rows, err := database.DBConn.Query(context.Background(), `
		SELECT DISTINCT EXERCISE_NAME FROM SETREPS
		INNER JOIN SESSION
		ON SETREPS.SESSION_ID = SESSION.ID
		INNER JOIN PLANS
		ON SESSION.PLAN_ID = PLANS.ID
		WHERE SESSION.PLAN_NAME = $1 AND USER_ID = $2 AND SESSION_ID = $3
	`, planName, userId, sessionId)

	if err != nil {
		return exerciseNames, err
	}

	defer rows.Close()

	var exerciseName string

	for rows.Next() {
		err := rows.Scan(&exerciseName)
		if err != nil {
			return exerciseNames, err
		}

		exerciseNames = append(exerciseNames, exerciseName)
	}

	return exerciseNames, nil
}

// get no of sets performed for each exercise

func GetNoOfSetsForAExercise(userId int, planName string, sessionId int, exerciseName string) (int, error) {
	var noOfSets int

	err := database.DBConn.QueryRow(context.Background(), `
		SELECT COUNT(SET_NUMBER) FROM SETREPS
		INNER JOIN SESSION
		ON SETREPS.SESSION_ID = SESSION.ID
		INNER JOIN PLANS
		ON SESSION.PLAN_ID = PLANS.ID
		WHERE SESSION.PLAN_NAME = $1 AND USER_ID = $2 AND SESSION_ID = $3 AND EXERCISE_NAME = $4
	`, planName, userId, sessionId, exerciseName).Scan(&noOfSets)

	if err != nil {
		return noOfSets, err
	}

	return noOfSets, nil
}

// userId, planName, sessionId , exerciseName, setNumber
// get reps and weight for each set

func GetRepsAndWeightsForASet(userId int , planName string, sessionId int, exerciseName string, setNumber int) (int, int, error) {

	var reps int
	var weight int

	err := database.DBConn.QueryRow(context.Background(), `
		SELECT REPS, WEIGHT FROM SETREPS
		INNER JOIN SESSION
		ON SETREPS.SESSION_ID = SESSION.ID
		INNER JOIN PLANS	
		ON SESSION.PLAN_ID = PLANS.ID
		WHERE SESSION.PLAN_NAME = $1 AND USER_ID = $2 AND SESSION_ID = $3 AND EXERCISE_NAME = $4 AND SET_NUMBER = $5
	`, planName, userId, sessionId, exerciseName, setNumber).Scan(&reps, &weight)

	if err != nil {
		return reps, weight, err
	}

	return  reps, weight, nil
}






// plan functions

func CreateAPlanInDB(userId int, planName string) (int, error) {
	var planId int

	err := database.DBConn.QueryRow(context.Background(), `
		INSERT INTO PLANS(USER_ID, PLAN_NAME, CREATED_AT, UPDATED_AT)
		VALUES($1, $2, NOW(), NOW())
		RETURNING ID
	`, userId, planName).Scan(&planId)

	if err != nil {
		return planId, err
	}
	return planId, nil
}

func InsertExercisesIntoPlan(planId int, exerciseNames []string) (error) {

	trnx, err := database.DBConn.Begin(context.Background())
	if err != nil {
		return err
	}

	for _, v := range(exerciseNames) {
	
		_, err = trnx.Exec(context.Background(), `
			INSERT INTO PLAN_EXERCISES(PLAN_ID, EXERCISE_NAMES)
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

func GetPlanIdFromDB(userId int, planName string) (int, error) {

	var planId int

	err := database.DBConn.QueryRow(context.Background(), `
		SELECT ID FROM PLANS
		WHERE USER_ID = $1 AND PLAN_NAME = $2
	`, userId, planName).Scan(&planId)

	if err != nil {
		return 0, err
	}

	return planId, nil
}

func GetAllUserExercisesByPlanNameFromDB(userId int, planName string) (pgx.Rows, error) {

	rows, err := database.DBConn.Query(context.Background(), `
		SELECT EXERCISE_NAMES FROM PLAN_EXERCISES
		INNER JOIN PLANS 
		ON PLAN_EXERCISES.PLAN_ID = PLANS.ID
		WHERE USER_ID = $1 AND PLAN_NAME = $2
	`, userId, planName)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func GetAllUserPlansFromDB(userId int) (pgx.Rows, error) {

	rows, err := database.DBConn.Query(context.Background(), `
		SELECT DISTINCT PLAN_NAME FROM PLANS
		WHERE USER_ID = $1
	`, userId)

	if err != nil {
		return nil, err
	}

	return rows, nil

}

func GetAllUserPlanIds(userId int) ([]int, error) {
	var ids []int

	rows, err := database.DBConn.Query(context.Background(), `
		SELECT ID FROM PLANS
		WHERE USER_ID = $1
	`, userId)
	if err != nil {
		return ids, err
	}

	defer rows.Close()

	var id int

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func GetAllUserExercisesIds(userId int) ([]int, error) {
	var ids []int

	rows, err := database.DBConn.Query(context.Background(), `
		SELECT PLAN_EXERCISES.ID FROM PLAN_EXERCISES
		INNER JOIN PLANS
		ON PLAN_EXERCISES.PLAN_ID = PLANS.ID
		WHERE USER_ID = $1
	`, userId)

	if err != nil{
		return ids, err
	}

	defer rows.Close()

	var id int

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}


	return ids, nil
}

func DeleteFromPlanExercisesById(ids []int) error {
		trnx, err := database.DBConn.Begin(context.Background())
	if err != nil {
		return err
	}

	for _, v := range(ids) {
		_, err := trnx.Exec(context.Background(), `
			DELETE FROM PLAN_EXERCISES
			WHERE ID = $1
		`, v)

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

func DeleteFromPlansById(ids []int) error {
		trnx, err := database.DBConn.Begin(context.Background())
	if err != nil {
		return err
	}

	for _, v := range(ids) {
		_, err := trnx.Exec(context.Background(), `
			DELETE FROM PLANS
			WHERE ID = $1
		`, v)

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

// exercise functions

func GetAllExercisesFromDB() (pgx.Rows, error) {

	rows, err := database.DBConn.Query(context.Background(), `
		SELECT * FROM EXERCISES
	`)
	if err != nil {
		return rows, err
	}

	return rows, nil
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
		`, exerciseName)

		if err != nil {
			return err
		}
	} else {

		trnx, err := database.DBConn.Begin(context.Background())
		if err != nil {
			return err
		}

		_, err = trnx.Exec(context.Background(), `
			DELETE FROM EXERCISES 
			WHERE EXERCISE_NAME = $1
		`, exerciseName)

		if err != nil {
			return err
		}

		_, err = trnx.Exec(context.Background(), `
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

		if err != nil {
			return err
		}

		err = trnx.Commit(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}


func GetExerciseIdFromTrackerInDB(exerciseName string) (int, error) {

	var exerciseId int

	err := database.DBConn.QueryRow(context.Background(), `
		SELECT ID FROM EXERCISE_TRACKER
		WHERE EXERCISE_NAME = $1
	`, exerciseName).Scan(&exerciseId)

	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, pgx.ErrNoRows
		} else {
			return exerciseId, err
		}
	}

	return exerciseId, nil
}


// stat functions

func GetExerciseByPlanName(userId int, planName string, sessionId int) {
	
}

func GetAllExercisesBySessionFromDB(userId int, planName string, sessionId int) (pgx.Rows, error) {

	rows, err := database.DBConn.Query(context.Background(), `
		SELECT SETREPS.EXERCISE_NAME, SETREPS.SET_NUMBER, SETREPS.REPS, SETREPS.WEIGHT FROM SETREPS
		INNER JOIN SESSION
		ON SETREPS.SESSION_ID = SESSION.ID
		INNER JOIN PLANS
		ON SESSION.PLAN_ID = PLANS.ID
		WHERE SESSION.PLAN_NAME = $1 AND USER_ID = $2 AND SESSION_ID = $3
	`, planName, userId, sessionId)

	if err != nil{
		return rows, err
	}

	return rows, err
}


