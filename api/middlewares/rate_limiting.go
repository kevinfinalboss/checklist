package middlewares

import (
	"net/http"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware() gin.HandlerFunc {
	limiter := tollbooth.NewLimiter(30.0/60.0, nil)
	limiter.SetMessage("Você atingiu a taxa limite de requisições, por favor, tente novamente mais tarde.")
	limiter.SetMessageContentType("application/json; charset=utf-8")
	limiter.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})
	limiter.SetOnLimitReached(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"message": "Limite de requisições atingido"}`, http.StatusTooManyRequests)
	})

	return tollbooth_gin.LimitHandler(limiter)
}
