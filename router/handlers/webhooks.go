package handler

import (
	"errors"
	model "wooble/models"

	"github.com/gin-gonic/gin"
	stripe "github.com/stripe/stripe-go"
)

// POSTWebhooks post webhooks (for Stripe)
func POSTWebhooks(c *gin.Context) {
	var evt stripe.Event

	if err := c.BindJSON(&evt); err != nil {
		c.Error(err)
		return
	}

	if evt.Type != "invoice.payment_succeeded" {
		c.Error(errors.New("Not a valid Stripe Webhook type"))
		return
	}

	if err := model.RenewUserPlan(evt.Data.Obj["customer"].(string), int64(evt.Data.Obj["period_start"].(float64)), int64(evt.Data.Obj["period_end"].(float64))); err != nil {
		c.Error(err)
	}
}
