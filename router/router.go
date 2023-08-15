package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kevinfinalboss/checklist-apps/api/middlewares"
	_ "github.com/kevinfinalboss/checklist-apps/docs"
	"github.com/kevinfinalboss/checklist-apps/pkg/controllers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.ErrorHandler())
	r.LoadHTMLGlob("./templates/*")

	r.Use(middlewares.RateLimitMiddleware())

	url := ginSwagger.URL("http://localhost:80/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.GET("/diag/health", controllers.HealthCheck)
	r.GET("/test/panic", func(c *gin.Context) {
		panic("Isso é um teste de pânico!")
	})

	return r
}
