package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kevinfinalboss/checklist-apps/api/middlewares"
	"github.com/kevinfinalboss/checklist-apps/pkg/controllers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.ErrorHandler())
	r.LoadHTMLGlob("./templates/*")

	r.GET("/diag/health", controllers.HealthCheck)
	r.GET("/test/panic", func(c *gin.Context) {
		panic("Isso é um teste de pânico!")
	})

	return r
}
