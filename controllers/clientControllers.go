package controllers

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var pdfUploadsCollection *mongo.Collection = database.OpenCollection(database.Client, "pdfUploads")
var validateClient = validator.New()

func HashPasswordClient(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPasswordClient(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := err == nil
	msg := ""
	if !check {
		msg = fmt.Sprintf("email or password not matched")
	}
	return check, msg
}

func SendRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.Request_to_admin

		// Parse the incoming JSON request body into the Requests struct
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Get the current time and format it as desired
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		request.Sended_At = &currentTime

		// Insert the request data into the MongoDB collection
		result, err := requestCollection.InsertOne(context.Background(), request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		// Respond with success message and inserted data
		c.JSON(http.StatusOK, gin.H{
			"msg":    "successfully request sent",
			"data":   request,
			"status": result,
		})
	}
}

func UploadPdf() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("UploadPdf handler function called")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Parse form
		fmt.Println("Parsing form...")
		err := c.Request.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to parse form",
			})
			fmt.Println("Error parsing form:", err.Error())
			return
		}

		// Get the PDF file from the request
		fmt.Println("Getting PDF file from request...")
		file, _, err := c.Request.FormFile("pdf_file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "PDF file not found in request",
			})
			fmt.Println("PDF file not found in request")
			return
		}
		defer file.Close()

		// Read the file content
		fmt.Println("Reading PDF file content...")
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read PDF file",
			})
			fmt.Println("Failed to read PDF file:", err.Error())
			return
		}

		// Get other fields from the request
		title := c.Request.FormValue("title")
		userEmail := c.Request.FormValue("user_email")

		id := primitive.NewObjectID()

		pdfUpload := models.PDFUploads{
			ID:        id,
			Title:     title,
			UserEmail: userEmail,
			PDFFile:   fileBytes,
		}

		fmt.Println("Validating PDFUploads struct...")
		if err := validateClient.Struct(pdfUpload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Validation error",
			})
			fmt.Println("Validation error:", err.Error())
			return
		}

		fmt.Println("Inserting PDFUploads data into MongoDB collection...")
		result, err := pdfUploadsCollection.InsertOne(ctx, pdfUpload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to insert data into MongoDB collection",
			})
			fmt.Println("Error inserting data into MongoDB collection:", err.Error())
			return
		}

		fmt.Println("PDF uploaded successfully")
		c.JSON(http.StatusOK, gin.H{
			"msg":    "PDF uploaded successfully",
			"data":   pdfUpload,
			"status": result,
		})
	}
}


func GetPdfDetailsByUserEmail() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Parse the JSON request body
        var requestBody map[string]string
        if err := c.BindJSON(&requestBody); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
            return
        }

        // Retrieve the user email from the request body
        userEmail, ok := requestBody["useremail"]
        if !ok || userEmail == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "User email is required"})
            return
        }

        // Retrieve the PDF details from the database based on the user email
        pdfDetails, err := GetPdfDetailsByUserEmailFromDatabase(userEmail)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve PDF details"})
            return
        }

        // Set response headers
        c.Header("Content-Type", "application/json")

        // Send the PDF details as the response
        c.JSON(http.StatusOK, pdfDetails)
    }
}

func GetPdfDetailsByUserEmailFromDatabase(userEmail string) ([]bson.M, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := pdfUploadsCollection.Find(ctx, bson.M{"useremail": userEmail})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var pdfDetails []bson.M
    err = cursor.All(ctx, &pdfDetails)
    if err != nil {
        return nil, err
    }

    return pdfDetails, nil
}




func GetAllPdfDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		pdfDetails, err := GetAllPdfDetailsFromDatabase()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve PDF details"})
			return
		}

		// Set response headers
		c.Header("Content-Type", "application/json")

		// Send the PDF details as the response
		c.JSON(http.StatusOK, pdfDetails)
	}
}

func GetAllPdfDetailsFromDatabase() ([]models.PDFUploads, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	cursor, err := pdfUploadsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pdfDetails []models.PDFUploads
	err = cursor.All(ctx, &pdfDetails)
	if err != nil {
		return nil, err
	}

	return pdfDetails,nil
}
func FetchFileById() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        requestIDParam := c.Param("id")
        requestID, err := primitive.ObjectIDFromHex(requestIDParam)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
            return
        }

        var request bson.M
        err = pdfUploadsCollection.FindOne(ctx, bson.M{"_id": requestID}).Decode(&request)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
            return
        }

        c.JSON(http.StatusOK, request)
    }
}
