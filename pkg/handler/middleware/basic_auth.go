package middleware

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

func BasicAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, password, hasAuth := ctx.Request.BasicAuth()

		if hasAuth && user == viper.GetString("auth.user") && password == os.Getenv("AUTH_BASIC_PASSWORD") {
			log.Infof("Authorised user: %s", user)
			ctx.Next()
		} else {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			ctx.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			log.Errorf("Authorisation failed for user: %s", user)
			return
		}
	}
}
