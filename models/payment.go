package model

// Payment is like a model but since it isn't related to the database
// it is preferable to put it in 'lib'. It's a third party so it has to be
// updatable easily.

import (
	"fmt"
	"strings"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/account"
	"github.com/stripe/stripe-go/bankaccount"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
	"github.com/stripe/stripe-go/transfer"
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

// RegisterBank links a bank account to a customer
func RegisterBank(email string, token string) (*stripe.Account, error) {
	accParams := &stripe.AccountParams{
		Managed: true,
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

// PayUser transfer money to the customer
func PayUser(accID string, amount uint64) (*stripe.Transfer, error) {
	tParams := &stripe.TransferParams{
		Amount:   int64(amount),
		Currency: "eur",
		Dest:     accID,
		Desc:     "Wooble.io funds",
	}

	return transfer.New(tParams)
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
