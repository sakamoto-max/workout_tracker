package models

type Exercise struct{
	Id int `json:"id"`
	ExerciseName string `json:"exercise_name"`
	Type string `json:"type"`
	BodyPart string `json:"body_part"`
}

type AllExercises struct{
	Exercises []Exercise `json:"exercises"`
}