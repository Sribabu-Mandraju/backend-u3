package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"backend/database"
	helper "backend/helpers"
	"backend/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "admin")
var clientCollection2 *mongo.Collection = database.OpenCollection(database.Client, "users")
var requestCollection *mongo.Collection = database.OpenCollection(database.Client, "requests")


var validate = validator.New()
var secretKey = []byte("HELLO_WORLD")






func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}


func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := err == nil
	msg := ""
	if !check {
		msg = fmt.Sprintf("email or password not matched")
	}
	return check, msg
}


func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.Admin

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		validationErr := validate.Struct(user)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return
		}

		countByEmail, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred",
			})
			return
		}

		countByContact, err := userCollection.CountDocuments(ctx, bson.M{"contact": user.Contact})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred",
			})
			return
		}

		// Check if either email or contact already exists
		if countByEmail > 0 && countByContact > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "this email or contact already exists",
			})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.ID = primitive.NewObjectID()
		userID := user.ID.Hex()
		user.User_id = &userID
		token, refreshToken := helper.GenerateAllTokens(
			*user.Email,
			*user.Name,
			*user.Company,
			*user.User_id,
			*user.Contact,
		)
		user.Token = &token
		user.Refresh_token = &refreshToken

		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("user item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": msg,
			})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, gin.H{
			"message": "User registered successfully",
			"user":    user,
		})

		return 
	}
}


func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.Admin
		var foundUser models.Admin

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "email or password not matched",
			})
			return
		}

		passwordIsValid, _ := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "hello world",
			})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "HELLO",
			})
			return
		}

		token, refreshToken := helper.GenerateAllTokens(
			*foundUser.Email,
			*foundUser.Name,
			*foundUser.Company,
			*foundUser.User_id,
			*foundUser.Contact,
		)
		helper.UpdateAllTokens(token, refreshToken, *foundUser.User_id)
		err = userCollection.FindOne(ctx, bson.M{"user_id": *foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}

func UpdateUserDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.Admin

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Validate the user input
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return
		}

		// Find the user by ID
		userID := c.Param("id") // Assuming the user ID is passed as a URL parameter
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		filter := bson.M{"_id": objID}
		update := bson.M{
			"$set": bson.M{
				"name":    user.Name,
				"email":   user.Email,
				"contact": user.Contact,
				// Add more fields to update as needed
			},
		}

		// Update the user details in the database
		result, err := userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error occurred while updating user details",
			})
			return
		}

		// Check if any user document was updated
		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		defer cancel()

		// Fetch the updated user details from the database
		var updatedUser models.Admin
		err = userCollection.FindOne(ctx, filter).Decode(&updatedUser)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error occurred while fetching updated user details",
			})
			return
		}

		// Return the updated user details in the response
		c.JSON(http.StatusOK, gin.H{
			"message": "User details updated successfully",
			"user":    updatedUser,
		})
	}
}

func GetAllAdmins() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.TODO()
	
		cursor, err := userCollection.Find(ctx, bson.M{"user_type":"ADMIN"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)
	
		var users []bson.M
		err = cursor.All(ctx, &users)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		c.JSON(http.StatusOK, users)
	}
}

func GetAdminByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.TODO()
	
		userIDParam := c.Param("id")
		userID, err := primitive.ObjectIDFromHex(userIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
	
		var user bson.M
		err = userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
	
		c.JSON(http.StatusOK, user)
	}
}




func GetUserInfo() gin.HandlerFunc{
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		fmt.Println("token string",tokenString)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not provided"})
			return
		}
	
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		fmt.Println("token",token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
	
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		fmt.Println("claims",claims)
	
		userID, ok := claims["User_id"].(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information"})
			return
		}
		fmt.Println("id",userID)
	
		var user models.Admin
		ctx := context.TODO()
		err = userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information"})
			return
		}
	
		c.JSON(http.StatusOK, user)
	}
	
}

func GetAllRequests() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Fetch all request elements from the database
        var requests []models.Request_to_admin
        cursor, err := requestCollection.Find(context.Background(), bson.M{}, options.Find().SetSort(bson.D{{"sended_at", -1}})) // Sorting by sended_at in ascending order
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "failed to fetch requests",
            })
            return
        }
        defer cursor.Close(context.Background())

        if err := cursor.All(context.Background(), &requests); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "failed to decode requests",
            })
            return
        }

        // Return the list of request elements to the client
        c.JSON(http.StatusOK, requests)
    }
}



func ApproveOrRejectRequest() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get the request ID from the URL path parameters
        requestID := c.Param("id")

        // Parse requestID to ObjectId (assuming MongoDB ObjectId)
        objectID, err := primitive.ObjectIDFromHex(requestID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "invalid request ID",
            })
            return
        }

        // Get the action (approve or reject) from the request body
        var action struct {
            Action string `json:"action" binding:"required"`
        }
        if err := c.BindJSON(&action); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "invalid request body",
            })
            return
        }

        // Update the status of the request in the database
        update := bson.M{"status_review": action.Action}
        _, err = requestCollection.UpdateOne(
            context.Background(),
            bson.M{"_id": objectID},
            bson.M{"$set": update},
        )
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "failed to update request",
            })
            return
        }

        // Respond with a success message
        c.JSON(http.StatusOK, gin.H{
            "msg": fmt.Sprintf("request %s successfully", action.Action),
        })
    }
}


func GetAllClients() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.TODO()
	
		cursor, err := clientCollection2.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)
	
		var users []bson.M
		err = cursor.All(ctx, &users)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		c.JSON(http.StatusOK, users)
	}
}

func GetClientByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.TODO()
	
		userIDParam := c.Param("id")
		userID, err := primitive.ObjectIDFromHex(userIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
	
		var user bson.M
		err = clientCollection2.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
	
		c.JSON(http.StatusOK, user)
	}
}

func SendEmailHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientEmail := c.PostForm("email")
		subject := c.PostForm("subject")
	
		err := sendEmail(clientEmail, subject)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	
		c.JSON(200, gin.H{"message": "Email sent successfully"})
	}
}

func sendEmail(to, subject string) error {
	// Sender email configuration
	from := "sribabumandraju@gmail.com"
	password := "63037sribabu" 
	smtpHost := "smtp.gmail.com"   
	smtpPort := "587"                  

	// Message body
	body := "This is the body of your email."

	// Message
	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	fmt.Println("Email sent successfully to", to)
	return nil
}