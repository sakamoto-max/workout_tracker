package service

import (
	"errors"

	"strings"
	"workout_tracker/models"
	"workout_tracker/repository"
	"workout_tracker/utils"

	"github.com/jackc/pgx/v5"
)

var (
	ErrOnlyAdminAccess = errors.New("only admin can modify this")
)

func GetAllExercisesService() (models.AllExercises, error) {
	var allExercises models.AllExercises
	rowsFromDb, err := repository.GetAllExercisesFromDB()
	
	if err != nil{
		return allExercises, err
	}
	
	defer rowsFromDb.Close()

	var id int
	var exerciseName string
	var exerciseType string
	var bodyPart string

	for rowsFromDb.Next() {
		
		rowsFromDb.Scan(&id, &exerciseName, &exerciseType, &bodyPart)
		
		exercise := models.Exercise{
			Id: id,
			ExerciseName: exerciseName,
			Type: exerciseType,
			BodyPart: bodyPart,
		}

		allExercises.Exercises = append(allExercises.Exercises, exercise)
	}

	return allExercises, nil
}


func UserSignupService(user models.User) (models.User, error ){
	var userDetailsFromDb models.User

	hashedPassword, err := utils.HashThePassword(user.Password)
	if err != nil {
		return userDetailsFromDb, err
	}

	user.Password = hashedPassword

	if user.Role == "" {
		user.Role = "user"
	}

	err = repository.CreateUserInDB(user)
	if err != nil {
		return userDetailsFromDb, err
	}

	userDetailsFromDb, err = repository.GetUserFromDb(user.Email)
	if err != nil{
		return userDetailsFromDb, err
	}

	return userDetailsFromDb, nil
}

func UserLoginService(userSentDetails models.User) (error) {

	// compare the passwords
	err := utils.PasswordMatcher(userSentDetails.Email, userSentDetails.Password)
	
	if err != nil{
		return err
	}else {
		return nil
	}	
}

func CreatePlanService(userSentExercises models.UserSentExercises) error {

	err := repository.EeAaO(userSentExercises.UserId, userSentExercises.PlanName, userSentExercises.ExercisesNames)
	if err != nil {
		return err
	}

	return nil
}

func InsertANewExerciseService(newExercise models.Exercise) (models.Exercise,error) {

	var response models.Exercise
	if strings.ToUpper(newExercise.UserRole) == "ADMIN" {
		responseFromdb, err := repository.InsertANewExerciseInDB(newExercise.ExerciseName, newExercise.Type, newExercise.BodyPart)
		if err != nil{
			return responseFromdb, err
		}

		response = responseFromdb
		response.UserRole = "Admin"
	}else {
		return response, ErrOnlyAdminAccess
	}

	return response, nil
}

func DeleteExerciseService(exerciseName string, role string) (error) {

	if strings.ToUpper(role) == "ADMIN" {

		var exerciseId int
		exerciseId, err := repository.GetExerciseIdFromTrackerInDB(exerciseName)
		if err != nil{
			if err == pgx.ErrNoRows{
				exerciseId = 0
			}else {
				return err
			}
		}
	
		err = repository.DeleteExerciseFromDb(exerciseName, exerciseId)
		if err != nil {
			return err
		}
	}else {
		return ErrOnlyAdminAccess
	}


	
	return nil
}

// func DeleteUserService(userId int) {



// }

func CreatePlanService2(userId int, ) {

}

func GetAllUserPlansService(userId int) (models.AllUserPlans, error) {

	var allPlans models.AllUserPlans
	plansExercises := make(map[string][]string)

	rows, err := repository.GetAllUserPlansFromDB(userId)
	if err != nil{
		return allPlans, err
	}

	defer rows.Close()

	var planName string

	for rows.Next() {
		err := rows.Scan(&planName)
		if err != nil{
			return allPlans, err
		}

		rows, err := repository.GetAllUserExercisesByPlanNameFromDB(userId, planName)
		if err != nil {
			return allPlans, err
		}
		defer rows.Close()

		var exerciseName string
		var allExerciseNames []string

		for rows.Next() {
			err := rows.Scan(&exerciseName)
			if err != nil{
				return allPlans, err
			}

			allExerciseNames = append(allExerciseNames, exerciseName)
		}

		
		// allPlans.UserPlans[planName] = AllExerciseNames
		plansExercises[planName] = allExerciseNames
	}

	allPlans.UserId = userId
	allPlans.UserPlans = plansExercises
	return allPlans, nil
}