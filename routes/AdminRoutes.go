package routes

import (
	"backend/controllers"
	middleware "backend/middlewares"

	"github.com/gin-gonic/gin"
)

func AdminAuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("admin/login", controllers.Login())
	incomingRoutes.POST("admin/signup", controllers.Register())	
}

func AdminRoutes(incomingRoutes *gin.Engine) {
	r := incomingRoutes.Group("/admin",middleware.Authenticate())
	{
		r.POST("/create-client",controllers.HandleSignin)
		r.GET("/allAdmins", controllers.GetAllAdmins())
		r.GET("/adminInfo", controllers.GetUserInfo())
		r.GET("/all-requests", controllers.GetAllRequests())
		r.GET("/adminsList", controllers.GetAllAdmins())
		r.GET("/adminByID/:id", controllers.GetAdminByID())
		r.GET("/client/all-clients", controllers.GetAllClients())
		r.GET("/client/:id", controllers.GetClientByID())
		r.GET("/job/",controllers.GetAllJobListings());
		r.GET("/job/:id",controllers.GetJobListingByID());


		r.POST("/sendRequest", controllers.SendRequest())
		r.POST("/job/new",middleware.Authenticate(),controllers.CreateJobListing());
		r.POST("/send-document",middleware.Authenticate(),controllers.UploadPdf())


		r.PUT("/modify-request/:id", controllers.ApproveOrRejectRequest())

		r.PATCH("/job/update/:id",middleware.Authenticate(),controllers.UpdateJobListingByID());

		r.DELETE("/job/delete/:id",middleware.Authenticate(),controllers.DeleteJobListingById());

		r.POST("/client/uploadPdf", controllers.UploadPdf())
		r.POST("/client/getPdfByEmail",controllers.GetPdfDetailsByUserEmail())
		r.GET("/client/all-documents",controllers.GetAllPdfDetails())
		r.GET("/admin/adminsList", controllers.GetAllAdmins())


	}
}


