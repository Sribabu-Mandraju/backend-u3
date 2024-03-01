package controllers

import (
	"backend/database"
	"backend/helpers"
	"backend/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FetchUser(c *gin.Context){
	var user models.User
	user_id := c.GetString("user_id");
	objectId,err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "objectId not found"})
		return
	}

	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	// get Db
	collection := database.DBinstance().Database("u3technologies").Collection("users");

	if err := collection.FindOne(ctx,bson.M{"_id":objectId}).Decode(&user) ; err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Credentials"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK,gin.H{
		"user":user,
	})
}

func HandleSignin(c *gin.Context){
	var user models.User
	var input models.InputUser

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	// get Db
	collection := database.DBinstance().Database("u3technologies").Collection("users");

	// find user
	err := collection.FindOne(ctx,bson.M{"email":input.Email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Credentials"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// checking password
	if err := helpers.CheckPassword(user.Password,input.Password);err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Credentials"})
		return
	}

	// generating jwt
	token, err := helpers.GenerateJwtToken(user.ID, time.Now().Add(24*time.Hour))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK,gin.H{
		"message":"Logged in successfully",
		"user":user,
		"token":token,
	})
}

func HandleSignup(c *gin.Context){
	var user models.User	

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := database.DBinstance().Database("u3technologies").Collection("users");

	// Check if email is already used
	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"User Already Exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Internal Server Error"})
		return
	}

	// Hash password and set timestamps
	hashedPassword, err := helpers.GenerateHash(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to hash password"})
		return
	}
	user.Password = hashedPassword
	user.SetUserTimeStamps()

	// Create new user in DB
	user.ID = primitive.NewObjectID()
	res, err := collection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to create new user"})
		return
	}

	userID := res.InsertedID.(primitive.ObjectID)

	// Create JWT token
	token, err := helpers.GenerateJwtToken(userID, time.Now().Add(24*time.Hour))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to Create Jwt"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusCreated, gin.H{
		"message": "New User created successfully",
		"user":    user,
		"token":   token,
	})
}