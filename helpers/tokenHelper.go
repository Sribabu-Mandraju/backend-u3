package helpers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"backend/database"
)

type SignedDetails struct {
	Email    string
	Name     string
	Contact  string
	Password string
	Company  string
	User_id  string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "admin")
var SECRET_KEY = "HELLO_WORLD"

func GenerateAllTokens(
	email string,
	name string,
	company string,
	user_id string,
	contact string,
) (signedToken string, signedRefreshToken string) {
	claims := &SignedDetails{
		Email: email,
		Name: name,
		Company: company,
		User_id: user_id,
		Contact: contact,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	var err error  // Declare the 'err' variable here

	// Use '=' instead of ':=' to assign to the existing 'err' variable
	signedToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	signedRefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	return signedToken, signedRefreshToken
}


func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updateObj := bson.D{
		{"$set", bson.D{
			{"token", signedToken},
			{"refresh_token", signedRefreshToken},
			{"updated_at", time.Now()},
		}},
	}

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.Update().SetUpsert(upsert)

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		updateObj,
		opt,
	)

	if err != nil {
		log.Panic(err)
		return
	}
}

func ValidateToken(signedToken string) (claims *SignedDetails , msg string){
	token ,err:= jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token)(interface{} ,error){
			return []byte("HELLO_WORLD"),nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	}
	claims , ok  := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token in invalid")
		msg = err.Error()
		return 
	}

	if claims.ExpiresAt < time.Now().Local().Unix(){
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return 
	}
	return  claims ,msg
}



