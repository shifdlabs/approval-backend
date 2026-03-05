package middleware

import (
	"Microservice/helper"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := helper.ExtractToken(c)

		err := helper.ValidateToken(accessToken)

		if err != nil {
			helper.PrintValue("Error Auth Middleware", "")
			fileName, atLine := helper.GetFileAndLine(*err)
			helper.ResponseError(c, helper.CustomError{
				Code:     401,
				Message:  "Check you credential.",
				FileName: fileName,
				AtLine:   atLine,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func DeserializeUser(DB *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Write code for Serialize Bearer Token.
		var access_token string
		authorization := ctx.GetHeader("Authorization")
		// accessTokenFromCookie, errorCookie := ctx.Cookie("access_token")
		access_token = strings.TrimPrefix(authorization, "Bearer ")

		// env, _ := config.LoadConfig(".")

		errTokenClaims := helper.ValidateToken(access_token)

		if errTokenClaims != nil {
			fileName, atLine := helper.GetFileAndLine(*errTokenClaims)
			helper.ResponseError(ctx, helper.CustomError{
				Code:     401,
				Message:  "Unauthorize.",
				FileName: fileName,
				AtLine:   atLine,
			})
			ctx.Abort()
		}

		ctx.Next()
	}
}
