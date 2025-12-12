package utils

import (
	"context"
	"errors"
	"workout_tracker/auth"
	"workout_tracker/middleware"
	"workout_tracker/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordMisMatch = errors.New("passwords doesn't match")
)

func HashThePassword(password string) (string, error) {
	hashedPassInbytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", err
	}

	return string(hashedPassInbytes), nil
}

func PasswordMatcher(UserSentEmail string ,UserSentPassword string) error {

	hashedPassword, err := repository.GetHashedPassFromDB(UserSentEmail)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(UserSentPassword))
	if err != nil{
		return ErrPasswordMisMatch
	}

	return nil

}

func GetClaimsFromRequest(c context.Context)  (*auth.UserClaims, bool) {
	claims, ok := c.Value(middleware.ClaimsKey).(*auth.UserClaims)
	return claims, ok
}