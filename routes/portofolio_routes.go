package routes

import (
	"github.com/gin-gonic/gin"

	"portofolio-api/controllers"
	"portofolio-api/middlewares"
)

func PortofolioRoutes(r *gin.Engine) {
	portofolio := r.Group("/portofolio")
	{
		portofolio.POST("/", middlewares.AuthMiddleware(), controllers.CreatePortofolio)
		portofolio.GET("/", middlewares.AuthMiddleware(), controllers.GetMyPortofolios)
		portofolio.GET("/user/:user_id", middlewares.AuthMiddleware(), controllers.GetPortofolioByUserID)
		portofolio.PUT("/:id", middlewares.AuthMiddleware(), controllers.UpdatePortofolio)
		portofolio.PATCH("/:portofolio_id/status", middlewares.AuthMiddleware(), controllers.UpdatePortofolioStatus)
		portofolio.DELETE("/:portofolio_id", middlewares.AuthMiddleware(), controllers.DeletePortofolio)

		// Public/All
		portofolio.GET("/all", controllers.GetAllPortofolios)
	}
}
