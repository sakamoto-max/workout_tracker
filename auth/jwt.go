package auth

import (
	"time"
	"workout_tracker/config"
	"workout_tracker/customerrors"
	"workout_tracker/models"
	"workout_tracker/repository"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwtToken(email string) (string, error) {

	userId, err := repository.GetUserIdFromDB(email)
	if err != nil {
		return "", err
	}

	userRole, err := repository.GetUserRoleFromDB(email)
	if err != nil {
		return "", err
	}

	newClaims := models.UserClaims{
		UserId: userId,
		Role:   userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "workout_tracker",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)

	signedToken, err := token.SignedString([]byte(config.Config.SecretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func VerifyJwtToken(token string) (*models.UserClaims, error) {

	claims := &models.UserClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		return []byte(config.Config.SecretKey), nil
	})

	if err != nil {
		return claims, err
	}

	if !parsedToken.Valid {
		return claims, customerrors.ErrTokenIsInvalid
	}

	return claims, nil
}
