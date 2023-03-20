package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"http-server/cmd/server/config"
	"http-server/internal/services"

	"github.com/gin-gonic/gin"
)

func MiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		config, _ := config.LoadConfig(".")
		jwtService := services.GetJwtService(&config)

		token, err := ctx.Cookie("token")
		if err != nil {
			print(err.Error())

			refreshToken, err := ctx.Cookie("refresh_token")
			if err != nil {

				print(err.Error())
				// Session timeout
				ctx.JSON(http.StatusBadGateway, gin.H{"message": "Require login"})
				ctx.Abort()
				return
			}

			// Gen new token if refresh token exist
			tokenExpHour, _ := strconv.Atoi(config.JWT_EXPIRY_HOUR)
			tokenExp := time.Now().Add(time.Hour * time.Duration(tokenExpHour)).Unix()
			newToken, _ := jwtService.RefreshToken(refreshToken, tokenExp)
			ctx.SetCookie("token", newToken, tokenExpHour, "/", "localhost", false, true)
			print("New Token: %s", newToken)
		}

		// Decode token
		id, err := jwtService.ExtractIDFromToken(token)
		if err != nil {
			fmt.Println("Cannot parse token")
			ctx.Abort()
			return
		}
		ctx.Set("user-uuid", id)
		ctx.Next()
	}
}
