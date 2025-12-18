package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"workout_tracker/auth"
	"workout_tracker/customerrors"
	"workout_tracker/models"
	"workout_tracker/service"
	"workout_tracker/utils"
)

// user handlers
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	
	claimsFromRequest, ok := utils.GetClaimsFromRequest(r.Context())

	if !ok {
		fmt.Printf("error occured : error getting token from request")
		response := map[string]string {
			"message" : "error occured",
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response, err := service.GetAllUsersService(claimsFromRequest.Role)
	if err != nil{
		fmt.Printf("error occured : %v", err)
		response := map[string]string {
			"message" : err.Error(),
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
		return 
	} 

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
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
		return
	} 

	jwtToken, err := auth.GenerateJwtToken(userSentDetails.Email)
	if err != nil {
		fmt.Printf("error occured : %v\n", err)
		return
	}

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

func UserUpdateDetails(w http.ResponseWriter, r *http.Request) {

	
	// can change name
	// can change email
	// can change password
	
	var userUpdateDetails models.User
	
	claims, ok := utils.GetClaimsFromRequest(r.Context())
	
	if !ok {
		response := map[string]string{
			"message": "error occured",
		}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	

	json.NewDecoder(r.Body).Decode(&userUpdateDetails)

	userUpdateDetails.Id = claims.UserId
	
	response, err := service.UserUpdateDetailsService(userUpdateDetails)

	if err != nil{
		response := map[string]string{
			"message" : err.Error(),
		}	

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func DeleteUserByUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := utils.GetClaimsFromRequest(r.Context())	
	
	if !ok {
		response := map[string]string{
			"message": "error occured",
		}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}


	err := service.DeleteUserByUserService(claims.UserId)
	if err != nil {
		fmt.Printf("error occred : %v", err)
		response := map[string]string {
			"message" : "error occured",
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return		
	}

	w.WriteHeader(http.StatusNotFound)
	
}


// exercise handlers
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
		return
	} 

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
		return
	}

	newExercise.UserRole = claims.Role

	response, err := service.InsertANewExerciseService(newExercise)
	if err != nil{
		if err == customerrors.ErrOnlyAdminAccess{
			response := map[string]string{
				"message" : customerrors.ErrOnlyAdminAccess.Error(),
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
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	
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
		return	
	}

	var exerciseName models.Exercise

	json.NewDecoder(r.Body).Decode(&exerciseName)

	exerciseName.UserRole = claims.Role

	err := service.DeleteExerciseService(exerciseName.ExerciseName, exerciseName.UserRole)
	if err != nil{
		if err == customerrors.ErrOnlyAdminAccess{

			response := map[string]string{
				"message" : customerrors.ErrOnlyAdminAccess.Error(),
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
		return
	}
	
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
}

// plan handlers

func CreatePlan(w http.ResponseWriter, r *http.Request) {

	claimsFromRequest, ok := utils.GetClaimsFromRequest(r.Context())

	if !ok {
		response := map[string]string{
			"message": "error occured",
		}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return

	}

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
		return
	}
	response := map[string]string{
		"message": "plan created successfully",
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

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
}

func GetUserPlan(w http.ResponseWriter, r *http.Request) {

	claims, ok := utils.GetClaimsFromRequest(r.Context())

	if !ok {
		response := map[string]string{
			"message" : "failed to get claims from request",
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	planName := r.PathValue("planname")

	response, err := service.GetUserPlanService(claims.UserId, planName)
	if err != nil {
		fmt.Printf("error occured : %v", err)
		response := map[string]string{
			"message": "error occured",
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return	
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// session plans

func CreateNewSession(w http.ResponseWriter, r *http.Request) {

	claimsFromRequest, ok := utils.GetClaimsFromRequest(r.Context())

	if !ok {
		response := map[string]string{
			"message" : "failed to get claims from request",
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	

	planName := r.PathValue("planname")

	response, err := service.CreateNewSessionService(claimsFromRequest.UserId, planName)
	if err != nil {
		fmt.Printf("error occured : %v\n", err)
		response := map[string]string{
			"message": err.Error(),
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func EndASession(w http.ResponseWriter, r *http.Request) {

	claimsFromRequest, ok := utils.GetClaimsFromRequest(r.Context())
	
	if !ok {
		response := map[string]string{
			"message" : "failed to get claims from request",
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	planName := r.PathValue("planname")


	response, err := service.EndASessionService(claimsFromRequest.UserId, planName)
	if err != nil{
		fmt.Printf("error occured : %v\n", err)
		response := map[string]string {
			"message" : err.Error(),
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func AddSetAndReps(w http.ResponseWriter, r *http.Request) {

	var addRepsAndWeights models.AddRepsWeights

	claimsFromRequest, ok := utils.GetClaimsFromRequest(r.Context())
	if !ok {
		fmt.Printf("error occured : error getting token from request")
		response := map[string]string {
			"message" : "error occured",
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	json.NewDecoder(r.Body).Decode(&addRepsAndWeights)

	planName := r.PathValue("planname")

	addRepsAndWeights.UserId = claimsFromRequest.UserId
	addRepsAndWeights.PlanName = planName
	


	response, err := service.AddSetAndRepsService(addRepsAndWeights) 

	if err != nil{
		if err == customerrors.ErrSessionIsClosed {
			response := map[string]string {
				"message" : err.Error(),
			}
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
		}else {
			fmt.Printf("error occured : %v\n", err)
			response := map[string]string {
				"message" : err.Error(),
			}
	
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
		}
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetAllUserSessionsByPlanName(w http.ResponseWriter, r *http.Request) {

	claims, ok := utils.GetClaimsFromRequest(r.Context())
	if !ok {
		fmt.Printf("error occured : error getting token from request")
		response := map[string]string {
			"message" : "error occured",
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}	

	planName := r.PathValue("planname")


	response, err := service.GetAllUserSessionsByPlanNameService(claims.UserId, planName) 
	if err != nil{
		response := map[string]string {
			"message" : err.Error(),
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}


func GetAllUserSessions(w http.ResponseWriter, r *http.Request) {

	claims, ok := utils.GetClaimsFromRequest(r.Context())

	if !ok {
		fmt.Printf("error occured : error getting token from request")
		response := map[string]string {
			"message" : "error occured",
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}	

	response, err := service.GetAllUserSessionsService(claims.UserId)
		if err != nil{
		response := map[string]string {
			"message" : err.Error(),
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)



}

// stats handlers

func GetStatsByExerciseName(w http.ResponseWriter, r *http.Request) {

	// claims, ok := utils.GetClaimsFromRequest(r.Context())

	// if !ok {
	// 	fmt.Printf("error occured : error getting token from request")
	// 	response := map[string]string {
	// 		"message" : "error occured",
	// 	}

	// 	w.Header().Set("Content-type", "application/json")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(response)
	// 	return
	// }	

	// exerciseName = r.PathValue("exercisename")


	
}








