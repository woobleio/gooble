package model

// Payment is like a model but since it isn't related to the database
// it is preferable to put it in 'lib'. It's a third party so it has to be
// updatable easily.

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/account"
	"github.com/stripe/stripe-go/bankaccount"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
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

// CaptureCharge applies the charge
func CaptureCharge(chargeID string) (*stripe.Charge, error) {
	return charge.Capture(chargeID, nil)
}

// RegisterBank links a bank account to a customer
func RegisterBank(email string, token string) (*stripe.Account, error) {
	accParams := &stripe.AccountParams{
		Type:    stripe.AccountTypeCustom,
		Country: "FR",
		Email:   email,
	}
	acc, err := account.New(accParams)
	if err != nil {
		return nil, err
	}

	bkParams := &stripe.BankAccountParams{
		AccountID: acc.ID,
		Token:     token,
	}

	_, err = bankaccount.New(bkParams)
	return acc, err
}

// SubCustomer subscribes the customer "custID" to the plan "plan"
func SubCustomer(custID string, plan string, token string) (*stripe.Sub, error) {
	subParams := &stripe.SubParams{
		Customer: custID,
		Plan:     plan,
		Card: &stripe.CardParams{
			Token: token,
		},
	}

	return sub.New(subParams)
}

// UnsubCustomer unsubscribes customer from his current plan
func UnsubCustomer(custID string) error {
	cust, err := customer.Get(custID, nil)
	if err != nil {
		return err
	}

	subID := cust.Subs.Values[0].ID

	subParams := &stripe.SubParams{
		EndCancel: true,
	}
	_, err = sub.Cancel(subID, subParams)

	return err
}
