package routes

import (
	"github.com/gin-gonic/gin"

	"portofolio-api/controllers"
	"portofolio-api/middlewares"
)

func SkillRoutes(r *gin.Engine) {
	skill := r.Group("/skill")
	{
		skill.POST("", middlewares.AuthMiddleware(), controllers.CreateSkill)
		skill.POST("/", middlewares.AuthMiddleware(), controllers.CreateSkill)
		skill.GET("/:id", middlewares.AuthMiddleware(), controllers.GetSkillByID)
		skill.GET("/user/:user_id", middlewares.AuthMiddleware(), controllers.GetSkillsByUserID)
		skill.PUT("/:id", middlewares.AuthMiddleware(), controllers.UpdateSkill)
		skill.DELETE("/:id", middlewares.AuthMiddleware(), controllers.DeleteSkill)
	}
}
