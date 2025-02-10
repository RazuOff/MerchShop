package handler

import (
	"net/http"

	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetInfo(c *gin.Context) {

	userID := c.GetString("id")
	if userID == "" {
		c.JSON(http.StatusInternalServerError, models.ErrorResponce{Errors: "context does not contains user ID"})
		return
	}

	resp, err := h.infoService.GenerateInfo(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponce{Errors: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)

}
