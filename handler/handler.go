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
				"message": "login successful",
			}

			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}

	}
}

func CreatePlan(w http.ResponseWriter, r *http.Request) {

	claimsFromRequest, ok := utils.GetClaimsFromRequest(r.Context())

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

func DeleteExercise(w http.ResponseWriter, r *http.Request) {

	claims, ok := utils.GetClaimsFromRequest(r.Context())
	if !ok {
		response := map[string]string{
			"message" : "cannot get claims",
		}
	
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)	
	}

	var exerciseName models.Exercise

	json.NewDecoder(r.Body).Decode(&exerciseName)

	exerciseName.UserRole = claims.Role

	err := service.DeleteExerciseService(exerciseName.ExerciseName, exerciseName.UserRole)
	if err != nil{
		if err == service.ErrOnlyAdminAccess{

			response := map[string]string{
				"message" : service.ErrOnlyAdminAccess.Error(),
			}
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
		}else {
			fmt.Printf("error occured : %v", err)
			response := map[string]string{
				"message" : "error occured",
			}
	
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
		}
	}else {

		response := map[string]string{
			"message" : "deleted successfully",
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
	}
}

// func DeleteUser(w http.ResponseWriter, r *http.Request) {
	
// 	userIdString := r.PathValue("id")

// 	userId, err := strconv.Atoi(userIdString)

// 	if err != nil{
// 		response := map[string]string{
// 			"message" : "error converting string to int",
// 		}

// 		w.Header().Set("Content-type", "application/json")
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(response)
// 		fmt.Println(userId)
// 		return
// 	}
// }

func GetAllUserPlans(w http.ResponseWriter, r *http.Request) {

	claimsFormRequest, ok := utils.GetClaimsFromRequest(r.Context())

	if !ok {
		response := map[string]string{
			"message" : "couldn't fetch claims from the request",
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	allPlans, err := service.GetAllUserPlansService(claimsFormRequest.UserId)
	
	if err != nil {
		fmt.Printf("error occured : %v", err)
		response := map[string]string{
			"message" : "error occured",
		}
	
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allPlans)
	return
}