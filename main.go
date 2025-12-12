package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"workout_tracker/database"
	"workout_tracker/handler"
	"workout_tracker/middleware"
	"workout_tracker/models"
	"workout_tracker/service"
	"workout_tracker/utils"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {

	database.InitDB()
	r := chi.NewRouter()	

	r.Use(chiMiddleware.Logger)

	r.Get("/workout/exercises", handler.GetAllExercises)
	r.Post("/workout/user/signup", handler.UserSignup)
	r.Post("/workout/user/login", handler.UserLogin)
	r.With(middleware.JwtMiddleware).Post("/workout/user/plan/create", handler.CreatePlan)
	// handler to add a new exercise
	r.With(middleware.JwtMiddleware).Post("/workout/user/exercise", InsertANewExercise)

	fmt.Println("server is starting at 5000.....")

	http.ListenAndServe(":5000", r)
}

// {
// 	"exercise_Nmae" : x,
// 	"type" : y,
// 	"body_part" : z
// }
// type 

func InsertANewExercise(w http.ResponseWriter, r *http.Request) {

	var newExercise models.Exercise

	json.NewDecoder(r.Body).Decode(&newExercise)

	claims, ok :=utils.GetClaimsFromRequest(r.Context())
	if !ok {
		response := map[string]string{
			"message" : "error in getting claims from Request",
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
	}

	newExercise.UserRole = claims.Role

	response, err := service.InsertANewExerciseService(newExercise)
	if err != nil{
		if err == service.ErrOnlyAdminAccess{
			response := map[string]string{
				"message" : service.ErrOnlyAdminAccess.Error(),
			}
			
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
		}else {
			fmt.Printf("error occured : %v\n", err)
			response := map[string]string{
				"message" : "error occured",
			}
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
		}
	}else {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
