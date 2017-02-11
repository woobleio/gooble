package lib

// Payment is like a model but since it isn't related to the database
// it is preferable to put it in 'lib'. It's a third party so it has to be
// updatable easily.

import (
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go"
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

// LoadPayment loads payment API
func LoadPayment() {
	stripe.Key = viper.GetString("stripe_key")
}
