package routes

import (
	"backend/controllers"
	middleware "backend/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	userRoute := r.Group("/user")
	{
		userRoute.GET("/",middleware.JwtMiddleware(),controllers.FetchUser)
		userRoute.GET("/jobs",controllers.GetAllJobListings())
		userRoute.POST("/signin",controllers.HandleSignin)
		userRoute.POST("/signup",controllers.HandleSignup) 
	}
}

