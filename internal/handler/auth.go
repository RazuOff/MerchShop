package handler

import (
	"net/http"

	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/gin-gonic/gin"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponce struct {
	Token string `json:"token"`
}

func (h *Handler) Auth(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponce{Errors: "Incorrect request body"})
		return
	}

	user, err := h.service.RegistrateOrLogin(creds.Username, creds.Password)
	if err != nil {
		c.JSON(err.Code, models.ErrorResponce{Errors: err.TextError})
		return
	}

	token, err := h.service.GenerateToken(creds.Username, user.ID, *h.config)
	if err != nil {
		c.JSON(err.Code, models.ErrorResponce{Errors: err.TextError})
		return
	}

	c.JSON(http.StatusOK, AuthResponce{Token: token})
}
