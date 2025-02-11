package handler

import (
	"net/http"

	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/gin-gonic/gin"
)

type requestBody struct {
	ToUserLogin string `json:"toUser"`
	Amount      int    `json:"amount"`
}

func (h *Handler) SendCoins(c *gin.Context) {
	var body requestBody

	userID := c.GetInt("id")
	if userID == 0 {
		c.JSON(http.StatusInternalServerError, models.ErrorResponce{Errors: "context does not contains user ID"})
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponce{Errors: "incorrect request body"})
		return
	}

	if err := h.coinService.SendCoins(userID, body.ToUserLogin, body.Amount); err != nil {
		c.JSON(err.Code, models.ErrorResponce{Errors: err.TextError})
		return
	}

	c.Status(http.StatusOK)
}
