package handler

import (
	"net/http"

	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponce{Errors: "Authorization is empty"})
			c.Abort()
			return
		}

		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return h.config.JwtKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponce{Errors: err.Error()})
			c.Abort()
			return
		}
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, models.ErrorResponce{Errors: "Token is not valid"})
			c.Abort()
			return
		}

		c.Set("id", claims.ID)
		c.Next()
	}
}
