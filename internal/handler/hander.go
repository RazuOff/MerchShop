package handler

import (
	"github.com/RazuOff/MerchShop/internal/config"
	"github.com/RazuOff/MerchShop/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
	config  *config.Config
}

func NewHandler(services *service.Service, config *config.Config) *Handler {
	return &Handler{
		service: services,
		config:  config,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	auth := router.Group("/api")
	{
		auth.POST("/auth", h.Auth)
	}

	api := router.Group("/api", h.AuthMiddleware())
	{
		api.GET("/info", h.GetInfo)
		api.POST("/sendCoin", h.SendCoins)
		api.GET("/buy/:item", h.BuyItem)
	}

	return router
}
