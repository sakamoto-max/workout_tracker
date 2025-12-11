package models

import "time"

type Exercise struct {
	Id           int    `json:"id"`
	ExerciseName string `json:"exercise_name"`
	Type         string `json:"type"`
	BodyPart     string `json:"body_part"`
}

type AllExercises struct {
	Exercises []Exercise `json:"exercises"`
}

type User struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
