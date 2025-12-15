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
	PlanId         string   `json:"plan_id"`
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