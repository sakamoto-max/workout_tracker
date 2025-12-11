package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"workout_tracker/database"
	"workout_tracker/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)



func main() {

	database.InitDB()
	r := chi.NewRouter()	
	r.With(middleware.Logger)

	r.Get("/exercises", GetAllExercises)

	fmt.Println("server is starting at 5000.....")

	http.ListenAndServe(":5000", r)
}



func GetAllExercises(w http.ResponseWriter, r *http.Request) {
	response, err := service.GetAllExercisesService()
	if err != nil {
		fmt.Printf("error occured : %v\n", err)
		response := map[string]string{
			"message" : "error occured",
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
	}else {
		
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

}
