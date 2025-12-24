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

	r.Get("/workout/exercises", handler.GetAllExercises)

	// user routes
	r.Post("/workout/signup", handler.UserSignup)
	r.Post("/workout/login", handler.UserLogin) 
	r.With(middleware.JwtMiddleware).Put("/workout/user/update", handler.UserUpdateDetails)
	r.With(middleware.JwtMiddleware).Delete("/workout/user", handler.DeleteUserByUser)
	r.With(middleware.JwtMiddleware).Get("/workout/allUsers", handler.GetAllUsers) 

	// plan routes
	r.With(middleware.JwtMiddleware).Post("/workout/plan", handler.CreatePlan)
	r.With(middleware.JwtMiddleware).Get("/workout/plan/{planname}", handler.GetUserPlan)
	r.With(middleware.JwtMiddleware).Get("/workout/plan", handler.GetAllUserPlans)
	r.With(middleware.JwtMiddleware).Patch("/workout/plan/{planname}", handler.AddExerciseToPlan)


	// session routes
	r.With(middleware.JwtMiddleware).Post("/workout/plan/{planname}/session", handler.CreateNewSession)
	r.With(middleware.JwtMiddleware).Post("/workout/plan/{planname}/session/end", handler.EndASession)
	r.With(middleware.JwtMiddleware).Post("/workout/plan/{planname}/session/addsetsandreps", handler.AddSetAndReps)
	r.With(middleware.JwtMiddleware).Get("/workout/plan/{planname}/session", handler.GetAllUserSessionsByPlanName)
	r.With(middleware.JwtMiddleware).Get("/workout/plan/session", handler.GetAllUserSessions)
	r.With(middleware.JwtMiddleware).Get("/workout/plan/{planname}/session/stats", handler.GetAllExercisesBySession)



	fmt.Println("server is starting at 5000.....")

	http.ListenAndServe(":"+config.Config.WebPort, r)
}


