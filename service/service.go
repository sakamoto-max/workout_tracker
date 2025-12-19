package service

import (
	"fmt"
	"strings"
	"workout_tracker/customerrors"
	"workout_tracker/models"
	"workout_tracker/repository"
	"workout_tracker/utils"

	"github.com/jackc/pgx/v5"
)

// user services
func GetAllUsersService(role string) ([]models.User ,error) {

	var users []models.User

	if role == "user" {
		return users, customerrors.ErrOnlyAdminAccess
	}

	rows, err := repository.GetAllUsersFromDB()
	if err != nil {
		return users, err
	}

	defer rows.Close()
    // ID, NAME, EMAIL, ROLE, CREATED_AT, UPDATED_AT
	for rows.Next() {
		var userDetails models.User
		err := rows.Scan(&userDetails.Id, &userDetails.Name, &userDetails.Email, &userDetails.Role ,&userDetails.CreatedAt, &userDetails.UpdatedAt)
		if err != nil{
			return users, err
		}

		userDetails.Password = "confidential"

		users = append(users, userDetails)
	}

	return users, nil

	
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
	}else {
		user.Role = strings.ToLower(user.Role)
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

	err := utils.PasswordMatcher(userSentDetails.Email, userSentDetails.Password)
	
	if err != nil{
		return err
	}else {
		return nil
	}	
}

func UserUpdateDetailsService(updationDetails models.User) (models.User, error) {

	var newDetails models.User
	// get user old details
	// match/
	// fill missing values
	// update
	
	
	originalDetails, err := repository.GetUserFromDBbyId(updationDetails.Id)
	if err != nil {
		return newDetails, err
	}	

	if updationDetails.Password != "" {

		updationDetails.Password, err = utils.HashThePassword(updationDetails.Password)
		if err != nil {
			return newDetails, err
		}

		originalDetails.Password = updationDetails.Password
	}


	if updationDetails.Name != "" {
		originalDetails.Name = updationDetails.Name

	}

	if updationDetails.Email != "" {
		originalDetails.Email = updationDetails.Email
	}

	newDetails, err = repository.UpdateUserDetailsInDB(originalDetails)
	if err != nil{
		return newDetails, err
	}

	newDetails.Password = "confidential"

	return newDetails, nil
}

func DeleteUserByUserService(userId int) (error) {

	sessionIds, err := repository.GetAllSessionIdsOfUser(userId)
	if err != nil{
		return err
	}


	planIds, err := repository.GetAllUserPlanIds(userId) 
	if err != nil {
		return err
	}

	exerciseIds, err := repository.GetAllUserExercisesIds(userId)
	if err != nil {
		return err
	}

	setRepIds, err := repository.GetAllIdsFromSetReps(userId) 
	if err != nil {
		return err
	}

	err = repository.DeleteUserByUserInDB(userId, setRepIds, sessionIds, exerciseIds, planIds)
	if err != nil {
		return err
	}
	return nil
}


// plan services

