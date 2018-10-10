package payments

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"core"
	"pjd"
)

type CardConnect struct {
	HTTP       pjd.HTTPClient
	MerchantID int64 // MerchantID is currently the same for all sites, this may change
}

func (c CardConnect) HealthCheck() error {
	// Create a copy so we don't try to unmarshal JSON
	healthCheckHTTP := c.HTTP
	healthCheckHTTP.ContentType = "application/xml"

	// Adding a trailing # so the HTTP client doesn't complain
	_, err := healthCheckHTTP.Get("/cardconnect/rest/#", struct{}{})
	return err
}

func (c CardConnect) AddCreditCard(card *core.CreditCard) error {
	body := profileBody{
		MerchantID: strconv.FormatInt(c.MerchantID, 10),
		CardName:   card.Name,
		CardNumber: card.Number,
		CardExpiry: card.Expiry,
		CardCVV:    card.CVV,
		CardZip:    card.ZipCode,
		CardType:   card.Type,
	}

	// Profile already exists
	if card.VaultID != "" {
		// Strip account number from vaultID
		if slashIndex := strings.Index(card.VaultID, "/"); slashIndex != -1 {
			card.VaultID = card.VaultID[:slashIndex]
		}

		body.ProfileID = card.VaultID
		body.IsUpdate = "Y"
	}

	if card.IsDefault || card.VaultID == "" {
		body.IsDefault = "Y"
	}

	_, err := c.HTTP.Put("/cardconnect/rest/profile", body, &body)
	if err != nil {
		return errors.Wrap(err, "create CardConnect profile")
	}

	card.Number = body.CardToken
	card.VaultID = body.ReturnedProfileID + "/" + body.AccountID
	card.IsDefault = body.IsDefault == "Y"
	return nil
}

func (c CardConnect) GetCreditCards(vaultID string) ([]core.CreditCard, error) {
	// The trailing slash is to indicate get all cards for a profile
	if !strings.Contains(vaultID, "/") {
		vaultID += "/"
	}

	cardConnectCards := []profileBody{}
	_, err := c.HTTP.Get(fmt.Sprintf("/cardconnect/rest/profile/%s/%d", vaultID, c.MerchantID), &cardConnectCards)
	if err != nil {
		return nil, errors.Wrap(err, "get CardConnect profile")
	}

	cards := make([]core.CreditCard, len(cardConnectCards))
	for i, c := range cardConnectCards {
		cards[i] = core.CreditCard{
			Name:      c.CardName,
			Number:    c.CardToken,
			Expiry:    c.CardExpiry,
			CVV:       c.CardCVV,
			ZipCode:   c.CardZip,
			Type:      c.CardType,
			VaultID:   c.ReturnedProfileID + "/" + c.AccountID,
			IsDefault: c.IsDefault == "Y",
		}
	}

	return cards, nil
}

func (c CardConnect) UpdateCreditCard(card *core.CreditCard) error {
	if card.VaultID == "" {
		return core.PaymentGatewayError{
			Reason: "update CardConnect profile: missing vaultID",
		}
	}

	body := profileBody{
		MerchantID: strconv.FormatInt(c.MerchantID, 10),
		ProfileID:  card.VaultID,
		CardName:   card.Name,
		CardNumber: card.Number,
		CardExpiry: card.Expiry,
		CardCVV:    card.CVV,
		CardZip:    card.ZipCode,
		CardType:   card.Type,
		IsUpdate:   "Y",
	}

	if card.IsDefault {
		body.IsDefault = "Y"
	}

	_, err := c.HTTP.Put("/cardconnect/rest/profile", body, &body)
	if err != nil {
		return errors.Wrap(err, "update CardConnect profile")
	}

	card.Number = body.CardToken
	card.VaultID = body.ReturnedProfileID + "/" + body.AccountID
	return nil
}

func (c CardConnect) DeleteCreditCard(vaultID string) error {
	// If an accountID is provided, only that account is deleted
	// Otherwise, the entire profile is deleted, but must end with a slash
	if !strings.Contains(vaultID, "/") {
		vaultID += "/"
	}

	_, err := c.HTTP.Delete(fmt.Sprintf("/cardconnect/rest/profile/%s/%d", vaultID, c.MerchantID))
	if err != nil {
		return errors.Wrap(err, "delete CardConnect profile")
	}

	return nil
}

type profileBody struct {
	MerchantID        string `json:"merchid"`
	ProfileID         string `json:"profile,omitempty"`
	ReturnedProfileID string `json:"profileid,omitempty"`
	AccountID         string `json:"acctid,omitempty"`
	CardName          string `json:"name"`
	CardNumber        string `json:"account,omitempty"`
	CardToken         string `json:"token,omitempty"`
	CardExpiry        string `json:"expiry"`
	CardCVV           string `json:"cvv2"`
	CardZip           string `json:"postal"`
	CardType          string `json:"accttype,omitempty"`
	IsUpdate          string `json:"profileupdate,omitempty"`
	IsDefault         string `json:"defaultacct,omitempty"`
}

func (c CardConnect) MakePayment(p core.TokenizedPayment) (string, error) {
	body := paymentBody{
		Amount:         strconv.Itoa(p.Amount + p.Tip),
		Currency:       "USD",
		ExpirationDate: p.Expiry,
		MerchantID:     strconv.FormatInt(c.MerchantID, 10),
		OrderID:        strconv.FormatInt(int64(p.OrderID), 10),
		ShouldCapture:  "Y",
		ProfileID:      p.Token,
	}

	_, err := c.HTTP.Put("/cardconnect/rest/auth", body, &body)
	if err != nil {
		return "", errors.Wrap(err, "sending payment to CardConnect")
	}

	return body.TransactionID, nil
}

type paymentBody struct {
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	ExpirationDate string `json:"expiry"`
	MerchantID     string `json:"merchid"`
	OrderID        string `json:"orderid"`
	ShouldCapture  string `json:"capture"`
	ProfileID      string `json:"profile"`
	TransactionID  string `json:"retref,omitempty"`
}
