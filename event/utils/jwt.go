package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "supersecret"

func GenerateToken(email string, userId int64) (string, error) {
	var (
		token *jwt.Token
	)

	// generate new token 
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email" : email,
		"userId" : userId,
		"exp" : time.Now().Add(time.Hour * 2).Unix(),
	}) // --> signature and data

	// signed token --> 2 string result (token + error)
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (int64, error) {

	var (
		parsedToken *jwt.Token
		tokenIsValid bool
		claims jwt.MapClaims
		userId int64
		ok bool
		err error
	)

	// jwt.Parse --> parse token, ambil key, cek signature
	parsedToken, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// check itu beneran HS256 apa bukan
		_, ok = token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return []byte(secretKey), nil
	} )

	if err != nil {
		return 0, errors.New("Could not parse token.")
	}

	tokenIsValid = parsedToken.Valid

	if !tokenIsValid {
		return 0, errors.New("Invalid token!")
	}

	// ambil value dari jwt
	claims, ok = parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Invalid token claims.")
	}

	// email = claims["email"].(string)
	userId = int64(claims["userId"].(float64))

	return userId, nil
}