func CreatePlanService(userSentExercises models.UserSentExercises) error {

	planId, err := repository.CreateAPlanInDB(userSentExercises.UserId, userSentExercises.PlanName)
	if err != nil{
		return err
	}

	err = repository.InsertExercisesIntoPlan(planId, userSentExercises.ExercisesNames)
	if err != nil {
		return err
	}

	return nil
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

func GetUserPlanService(userId int, planName string) (models.UserSentExercises ,error) {
	var planNameExercises models.UserSentExercises
	var exerciseNames []string

	planId, err := repository.GetPlanIdFromDB(userId, planName)
	if err != nil {
		return planNameExercises, err
	}



	rows, err := repository.GetAllUserExercisesByPlanNameFromDB(userId, planName)
	if err != nil {
		return planNameExercises, err
	}

	var exerciseName string

	for rows.Next() {
		err := rows.Scan(&exerciseName)
		if err != nil {
			return planNameExercises, err
		}

		exerciseNames = append(exerciseNames, exerciseName)
	}

	planNameExercises.PlanId = planId
	planNameExercises.PlanName = planName
	planNameExercises.ExercisesNames = exerciseNames
	planNameExercises.UserId = userId

	return planNameExercises, nil
}


// session services

func CreateNewSessionService(userId int, planName string) (models.Session, error) {

	var session models.Session
	planId ,err := repository.GetPlanIdFromDB(userId, planName)
	if err != nil {
		return session, err
	}

	err = repository.CheckIfSessionIsOpen(userId, planName)

	if err != nil {
		return session, err
	}

	session, err = repository.CreateANewSessionInDB(planId, planName)
	if err != nil {
		return session, err
	}

	return session, nil




	// if err != nil {
	// 	if err == pgx.ErrNoRows{

	// 	}else {
	// 		return session, err
	// 	}
	// }

	// if open {
	// 	return session, customerrors.ErrDuplicateSession
	// }

	// session, err = repository.CreateANewSessionInDB(planId, planName)
	// if err != nil {
	// 	return session, err
	// }

	// return session, nil

}

func AddSetAndRepsService(addRepsAndWeights models.AddRepsWeights) (models.AddRepsWeights, error) {

	var response models.AddRepsWeights
	sessionId, open, err := repository.GetSessionIdFromDB(addRepsAndWeights.UserId, addRepsAndWeights.PlanName)
	if err != nil {
		if err == pgx.ErrNoRows{
			return response, customerrors.ErrSessionIsClosed
		}
		return response, err
	}

	addRepsAndWeights.SessionId = sessionId

	if open {

		// var lastSetNumber int
		lastSetNumber, err := repository.GetLastSetNumber(addRepsAndWeights.UserId, addRepsAndWeights.ExerciseName, addRepsAndWeights.SessionId)
		fmt.Printf("last set number : %v\n", lastSetNumber)
		
		if err != nil {
			if err == pgx.ErrNoRows{
				lastSetNumber = 0
			}else {
				return response, err
			}
		}
		addRepsAndWeights.SessionId = sessionId
		newSetNumber := lastSetNumber + 1
		fmt.Printf("new set number : %v\n", newSetNumber)
		addRepsAndWeights.SetNumber = newSetNumber
	
		response, err = repository.AddSetAndRepsInDB(addRepsAndWeights)
		if err != nil {
			return response, err
		}
		response.PlanName = addRepsAndWeights.PlanName
		response.UserId = addRepsAndWeights.UserId
	}else {
		return response, customerrors.ErrSessionIsClosed
	}

	return response, nil

}

func EndASessionService(userId int, planName string) (models.EndSession, error) {
	// get planId

	var response models.EndSession
	planId, err := repository.GetPlanIdFromDB(userId, planName)
	if err != nil {
		return response, err
	}

	response, err = repository.EndASessionInDB(planId, planName)
	if err != nil{
		return response, err
	}

	return response, nil
}

func GetAllUserSessionsByPlanNameService(userId int, planName string) ([]models.Session, error) {

	var allSessions []models.Session

	rows, err := repository.GetAllUserSessionsByPlanNameFromDb(userId, planName)
	if err != nil{
		return allSessions, err

	}

	var a models.Session

	for rows.Next() {
		err := rows.Scan(&a.SessionId, &a.PlanName, &a.StartTime ,&a.EndTime, &a.Open)
		if err != nil{
			return allSessions, err
		}

		allSessions = append(allSessions, a)
	}

	return allSessions, nil
	
}

func GetAllUserSessionsService(userId int) ([]models.Session, error) {
	
	var allSessions []models.Session
	
	rows, err := repository.GetAllUserSessions(userId)
	if err != nil{
		return allSessions, err
	
	}
	
	var a models.Session
	
	for rows.Next() {
		err := rows.Scan(&a.SessionId, &a.PlanName, &a.StartTime ,&a.EndTime, &a.Open)
		if err != nil{
			return allSessions, err
		}
	
		allSessions = append(allSessions, a)
	}
	
	return allSessions, nil
}



// exercise services
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
		return response, customerrors.ErrOnlyAdminAccess
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
		return customerrors.ErrOnlyAdminAccess
	}


	
	return nil
}


// stats services

func GetStatsByExerciseNameService(userId int, exerciseName string) {

	


}



// type ExerciseNameSet struct {
// 	ExerciseName map[string]SetRepsWeights `json:"exercise_name"`
// }

// {
// 	"leg_curls" : [
// 		"set_1" : [
// 			"reps" : 10,
// 			"weight" : 20
// 		]
// 	]
// }

type New struct {

}





func GetAllExercisesBySessionService(userId int, planName string) ([]models.ExerciseNameSet, error) {
	var AllExercises []models.ExerciseNameSet

	noOfSetsPerExercise := make(map[string]int)


	sessionId, _,  err := repository.GetSessionIdFromDBTwo(userId, planName)

	if err != nil {
		return AllExercises, err
	}

	exerciseNames, err := repository.GetAllExercisesBySession(userId, planName, sessionId)
	if err != nil {
		return AllExercises, err
	}

	for _, v := range(exerciseNames) {
		noOfSets, err := repository.GetNoOfSetsForAExercise(userId, planName, sessionId, v)
		if err != nil {
			return AllExercises, err
		}

		noOfSetsPerExercise[v] = noOfSets
	}

	for exerciseName, noOfSets := range(noOfSetsPerExercise) {

		exerciseNameSet := make(models.ExerciseNameSet)

		var setsArray models.SetsArray

		for j := range(noOfSets) {

			repsWeightsMap := make(map[string]int)
			set := make(models.SetsRepsWeight)

			setName := fmt.Sprintf("set_%v", j+1)
			reps, weight, err := repository.GetRepsAndWeightsForASet(userId, planName, sessionId, exerciseName, j+1)
			if err != nil {
				return AllExercises, err
			}

			repsWeightsMap["reps"] = reps
			repsWeightsMap["weight"] = weight

			set[setName] = repsWeightsMap

			setsArray = append(setsArray, set)
		}

		exerciseNameSet[exerciseName] = setsArray

		AllExercises = append(AllExercises, exerciseNameSet)
	}

	return AllExercises, nil
}






