package core

import (
	"time"

	"github.com/nu7hatch/gouuid"
)

type Token struct {
	ID         DatabaseID `json:"-"`
	Service    string     `json:"service"`
	Name       string     `json:"name"` // Name is what type of token: access, refresh, etc
	Token      string     `json:"token"`
	CustomerID DatabaseID `json:"customer_id"`
	Expiry     time.Time  `json:"expiry"`
}

type TokenEmail struct {
	Token  string
	Expiry string
}

func (app AppContext) CreateTokenForCustomerWithID(service, name string, customerID DatabaseID, expiry time.Time) (Token, error) {
	token := Token{}

	u4, err := uuid.NewV4()
	if err != nil {
		return token, err
	}

	token.Service = service
	token.Name = name
	token.Token = u4.String()
	token.CustomerID = customerID
	token.Expiry = expiry
	err = app.DB.InsertToken(&token)
	return token, err
}

func (app AppContext) ReplaceTokensForCustomerWithID(service, name string, customerID DatabaseID, expiry time.Time) (Token, error) {
	//Clear all the tokens then add one. This handles the edge case for forgetting password currently, will sign out all clients with a given customerID for password reset
	app.DB.DeleteTokens(DatabaseID(customerID))

	return app.CreateTokenForCustomerWithID(service, name, customerID, expiry)
}
