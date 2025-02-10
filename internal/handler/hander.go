package handler

import (
	"github.com/RazuOff/MerchShop/internal/config"
	"github.com/RazuOff/MerchShop/internal/middleware"
	"github.com/RazuOff/MerchShop/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	authService service.AuthService
	infoService service.InfoService
	coinService service.CoinService

	config *config.Config
}

func NewHandler(services *service.Service, config *config.Config) *Handler {
	return &Handler{
		authService: services.AuthService,
		infoService: services.InfoService,
		coinService: services.CoinService,

		config: config,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/api")
	{
		auth.POST("/auth", h.Auth)
	}

	api := router.Group("/api", middleware.AuthMiddleware(*h.config))
	{
		api.GET("/info", h.GetInfo)
		api.POST("/sendCoin", h.SendCoins)
		api.GET("/buy/:item", h.BuyItem)
	}

	return router
}
