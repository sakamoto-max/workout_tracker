package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"workout_tracker/auth"
	"workout_tracker/models"
	"workout_tracker/service"
	"workout_tracker/utils"
)

func GetAllExercises(w http.ResponseWriter, r *http.Request) {
	response, err := service.GetAllExercisesService()
	if err != nil {
		fmt.Printf("error occured : %v\n", err)
		response := map[string]string{
			"message": "error occured",
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
	} else {

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

}
func UserSignup(w http.ResponseWriter, r *http.Request) {
	var user models.User

	json.NewDecoder(r.Body).Decode(&user)

	response, err := service.UserSignupService(user)
	if err != nil {
		fmt.Printf("error occured : %v", err)
		response := map[string]string{
			"message": "error occured",
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}

}
func UserLogin(w http.ResponseWriter, r *http.Request) {

	var userSentDetails models.User

	json.NewDecoder(r.Body).Decode(&userSentDetails)

	err := service.UserLoginService(userSentDetails)
	if err != nil {
		fmt.Printf("error occured : %v\n", err)

		response := map[string]string{
			"message": "error occured",
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
	} else {

		jwtToken, err := auth.GenerateJwtToken(userSentDetails.Email)
		if err != nil {
			fmt.Printf("error occured : %v\n", err)
		} else {

			myCookie := http.Cookie{
				Name:     "jwtToken",
				Value:    jwtToken,
				Secure:   true,
				HttpOnly: true,
				Expires:  time.Now().Add(time.Minute * 10),
			}

			http.SetCookie(w, &myCookie)

			response := map[string]string{
				"message": "user login successful",
			}

			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}

	}
}

func CreatePlan(w http.ResponseWriter, r *http.Request) {

	claimsFromRequest, ok := utils.GetClaimsFromRequest(r.Context())
	fmt.Println(claimsFromRequest)

	if !ok {
		response := map[string]string{
			"message": "error occured",
		}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)

	} else {

		var exercises models.UserSentExercises

		json.NewDecoder(r.Body).Decode(&exercises)

		exercises.UserId = claimsFromRequest.UserId

		err := service.CreatePlanService(exercises)

		if err != nil {
			fmt.Printf("error occured : %v\n", err)

			response := map[string]string{
				"message": "error occured",
			}
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
		} else {
			response := map[string]string{
				"message": "exercises uploaded successfully",
			}

			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(response)

		}

	}

}
