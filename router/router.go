package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kevinfinalboss/checklist-apps/pkg/controllers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/diag/health", controllers.HealthCheck)

	return r
}
