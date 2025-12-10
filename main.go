package main

import (
	"context"
	"net/http"
	"workout_tracker/database"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)



func main() {

	database.InitDB()
	r := chi.NewRouter()	

	r.Get("/exercises", GetAllExercises)


	http.ListenAndServe(":5000", r)
}

type Exercise struct{
	Id int `json:"id"`
	ExerciseName string `json:"exercise_name"`
	Type string `json:"type"`
	BodyPart string `json:"body_part"`
}

type ExercisesStruct struct{
	Exercises []Exercise `json:"exercises"`
}

func GetAllExercises(w http.ResponseWriter, r *http.Request) {
	

}
