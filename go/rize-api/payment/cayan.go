package payments

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"core"
	"pjd"
)

// online documentation:
// https://smbpos.transactiongateway.com/merchants/resources/integration/integration_portal.php?tid=f736bd0b4883b509c0ad26f73a294229#transaction_response_variables

const paymentGatewayTransactionApproved = "1"

type Cayan struct {
	HTTP     pjd.HTTPClient
	UserName string
	Password string
	NoOp     bool // set this to true to prevent payments from processing.  usefull for testing on non-production environments
}

type keyResponse struct {
	Error  string `xml:"error_response"`
	SDKKey string `xml:"sdk_key"`
}

func (c Cayan) GetSDKKey() (string, error) {
	response := keyResponse{}

	params := url.Values{}
	params.Add("username", c.UserName)
	params.Add("password", c.Password)
	params.Add("report_type", "sdk_key")

	_, err := c.HTTP.Get("/query.php?"+params.Encode(), &response)
	if err != nil {
		return "", core.PaymentGatewayError{
			Reason: errors.Wrap(err, "payment_gateway: error getting sdk_key").Error(),
		}
	}

	return response.SDKKey, nil
}

func (c Cayan) MakePayment(paymentInfo core.LegacyPaymentInfo) (string, string, error) {
	if c.NoOp {
		return fmt.Sprintf("%d-noop-transaction-id", paymentInfo.OrderID), "noop-vault-id", nil
	}

	queryString := paymentQueryString(paymentInfo)
	if len(queryString) == 0 {
		return "", "", errors.New("payment_gateway: malformed payment request")
	}

	reqBody := bytes.NewReader([]byte(queryString))
	req, err := http.NewRequest("POST", "https://secure.networkmerchants.com/api/transact.php", reqBody)
	if err != nil {
		return "", "", errors.Wrap(err, "payment_gateway: unable to create request")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	log.Println(pjd.MustDumpRequest(req))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", "", errors.Wrap(err, "payment_gateway: unable to save and/or charge card")
	}
	defer res.Body.Close()

	log.Println(pjd.MustDumpResponse(res))

	resBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", errors.Wrap(err, "payment_gateway: unable to read response body")
	}

	results, err := url.ParseQuery(string(resBodyBytes))
	if err != nil {
		return "", "", errors.Wrap(err, "payment_gateway: unable to process response values")
	}

	responseNumber := results.Get("response")
	if responseNumber != paymentGatewayTransactionApproved {
		responseText := results.Get("responsetext")
		return "", "", core.PaymentGatewayError{Reason: responseText}
	}

	return results.Get("transactionid"), results.Get("customer_vault_id"), nil
}

func paymentQueryString(p core.LegacyPaymentInfo) string {
	var action string
	if len(p.VaultID) > 0 {
		action = "update_customer"
	} else {
		action = "add_customer"
	}

	amount := fmt.Sprintf("%.2f", pjd.Round(float64(p.Amount+p.Tip())/100.0, .005, 2))
	params := map[string]string{
		"username": os.Getenv("CAYAN_USERNAME"),
		"password": os.Getenv("CAYAN_PASSWORD"),
		"amount":   amount,
		"type":     "sale",
		"currency": "USD",
		"orderid":  strconv.FormatInt(int64(p.OrderID), 10),
	}

	params["encrypted_payment"] = p.Token
	params["zip"] = p.CardZip
	if p.CardDefault {
		params["customer_vault"] = action
	}
	if len(p.VaultID) > 0 {
		params["customer_vault_id"] = p.VaultID
	}
	nameParts := strings.Split(p.CardName, " ")
	if len(nameParts) > 0 {
		params["first_name"] = nameParts[0]
		if len(nameParts) > 1 {
			params["last_name"] = strings.Join(nameParts[1:len(nameParts)-1], " ")
		}
	}

	query := ""
	for k, v := range params {
		query += k + "=" + url.QueryEscape(v) + "&"
	}

	if len(query) == 0 {
		return ""
	}
	return query[:len(query)-1]
}
