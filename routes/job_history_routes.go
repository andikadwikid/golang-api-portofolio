package routes

import (
	"github.com/gin-gonic/gin"

	"portofolio-api/controllers"
	"portofolio-api/middlewares"
)

func JobHistoryRoutes(r *gin.Engine) {
	jobHistory := r.Group("/job-history")
	{
		jobHistory.POST("/", middlewares.AuthMiddleware(), controllers.CreateJobHistory)
		jobHistory.GET("/:id", middlewares.AuthMiddleware(), controllers.GetJobHistoryByID)
		jobHistory.GET("/user/:user_id", middlewares.AuthMiddleware(), controllers.GetJobHistoryByUserID)
		jobHistory.PUT("/:id", middlewares.AuthMiddleware(), controllers.UpdateJobHistory)
		jobHistory.DELETE("/:id", middlewares.AuthMiddleware(), controllers.DeleteJobHistory)
	}
}
