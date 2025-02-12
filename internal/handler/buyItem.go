package handler

import (
	"net/http"

	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) BuyItem(c *gin.Context) {
	userID := c.GetInt("id")
	if userID == 0 {
		c.JSON(http.StatusInternalServerError, models.ErrorResponce{Errors: "context does not contains user ID"})
		return
	}

	itemName := c.Param("item")
	if itemName == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponce{Errors: "request does not contains 'item' query param"})
		return
	}

	err := h.service.BuyItem(itemName, userID)

	if err != nil {
		c.JSON(err.Code, models.ErrorResponce{Errors: err.TextError})
		return
	}

	c.Status(http.StatusOK)
}
