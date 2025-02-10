package handler

import (
	"net/http"

	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/gin-gonic/gin"
)

type requestBody struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

func (h *Handler) SendCoins(c *gin.Context) {
	var body requestBody

	userID := c.GetString("id")
	if userID == "" {
		c.JSON(http.StatusInternalServerError, models.ErrorResponce{Errors: "context does not contains user ID"})
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponce{Errors: "incorrect request body"})
		return
	}

	if err := h.coinService.SendCoins(userID, body.ToUser, body.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponce{Errors: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
