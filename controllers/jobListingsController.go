package controllers

import (
	"backend/database"
	"backend/models"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var jobListingCollection *mongo.Collection = database.OpenCollection(database.Client, "jobListings")

func GetTest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "hello test",
		})
	}
}


// jobs

func CreateJobListing() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var jobListing models.JobListing
		if err := ctx.Bind(&jobListing); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "failed to bind json",
			})
			return
		}

		//generatea unique id for the job listing
		jobListing.ID = primitive.NewObjectID()
		_, err := jobListingCollection.InsertOne(context.Background(), jobListing)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create a Job Listing",
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Job listing created successfully", "jobListing": jobListing,
		})
	}
}

func GetAllJobListings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var jobListings []models.JobListing
		log.Println(jobListings);
		cursor, err := jobListingCollection.Find(context.Background(), bson.M{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch job listings",
			})
			return
		}

		defer cursor.Close(context.Background())

		err = cursor.All(context.Background(), &jobListings)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to decode the job listings",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"jobListing": jobListings,
		})

	}
}

func GetJobListingByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("inside the function")
		// Retrieve the ID from the URL parameters
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
			return
		}

		// Convert the ID string to ObjectID
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job listing ID"})
			return
		}

		// Find the job listing with the given ID
		var jobListing models.JobListing
		err = jobListingCollection.FindOne(context.Background(), bson.M{"_id": objId}).Decode(&jobListing)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Job listing not found"})
			return
		}

		// Return the job listing with status OK
		ctx.JSON(http.StatusOK, gin.H{"jobListing": jobListing})
	}
}

func UpdateJobListingByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Retrieve the ID from the URL parameters
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
			return
		}

		// Convert the ID string to primitive.ObjectID
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		// Parse the JSON request body into a job listing struct
		var updateData models.JobListing
		if err := ctx.BindJSON(&updateData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update data"})
			return
		}

		// Prepare update document
		update := bson.M{}
		if updateData.Role != "" {
			update["role"] = updateData.Role
		}
		if updateData.Location != "" {
			update["location"] = updateData.Location
		}
		if updateData.Link != "" {
			update["link"] = updateData.Link
		}
		if updateData.Company != "" {
			update["company"] = updateData.Company
		}
		if updateData.Description != "" {
			update["desc"] = updateData.Description
		}
		if len(updateData.Requirements) > 0 {
			update["requirements"] = updateData.Requirements
		}
		update["active"] = updateData.Active

		// Perform the update operation
		_, err = jobListingCollection.UpdateOne(
			context.Background(),
			bson.M{"_id": objID},
			bson.M{"$set": update},
		)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job listing"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Job listing updated successfully", "jobListing": updateData})
	}
}

func DeleteJobListingById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Retrieve the ID from the URL parameters
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
			return
		}

		// Convert the ID string to primitive.ObjectID
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		// Delete the job listing with the given ID
		_, err = jobListingCollection.DeleteOne(context.Background(), bson.M{"_id": objId})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the job listing"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Job listing deleted successfully"})
	}
}
