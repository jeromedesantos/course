package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	var (
		bytes []byte
		err   error
	)

	// convert plain -> hash pass (in byte)
	bytes, err = bcrypt.GenerateFromPassword([]byte(password), 14) // 14 = salt
	return string(bytes), err
}

func CheckPasswordHash(password, hashedPassword string) bool {
	var (
		err	error
	)

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}