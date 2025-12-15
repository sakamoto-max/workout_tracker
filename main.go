package main

import (
	"fmt"
	"net/http"
	"workout_tracker/database"
	"workout_tracker/handler"
	"workout_tracker/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

)

func main() {

	database.InitDB()

	r := chi.NewRouter()	

	r.Use(chiMiddleware.Logger)

	r.Get("/workout/exercises", handler.GetAllExercises)
	r.Post("/workout/signup", handler.UserSignup)
	r.Post("/workout/login", handler.UserLogin)
	r.With(middleware.JwtMiddleware).Post("/workout/plan/create", handler.CreatePlan)
	r.With(middleware.JwtMiddleware).Post("/workout/exercise", handler.InsertANewExercise)
	r.With(middleware.JwtMiddleware).Delete("/workout/exercise", handler.DeleteExercise)
	r.With(middleware.JwtMiddleware).Get("/workout/plan", handler.GetAllUserPlans)
	// User_id_plan_name_exercise_name
	

	fmt.Println("server is starting at 5000.....")

	http.ListenAndServe(":5000", r)


	// list all the exercises	// done       -> workout/exercises
	// user signUp	// done                   -> workout/user/signup
	// user login // done                     -> workout/user/login
	// user create a workoutplan // done      -> workout/user/plan/create
	// admin insert a new exercise // done    -> workout/user/exercise
	// admin delete an exercise               -> workout/user/exercise


	//todo :
	// user update a workout plan             -> workout/user/plan/{planname}
	// admin delete an user                   -> workout/user


 
	// user delete all plans             -> workout/user/plan/delete
	// user delete one plan             -> workout/user/plan/{planname}
	// user get his workoutplan               -> workout/user/plan/{planname}
	// user add sets and reps				  -> workout/user/plan/{planname}
	// user get reports of his workouts       -> workout/user/progress
}


