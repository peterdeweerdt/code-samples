package payments

import (
	"fmt"

	"github.com/pkg/errors"
	"core"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

type Stripe struct {
	SDKKey string
	NoOp   bool
}

func (s Stripe) MakePayment(p core.TokenizedPayment) (string, error) {
	if s.NoOp {
		return fmt.Sprintf("%d-noop-transaction-id", p.OrderID), nil
	}

	chargeParams := stripe.ChargeParams{
		Amount:   uint64(p.Amount + p.Tip),
		Currency: "usd",
		Source:   &stripe.SourceParams{Token: p.Token},
	}
	chargeParams.AddMeta("tip", fmt.Sprintf("%d", p.Tip))
	chargeParams.AddMeta("order_id", fmt.Sprintf("%d", p.OrderID))
	if p.CustomerID != 0 {
		chargeParams.AddMeta("customer_id", fmt.Sprintf("%d", p.CustomerID))
	}

	stripe.Key = s.SDKKey
	ch, err := charge.New(&chargeParams)
	if err != nil {
		return "", errors.Wrapf(err, "paymentHandler: error making stripe payment")
	}

	return ch.ID, nil
}
