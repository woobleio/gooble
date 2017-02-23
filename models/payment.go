package model

// Payment is like a model but since it isn't related to the database
// it is preferable to put it in 'lib'. It's a third party so it has to be
// updatable easily.

import (
	"fmt"
	"strings"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
)

// NewCustomer creates a new customer in the payment system
func NewCustomer(custEmail string, plan string, token string) (*stripe.Customer, error) {
	custParams := &stripe.CustomerParams{
		Email: custEmail,
		Plan:  plan,
	}

	if token != "" {
		custParams.SetSource(token)
	}

	return customer.New(custParams)
}

// ChargeCustomerForCreations charges the customer "custID" for the given creations "objIDs"
func ChargeCustomerForCreations(custID string, amount uint64, objIDs []string) (*stripe.Charge, error) {
	params := &stripe.ChargeParams{
		Amount:    amount,
		Currency:  "eur",
		Desc:      fmt.Sprintf("Creations:%s", strings.Join(objIDs, " | ")),
		Customer:  custID,
		NoCapture: true,
	}

	return charge.New(params)
}

// ChargeOneTimeForCreations charges an account without recording it
func ChargeOneTimeForCreations(amount uint64, objIDs []string, token string) (*stripe.Charge, error) {
	params := &stripe.ChargeParams{
		Amount:    amount,
		Currency:  "eur",
		Desc:      fmt.Sprintf("Creations:%s", strings.Join(objIDs, " | ")),
		NoCapture: true,
	}
	params.SetSource(token)

	return charge.New(params)
}

// CaptureCharge applies the charge
func CaptureCharge(chargeID string) (*stripe.Charge, error) {
	return charge.Capture(chargeID, nil)
}
