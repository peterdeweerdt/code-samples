package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/peterdeweerdt/rize-api-go/core"
	"github.com/peterdeweerdt/rize-api-go/handlers"
	"github.com/peterdeweerdt/rize-api-go/payments"
	"github.com/peterdeweerdt/rize-api-go/pos"
	"github.com/peterdeweerdt/rize-api-go/sendgrid"
	"github.com/peterdeweerdt/rize-api-go/sk"
)

type environment struct {
	serverURL             string
	dbURL                 string
	kountaKey             string
	cayanUserName         string
	cayanPassword         string
	cardConnectURL        string
	cardConnectUserName   string
	cardConnectPassword   string
	cardConnectMerchantID int64
	sendgridAPIKey        string
	stripeKey             string
	siteWhitelist         []int64
	isNoOpPayments        bool
}

func main() {
	log.Println("Starting up...")
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("[RIZE] ")

	env := loadEnvironment()

	// Database
	pg, err := core.InitPG(env.dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	// Point of sale
	kounta := pos.NewKounta(
		sk.HTTPClient{
			BaseURL:     "https://api.kounta.com/v1/companies/20246",
			BasicAuth:   env.kountaKey,
			ContentType: "application/json",
			Logging:     true,
		},
		env.serverURL)

	// Payment gateways
	cayan := payments.Cayan{
		HTTP: sk.HTTPClient{
			BaseURL:     "https://smbpos.transactiongateway.com/api",
			ContentType: "application/xml",
			Logging:     true,
		},
		UserName: env.cayanUserName,
		Password: env.cayanPassword,
		NoOp:     env.isNoOpPayments,
	}

	stripe := payments.Stripe{
		SDKKey: env.stripeKey,
		NoOp:   env.isNoOpPayments,
	}

	cardConnectAuth := base64.StdEncoding.EncodeToString([]byte(env.cardConnectUserName + ":" + env.cardConnectPassword))
	cardConnect := payments.CardConnect{
		HTTP: sk.HTTPClient{
			BaseURL:     env.cardConnectURL,
			BasicAuth:   cardConnectAuth,
			ContentType: "application/json",
			Logging:     true,
		},
		MerchantID: env.cardConnectMerchantID,
	}

	app := core.AppContext{
		DB:            pg,
		Kounta:        kounta,
		Cayan:         cayan,
		Stripe:        stripe,
		CardConnect:   cardConnect,
		Mailer:        sendgrid.API{Key: env.sendgridAPIKey},
		SiteWhitelist: env.siteWhitelist,
	}

	log.Println("Updating all menus")
	go app.UpdateAllMenus()

	log.Println("Starting server")
	err = http.ListenAndServe(":"+os.Getenv("PORT"), handlers.APIHandlers(app))
	log.Fatal(err)
}

func loadEnvironment() environment {
	herokuAppName := os.Getenv("HEROKU_APP_NAME")
	if herokuAppName == "" {
		log.Fatal("HEROKU_APP_NAME required!")
	}
	serverURL := fmt.Sprintf("https://%s.herokuapp.com", herokuAppName)
	fmt.Printf("Setting SERVER_URL to %s\n", serverURL)

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL required!")
	}

	kountaKey := os.Getenv("KOUNTA_KEY")
	if kountaKey == "" {
		log.Fatal("KOUNTA_KEY is required!")
	}

	cayanUserName := os.Getenv("CAYAN_USERNAME")
	cayanPassword := os.Getenv("CAYAN_PASSWORD")
	if cayanUserName == "" || cayanPassword == "" {
		log.Fatal("CAYAN_USERNAME and CAYAN_PASSWORD are required!")
	}

	cardConnectURL := os.Getenv("CC_URL")
	cardConnectUserName := os.Getenv("CC_USERNAME")
	cardConnectPassword := os.Getenv("CC_PASSWORD")
	cardConnectMerchantID, _ := strconv.ParseInt(os.Getenv("CC_MERCHANT_ID"), 10, 64)
	if (cardConnectURL == "") ||
		(cardConnectUserName == "") ||
		(cardConnectPassword == "") ||
		(cardConnectMerchantID == 0) {
		log.Fatal("CC_URL, CC_USERNAME, CC_PASSWORD, and CC_MERCHANT_ID are required!")
	}

	sendgridAPIKey := os.Getenv("SENDGRID_API_KEY")
	if sendgridAPIKey == "" {
		log.Fatal("SENDGRID_API_KEY is required!")
	}

	stripeKey := os.Getenv("STRIPE_KEY")
	if stripeKey == "" {
		log.Fatal("STRIPE_KEY is required!")
	}

	whitelist := os.Getenv("LRS_SITE_ID_WHITELIST")
	siteIDs, err := sk.MapStringsToInt64s(strings.Split(whitelist, ","))
	if err != nil {
		log.Fatal(err)
	}

	noopPayments := os.Getenv("NOOP_PAYMENTS")
	isNoOpPayments := noopPayments != "" &&
		strings.ToLower(noopPayments) != "false" &&
		strings.ToLower(noopPayments) != "off" &&
		strings.ToLower(noopPayments) != "disabled"

	return environment{
		serverURL:             serverURL,
		dbURL:                 dbURL,
		kountaKey:             kountaKey,
		cayanUserName:         cayanUserName,
		cayanPassword:         cayanPassword,
		cardConnectURL:        cardConnectURL,
		cardConnectUserName:   cardConnectUserName,
		cardConnectPassword:   cardConnectPassword,
		cardConnectMerchantID: cardConnectMerchantID,
		sendgridAPIKey:        sendgridAPIKey,
		stripeKey:             stripeKey,
		siteWhitelist:         siteIDs,
		isNoOpPayments:        isNoOpPayments,
	}
}
