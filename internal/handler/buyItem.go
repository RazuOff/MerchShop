package handler

import (
	"net/http"

	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) BuyItem(c *gin.Context) {
	userID := c.GetString("id")
	if userID == "" {
		c.JSON(http.StatusInternalServerError, models.ErrorResponce{Errors: "context does not contains user ID"})
		return
	}

	itemName := c.Param("item")
	if itemName == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponce{Errors: "request does not contains 'item' query param"})
		return
	}

	exist, err := h.coinService.BuyItem(itemName, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponce{Errors: err.Error()})
		return
	}

	if !exist {
		c.JSON(http.StatusBadRequest, models.ErrorResponce{Errors: "this item does not exists"})
		return
	}

	c.Status(http.StatusOK)
}
