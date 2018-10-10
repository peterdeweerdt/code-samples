package core

const (
	GatewayCardConnect = "card_connect"
	GatewayStripe      = "stripe"
)

type TokenizedPayment struct {
	Gateway    string     `json:"gateway"`
	Token      string     `json:"token"`
	SiteID     KountaID   `json:"site_id"`
	OrderID    DatabaseID `json:"order_id"`
	CustomerID DatabaseID `json:"-"` // CustomerID not transmitted over JSON
	Amount     int        `json:"amount"`
	Tip        int        `json:"tip"`
	Expiry     string     `json:"expiry"` // Expiry only needed for CardConnect
}

type Cayan interface {
	GetSDKKey() (key string, err error)
	MakePayment(info LegacyPaymentInfo) (transactionID string, vaultID string, err error)
}

type Stripe interface {
	MakePayment(p TokenizedPayment) (transactionID string, err error)
}

type CardConnect interface {
	HealthCheck() error
	// AddCreditCard will update the VaultID on success
	AddCreditCard(info *CreditCard) error
	GetCreditCards(vaultID string) ([]CreditCard, error)
	// UpdateCreditCard will update the VaultID on success
	UpdateCreditCard(info *CreditCard) error
	DeleteCreditCard(vaultID string) error
	MakePayment(p TokenizedPayment) (transactionID string, err error)
}

type PaymentGatewayError struct {
	Reason string
}

func (e PaymentGatewayError) Error() string {
	return e.Reason
}
