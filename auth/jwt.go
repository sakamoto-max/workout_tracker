package auth

import (
	"errors"
	"time"
	"workout_tracker/repository"

	"github.com/golang-jwt/jwt/v5"
)

// what should the jwt token contain?
// user_id

type UserClaims struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var (
	ErrTokenIsInvalid = errors.New("token is invalid")
)
var SECRET_KEY string = "adsliasd2kajdk#2jdoiaj"

func GenerateJwtToken(email string) (string, error) {

	userId, err := repository.GetUserIdFromDB(email)
	if err != nil {
		return "", err
	}

	userRole, err := repository.GetUserRoleFromDB(email)
	if err != nil {
		return "", err
	}

	newClaims := UserClaims{
		UserId: userId,
		Role:   userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "workout_tracker",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)

	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func VerifyJwtToken(token string) (*UserClaims, error) {

	claims := &UserClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return claims, err
	}

	if !parsedToken.Valid {
		return claims, ErrTokenIsInvalid
	}

	return claims, nil
}
