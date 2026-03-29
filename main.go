package main

import (
	"portofolio-api/database"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	r := gin.Default()

	r.Run(":8080")
}
