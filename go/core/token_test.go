package core_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"core"
	"handlers"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateTokenForCustomerWithID(t *testing.T) {
	var app core.AppContext
	defer testServer(&app)()

	c := createTestCustomer(app)
	token, err := app.CreateTokenForCustomerWithID("test-service", "test-name", c.ID, time.Now())
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, token.CustomerID, c.ID)
}

func TestReplaceTokensForCustomerWithID(t *testing.T) {
	var app core.AppContext
	defer testServer(&app)()

	c := createTestCustomer(app)
	token, _ := app.CreateTokenForCustomerWithID("rize", "access", c.ID, time.Now())
	newToken, err2 := app.ReplaceTokensForCustomerWithID("rize", "access", c.ID, time.Now())
	assert.NoError(t, err2)

	cases := []struct {
		token core.Token
		code  int
	}{
		{token, 401},
		{newToken, 200},
	}

	for _, c := range cases {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", c.token.Token)
		res := httptest.NewRecorder()
		handlers.APIHandlers(app).ServeHTTP(res, req)
		assert.Equal(t, c.code, res.Code)
	}
}

func createTestCustomer(app core.AppContext) *core.Customer {
	//hashing in the test could go away if we move it down to the DB layer
	hash, _ := bcrypt.GenerateFromPassword([]byte("testCustomerPassword"), bcrypt.DefaultCost)

	c := &core.Customer{
		FirstName:  "Bob",
		LastName:   "Smith",
		Email:      "bob@smith.com",
		Password:   string(hash),
		ExternalID: "4321"}
	app.AddCustomer(c)

	return c
}
