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

	r.LoadHTMLGlob("./templates/*.html")

	r.Static("/assets", "./assets")
	r.Static("/emails", "./emails")

	publicRoutes(r)

	authorizedRoutes(r)

	return r
}

func publicRoutes(r *gin.Engine) {
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	r.POST("/user/register", controllers.Register)
	r.POST("/login", controllers.Login)
}

func authorizedRoutes(r *gin.Engine) {
	authorized := r.Group("/")
	authorized.Use(middlewares.AuthMiddleware())
	{
		authorized.GET("/home", homePage)

		authorized.GET("/docs", redirectToDocs)
		authorized.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		authorized.GET("/diag/health", controllers.HealthCheck)
		authorized.GET("/test/panic", testPanic)

		authorized.GET("/user/email/:email", controllers.GetUserByEmail)
		authorized.GET("/user/all", controllers.GetUserAll)
	}
}

func redirectToDocs(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/docs/index.html")
}

func testPanic(c *gin.Context) {
	panic("Isso é um teste de pânico!")
}

func homePage(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", nil)
}
