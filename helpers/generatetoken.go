package helpers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type customClaims struct {
	UserId 	primitive.ObjectID 	`json:"user_id"`
	jwt.StandardClaims
}

func GenerateJwtToken(UserId primitive.ObjectID, ExpirationDate time.Time) (string, error) {
	claims := customClaims {
		UserId: UserId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt : ExpirationDate.Unix(),
			IssuedAt : time.Now().Unix(),
		},
	}

	// creating token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	// parsing token with secret key
	jwtSecret := "aa3cd62b7b442634bc91c8df818ee5633152154c9208f2301ccd6b19f0e8b675"
	tokenString,err := token.SignedString([]byte(jwtSecret));
	if err != nil {
		return "",err
	}

	return tokenString,nil
}