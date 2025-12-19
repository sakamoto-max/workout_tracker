package validations

import (
	"errors"
	"workout_tracker/models"
)

type ErrorStruct struct{
	Path string `json:"path"`
	Message string `json:"message"`
}

var(
	ErrUserNameReq = ErrorStruct{Path: "name", Message: "required"}
	ErrUserEmailReq = ErrorStruct{Path: "email", Message: "required"}
	ErrUserPassReq = ErrorStruct{Path: "password", Message: "required"}
	ErrUserRoleOptional = ErrorStruct{Path: "role", Message: "optional"}
	ErrExerciseNameReq = ErrorStruct{Path: "exercise_name", Message: "required"}
	ErrBodyPartReq = ErrorStruct{Path: "body_part", Message: "required"}
	ErrTypeReq = ErrorStruct{Path: "type", Message: "required"}
	ErrPlanNameReq = ErrorStruct{Path: "plan_name", Message: "required"}
	ErrExerciseNamesReq = ErrorStruct{Path: "exercise_names", Message: "required"}
	ErrRepCountReq = ErrorStruct{Path: "rep_count", Message: "required"}
	ErrWeightOptional = ErrorStruct{Path: "rep_count", Message: "optional"}
	ErrCommentsOptional = ErrorStruct{Path: "rep_count", Message: "optional"}
	
	ErrSomeError = errors.New("chumma error")
)



func UserSignUpValidator(userSentDetails models.User) ([]ErrorStruct, error) {

	var errors []ErrorStruct


	if userSentDetails.Name == "" && userSentDetails.Email == "" && userSentDetails.Password == "" && userSentDetails.Role == "" {
		errors = append(errors, ErrUserNameReq, ErrUserEmailReq, ErrUserPassReq, ErrUserRoleOptional)
		return errors, ErrSomeError
	}

	if userSentDetails.Name == "" || userSentDetails.Email == "" || userSentDetails.Password == ""{

		if userSentDetails.Name == "" {
			a := ErrorStruct{Path: "name", Message: "required"}
	
			errors = append(errors, a)
		}
		if userSentDetails.Email == "" {
			a := ErrorStruct{Path: "email", Message: "required"}
			errors = append(errors, a)
		}
		if userSentDetails.Password == "" {
			a := ErrorStruct{Path: "password", Message: "required"}
			errors = append(errors, a)
		}

		return errors, ErrSomeError
	}

	return errors, nil
}

func UserLoginValidator(userSentDetails models.User) ([]ErrorStruct, error) {
	var errors []ErrorStruct

	if userSentDetails.Email == "" || userSentDetails.Password == "" {
		if userSentDetails.Email == "" {
			errors = append(errors, ErrUserEmailReq)
		}

		if userSentDetails.Password == "" {
			errors = append(errors, ErrUserPassReq)
		}

		return errors, ErrSomeError
	}

	return errors, nil
}

func InsertNewExerciseValidator(userSentDetails models.Exercise) ([]ErrorStruct, error){

	var errors []ErrorStruct

	if userSentDetails.ExerciseName == "" || userSentDetails.BodyPart == "" || userSentDetails.Type == "" {
		if userSentDetails.ExerciseName == "" {
			errors = append(errors, ErrExerciseNameReq)
		}

		if userSentDetails.BodyPart == "" {
			errors = append(errors, ErrBodyPartReq)
		}

		if userSentDetails.Type == "" {
			errors = append(errors, ErrTypeReq)
		}

		return errors, ErrSomeError
	}

	return errors, nil

}

// DeleteExercise, CreatePlan, AddSetAndReps

func DeleteExerciseValidator() () {


}

// {
// 	"plan_name" :
// 	"exercise_name" : []
// }

func CreatePlanValidator(userSentDetails models.UserSentExercises) ([]ErrorStruct, error) {
	var errors []ErrorStruct

	if userSentDetails.PlanName == "" || len(userSentDetails.ExercisesNames) == 0 {

		if userSentDetails.PlanName == "" {
		   errors = append(errors, ErrPlanNameReq)
	   }
   
	   if len(userSentDetails.ExercisesNames) == 0 {
		   errors = append(errors, ErrExerciseNamesReq)
	   }

	   return errors, ErrSomeError
	}

	return errors, nil
}

// {
//     "exercise_name" : "incline_lateral_raises",
//     "rep_count" : 8,
//     "weight" : 5
// }

func AddRepsWeightsValidator(userSentDetails models.AddRepsWeights) ([]ErrorStruct, error) {
	var errors []ErrorStruct

	if userSentDetails.ExerciseName == "" && userSentDetails.RepCount == 0 && userSentDetails.Comments == "" && userSentDetails.Weight == 0{
		
		errors = append(errors, ErrExerciseNameReq, ErrRepCountReq, ErrCommentsOptional, ErrWeightOptional)
		return errors, ErrSomeError

	}
	if userSentDetails.ExerciseName == "" || userSentDetails.RepCount == 0 {
		
		if userSentDetails.ExerciseName == "" {
			errors = append(errors, ErrExerciseNameReq)
		}
		
		if userSentDetails.RepCount == 0 {
			errors = append(errors, ErrRepCountReq)
		}

		return errors, ErrSomeError

	}

	return errors, nil
}