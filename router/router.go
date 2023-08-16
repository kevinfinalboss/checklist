package router

import (
	"net/http"

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
	r.Static("/assets", "./assets")

	r.Use(middlewares.RateLimitMiddleware())

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	r.POST("/login", controllers.Login)

	authorized := r.Group("/")
	authorized.Use(middlewares.AuthMiddleware())
	{
		authorized.GET("/diag/health", controllers.HealthCheck)
		authorized.GET("/test/panic", func(c *gin.Context) {
			panic("Isso é um teste de pânico!")
		})

		authorized.GET("/home", func(c *gin.Context) {
			c.HTML(http.StatusOK, "home.html", nil)
		})
	}

	return r
}
