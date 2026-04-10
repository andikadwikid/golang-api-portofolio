package routes

import (
	"github.com/gin-gonic/gin"

	"portofolio-api/controllers"
	"portofolio-api/middlewares"
)

func CertificateRoutes(r *gin.Engine) {
	certificate := r.Group("/certificate")
	{
		certificate.POST("/", middlewares.AuthMiddleware(), controllers.CreateCertificate)
		certificate.GET("/:id", middlewares.AuthMiddleware(), controllers.GetCertificateByID)
		certificate.GET("/user/:user_id", middlewares.AuthMiddleware(), controllers.GetCertificatesByUserID)
		certificate.PUT("/:id", middlewares.AuthMiddleware(), controllers.UpdateCertificate)
		certificate.DELETE("/:id", middlewares.AuthMiddleware(), controllers.DeleteCertificate)
	}
}
