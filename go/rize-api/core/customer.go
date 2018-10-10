package core

import (
	"time"

	"github.com/pkg/errors"
)

// A Customer is a person who is buying a pizza in this system.
type Customer struct {
	ID          DatabaseID `json:"id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	ImageURL    string     `json:"image"`
	Phone       string     `json:"phone"`
	Password    string     `json:"-"`
	ExternalID  string     `json:"-"`
	PosID       KountaID   `json:"-"`
	ServiceName string     `json:"service_name,omitempty"`
}

// Credential represents the login information
type Credential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// PasswordUpdate represents a request to set a password after forgotten
type PasswordUpdate struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// Connection represents an external account.  ex: Facebook
type Connection struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	Service            string    `json:"service"`
	AccessToken        string    `json:"access_token"`
	AccessTokenExpiry  time.Time `json:"access_token_expiry"`
	RefreshToken       string    `json:"refresh_token"`
	RefreshTokenExpiry time.Time `json:"refresh_token_expiry"`
}

func (app AppContext) AddCustomer(c *Customer) error {
	posCustomer, err := app.Kounta.GetCustomerByEmail(c.Email)
	if err != nil {
		return errors.Wrap(err, "add customer")
	}

	if posCustomer == nil {
		posCustomer, err = app.Kounta.CreateCustomer(c.Email, c.Email, "", c.Phone, c.ID)
		if err != nil {
			return errors.Wrap(err, "add customer")
		}
	}

	c.PosID = posCustomer.GetPosID()

	if err := app.DB.InsertCustomer(c); err != nil {
		return errors.Wrap(err, "add customer")
	}

	return nil
}
