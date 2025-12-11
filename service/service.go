package service

import (
	"workout_tracker/models"
	"workout_tracker/repository"
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