package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getOrderById(c *gin.Context) {

	orderId := c.Param("id")
	if orderId == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid param")
		return
	}

	order, err := h.services.OrderCRUD.GetById(orderId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, order)
}
