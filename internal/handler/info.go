package handler

import (
	"net/http"

	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetInfo(c *gin.Context) {

	userID := c.GetInt("id")
	if userID == 0 {
		c.JSON(http.StatusInternalServerError, models.ErrorResponce{Errors: "context does not contains user ID"})
		return
	}

	resp, err := h.service.GenerateInfo(userID)
	if err != nil {
		c.JSON(err.Code, models.ErrorResponce{Errors: err.TextError})
		return
	}

	c.JSON(http.StatusOK, resp)

}
