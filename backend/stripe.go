package main

import (
	"fmt"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
)

const (
	STRIPE_CUSTOMER_DESC = "%s %s (id: %d)"
)

// Initializes the Stripe API client
func SetupStripe(env *Environment) {
	stripe.Key = env.stripeAPIKey
}

// Creates a new Stripe customer; returns the customer id
func NewStripeCustomerId(email string, id int64, firstName string, lastName string) (string, error) {
	params := &stripe.CustomerParams{
		Email: email,
		Desc:  fmt.Sprintf(STRIPE_CUSTOMER_DESC, firstName, lastName, id),
	}
	newCust, err := customer.New(params)
	if err != nil {
		return "", err
	} else {
		return newCust.ID, nil
	}
}
