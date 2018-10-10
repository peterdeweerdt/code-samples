package core

import (
	"strconv"
)

type PaymentMethod string

var PaymentMethodCreditCard PaymentMethod = "Credit Card"
var PaymentMethodApplePay PaymentMethod = "Apple Pay"
var PaymentMethodAndroidPay PaymentMethod = "Android Pay"

type LegacyPaymentInfo struct {
	CardType    string                 `json:"card_type,omitempty"`
	CardNumber  string                 `json:"number,omitempty"`
	CardExpiry  string                 `json:"expiry,omitempty"`
	CardName    string                 `json:"name,omitempty"`
	CardCVV     string                 `json:"cvv,omitempty"`
	CardZip     string                 `json:"zip,omitempty"`
	CardLast4   string                 `json:"card_last_4,omitempty"`
	CardDefault bool                   `json:"make_default"`
	Token       string                 `json:"token"`
	StripeToken string                 `json:"stripe_token"`
	Amount      int                    `json:"amount"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
	VaultID     string                 `json:"vault_id"`
	CustomerID  DatabaseID
	OrderID     DatabaseID
}

func (p *LegacyPaymentInfo) Tip() int {
	tip, err := strconv.Atoi(p.Metadata["tip"].(string))
	if err != nil {
		tip = 0
	}
	return tip
}

func (p *LegacyPaymentInfo) Method() PaymentMethod {
	switch p.CardType {
	case string(PaymentMethodApplePay):
		return PaymentMethodApplePay
	case string(PaymentMethodAndroidPay):
		return PaymentMethodAndroidPay
	default:
		return PaymentMethodCreditCard
	}
}
