package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Exercise struct {
	UserRole     string `json:"user_role"`
	Id           int    `json:"id"`
	ExerciseName string `json:"exercise_name"`
	Type         string `json:"type"`
	BodyPart     string `json:"body_part"`
}

type AddExerciseToPlanStruct struct {
	UserId int `json:"user_id"`
	ExerciseName string `json:"exercise_name"`
	PlanName string `json:"plan_name"`
}

type AllExercises struct {
	Exercises []Exercise `json:"exercises"`
}

type User struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserSentExercises struct {
	PlanName       string   `json:"plan_name"`
	PlanId         int   `json:"plan_id"`
	UserId         int      `json:"user_id"`
	ExercisesNames []string `json:"exercise_names"`
}

type UserClaims struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type AllUserPlans struct {
	UserId int `json:"user_id"`
	UserPlans map[string][]string `json:"plan_names"`
}

type SetsWeight struct {
	ExerciseName string `json:"exercise_name"`
	SetsAndWeigts [3][2]int `json:"sets_weights"`
}

type Session struct {
	SessionId int `json:"session_id"`
	PlanName string `json:"plan_name"`
	StartTime time.Time `json:"started_at"`
	EndTime time.Time `json:"ended_at"`
	Open bool `json:"open"`
}

type AddRepsWeights struct {
	SetsAndRepsId int `json:"id"`
	UserId int `json:"user_id"`
	PlanName string `json:"plan_name"`
	SessionId int `json:"session_id"`
	ExerciseName string `json:"exercise_name"`
	SetNumber int `json:"set_number"`
	RepCount int `json:"rep_count"`
	Weight int `json:"weight"`
	Comments string `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
}

type EndSession struct {
	SessionId int `json:"session_id"`
	PlanName string `json:"plan_name"`
	StartedAt time.Time `json:"started_at"`
	EndedAt time.Time `json:"ended_at"`
}

type RepsAndWeight struct {
	Reps int `json:"rep_count"`
	Weight int `json:"weight"`
}

type SetRepsWeights struct {
	Set map[int]RepsAndWeight `json:"set"`
}


type RepsWeightsMap map[string]int

type SetsRepsWeight map[string]RepsWeightsMap

type ExerciseSetsRepsWeights map[string]SetRepsWeights

type SetsArray []SetsRepsWeight

type ExerciseNameSet map[string]SetsArray



// {
// 	"exercise_name" : "x",
// 	tracker := 
// }