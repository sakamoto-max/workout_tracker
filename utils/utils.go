package utils

import "golang.org/x/crypto/bcrypt"

func HashThePassword(password string) (string, error) {
	hashedPassInbytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", err
	}

	return string(hashedPassInbytes), nil
}