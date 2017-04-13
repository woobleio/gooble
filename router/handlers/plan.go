package handler

import (
	model "wooble/models"

	"github.com/gin-gonic/gin"
)

// GETPlans is handlers that returns all Wooble plans
func GETPlans(c *gin.Context) {
	plans, err := model.AllPlans()

	if err != nil {
		c.Error(err).SetMeta(ErrServ)
		return
	}

	c.JSON(OK, NewRes(plans))
}
