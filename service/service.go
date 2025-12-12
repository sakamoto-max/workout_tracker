package service

import (
	"errors"
	"strings"
	"workout_tracker/models"
	"workout_tracker/repository"
	"workout_tracker/utils"
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