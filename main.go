package main

import (
	"fmt"
	"net/http"
	"workout_tracker/config"
	"workout_tracker/database"
	"workout_tracker/handler"
	"workout_tracker/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)


func main() {

	config.InitializeConfig()
	database.InitDB()
	r := chi.NewRouter()	

	r.Use(chiMiddleware.Logger)
	// users - C R U D // DONE
	// plans -  U D
	// sessions - U D
	// admin exercises - C R U D
	// get stats



	r.Get("/workout/exercises", handler.GetAllExercises)
	r.Post("/workout/signup", handler.UserSignup)
	r.Post("/workout/login", handler.UserLogin) 
	r.With(middleware.JwtMiddleware).Put("/workout/user/update", handler.UserUpdateDetails)

	r.With(middleware.JwtMiddleware).Delete("/workout/user", handler.DeleteUserByUser)

	r.With(middleware.JwtMiddleware).Get("/workout/allUsers", handler.GetAllUsers) 
	r.With(middleware.JwtMiddleware).Post("/workout/plan", handler.CreatePlan)
	r.With(middleware.JwtMiddleware).Get("/workout/plan/{planname}", handler.GetUserPlan)
	r.With(middleware.JwtMiddleware).Get("/workout/plan", handler.GetAllUserPlans)

	r.With(middleware.JwtMiddleware).Post("/workout/plan/{planname}/session", handler.CreateNewSession)
	r.With(middleware.JwtMiddleware).Post("/workout/plan/{planname}/session/end", handler.EndASession)
	r.With(middleware.JwtMiddleware).Post("/workout/plan/{planname}/session/addsetsandreps", handler.AddSetAndReps)

	r.With(middleware.JwtMiddleware).Get("/workout/plan/{planname}/session", handler.GetAllUserSessionsByPlanName)
	r.With(middleware.JwtMiddleware).Get("/workout/plan/session", handler.GetAllUserSessions)
	// r.With(middleware.JwtMiddleware).Get("/workout/plan/{planname}/session/{exercisename}", handler.GetStatsByExerciseName)
	r.With(middleware.JwtMiddleware).Get("/workout/plan/{planname}/session/stats", handler.GetAllExercisesBySession)

	// r.Get("/workout/plan/{exercisename}/getstats", handler.GetStatsByExerciseName)

	// get stats by exercise name
	// exercise_name -> date reps weight




	fmt.Println("server is starting at 5000.....")

	http.ListenAndServe(":"+config.Config.WebPort, r)


	// list all the exercises	// done       -> workout/exercises
	// user signUp	// done                   -> workout/signup
	// user login // done                     -> workout/login
	// user create a workoutplan // done      -> workout/plan/create
	// admin insert a new exercise // done    -> workout/exercise
	// admin delete an exercise               -> workout/exercise


	// user update a workout plan             -> workout/plan/{planname}
	// admin delete an user                   -> workout/user


 
	// user delete all plans             -> workout/plan/delete
	// user delete one plan             -> workout/plan/{planname}
	// user get his workoutplan               -> workout/plan/{planname}
	// user add sets and reps				  -> workout/plan/{planname}

	// user get reports of his workouts       -> workout/progress
}



