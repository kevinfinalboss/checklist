package middlewares

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	r := gin.Default()

	r.Use(ErrorHandler())

	r.GET("/test-route", func(c *gin.Context) {
		c.String(200, "OK")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test-route", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
