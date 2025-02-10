package middleware

import (
	"net/http"

	"github.com/RazuOff/MerchShop/internal/config"
	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// func GenerateToken(username string, role string, userID int) (string, error) {
// 	expirationTime := time.Now().Add(5 * time.Minute)
// 	claims := &Claims{
// 		ID:       userID,
// 		Username: username,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: expirationTime.Unix(),
// 		},
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(jwtKey)
// }

func AuthMiddleware(config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JwtKey, nil
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
