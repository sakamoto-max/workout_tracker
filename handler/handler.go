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
	"workout_tracker/validations"
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

	validationErr, err := validations.UserSignUpValidator(user)
	if err != nil {

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationErr)
		return
	}

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

	validationErr, err := validations.UserLoginValidator(userSentDetails)
	if err != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationErr)
		return
	}
		
	err = service.UserLoginService(userSentDetails)
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
		
	jwtToken, err := auth.GenerateJwtToken(userSentDetails.Email)
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

	validationErr, err := validations.CreatePlanValidator(exercises)
	if err != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationErr)
		return
	}

	exercises.UserId = claimsFromRequest.UserId

	err = service.CreatePlanService(exercises)

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

func AddExerciseToPlan(w http.ResponseWriter, r *http.Request) {

	var exercise models.Exercise

	// exerciseName := r.PathValue("exercisename")
	planName := r.PathValue("planname")

	
	claimsFromRequest, ok := utils.GetClaimsFromRequest(r.Context())

	if !ok {
		response := map[string]string{
			"message": "error fetching claims from req",
		}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	json.NewDecoder(r.Body).Decode(&exercise)

	Errors, err := validations.AddExerciseToPlanValidator(exercise.ExerciseName)
	if err != nil{

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Errors)		
	}



	err = service.AddExerciseToPlanService(claimsFromRequest.UserId, planName, exercise.ExerciseName)
	if err != nil {
		fmt.Printf("error occured : %v", err)
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
}

// session plans

func GetAllExercisesBySession(w http.ResponseWriter, r *http.Request) {

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

	response, err := service.GetAllExercisesBySessionService(claimsFromRequest.UserId, planName)
	if err != nil {
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

	validationErr, err := validations.AddRepsWeightsValidator(addRepsAndWeights)
	if err != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationErr)
		return
	}

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








