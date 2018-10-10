package payments

import (
	"testing"

	"core"
	"github.com/stretchr/testify/assert"
)

func TestNoOp(t *testing.T) {
	api := Cayan{
		NoOp: true,
	}
	transactionID, vaultID, err := api.MakePayment(core.LegacyPaymentInfo{OrderID: 123})
	assert.NoError(t, err)
	assert.Equal(t, "123-noop-transaction-id", transactionID)
	assert.Equal(t, "noop-vault-id", vaultID)
}
