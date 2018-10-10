package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	_ "github.com/lib/pq"
)

var (
	a, b, c, d            int
	rfc3339, rfc3339Micro string
)

type orderEvents struct {
	siteID    int
	notes     string
	orderID   int
	status    string
	total     float64
	totalTax  float64
	paid      float64
	tips      float64
	customer  int
	priceVar  float64
	deleted   bool
	createdAt time.Time
	updatedAt time.Time
}

func main() {

	// --------------------check postgres database connection------------------
	timeStart := time.Now()
	fmt.Println("... validate the connection arguments to rize analytics postgres database")
	db, err := sql.Open("postgres", "user=postgres dbname=rize_analytics port=5432 password=April2016! sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic(err.Error)
	}
	defer db.Close()
	fmt.Println("... responseTime - sql.Open(): ", time.Since(timeStart))

	timeStart = time.Now()
	fmt.Println("... validate the connection to the rize analytics postgress database")
	err = db.Ping()
	if err != nil {
		panic(err.Error)
	}
	fmt.Println("... responseTime - db.Ping(): ", time.Since(timeStart))

	// ------------------------------------------------------------------------
	file, err := os.Open("order.log")
	if err != nil {
		fmt.Println("err is: ", err)
	}
	defer file.Close()

	i := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		singleLine := scanner.Text()
		numberOfCharacters := len(singleLine)
		if numberOfCharacters > 1 {

			fmt.Println("-------------------------------------------------------------------------------")
			fmt.Println("working on line number: ", i)
			// fmt.Println(singleLine)

			// SiteID:
			siteID := extractInt(singleLine, "SiteID:", " CallbackURI:")
			fmt.Printf("type: %T        siteID:              %v \n", siteID, siteID)

			// Notes:
			notes := extractString(singleLine, "Notes:", " ID:")
			fmt.Printf("type: %T       notes:               %v \n", notes, notes)

			// Extracting Table Number from Notes
			a = strings.Index(notes, "[") + 2
			b = strings.Index(notes, "]")
			// fmt.Println("a:", a, "b:", b)
			var tableNumberString string
			if a != -1 && b != -1 {
				tableNumberString = notes[a:b]
				if err != nil {
					fmt.Println(err)
				}
			}
			tableNumber, err := strconv.ParseInt(tableNumberString, 10, 64)
			fmt.Printf("type: %T        tableNumber:         %v \n", tableNumber, tableNumber)

			// Extract Pager Number from Notes

			var digit1String, digit2String string
			var digit1 bool
			var digit2 bool
			var pagerNumber int64

			runeNotes := []rune(notes)
			if len(runeNotes) > 1 {
				digit1 = unicode.IsDigit(runeNotes[0])
				digit2 = unicode.IsDigit(runeNotes[1])
			}
			// fmt.Println("digit1:", digit1)
			// fmt.Println("digit2:", digit2)

			if digit1 {
				digit1String = notes[0:1]
				//	fmt.Println("digit1String:", digit1String)
			}
			if digit2 {
				digit2String = notes[1:2]
				//	fmt.Println("digit1String:", digit2String)
			}
			pagerNumberString := digit1String + digit2String
			pagerNumber, err = strconv.ParseInt(pagerNumberString, 10, 64)

			fmt.Printf("type: %T        pagerNumber:         %v \n", pagerNumber, pagerNumber)

			// Order ID:
			a = len(" ID:")
			b = strings.Index(singleLine, " ID:")
			d = b + a
			c = d + 9
			orderIDString := singleLine[d:c]
			orderID, err := strconv.Atoi(orderIDString)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("type: %T          orderID:             %v \n", orderID, orderID)

			// Status:
			a = len(" Status:")
			b = strings.LastIndex(singleLine, " Status:")
			c = strings.Index(singleLine, " Total:")
			d = b + a
			status := singleLine[d:c]
			fmt.Printf("type: %T       status:              %v \n", status, status)

			// Total:
			total := extractFloat(singleLine, "Total:", " TotalTax:")
			fmt.Printf("type: %T      total:           %9.2f \n", total, total)

			// TotalTax:
			// totalTax := extractFloat(singleLine, "TotalTax:", " Paid:")
			// fmt.Printf("type: %T      totalTax:        %9.2f \n", totalTax, totalTax)

			// Paid:
			// a = len(" Paid:")
			// b = strings.LastIndex(singleLine, " Paid:")
			// c = strings.Index(singleLine, " Tips:")
			// d = b + a
			// paidString := singleLine[d:c]
			// paid, err := strconv.ParseFloat(paidString, 64)
			// if err != nil {
			//	fmt.Println("error with strconv.ParseFloat(stringData)", err)
			// }
			// fmt.Printf("type: %T      paid:            %9.2f \n", paid, paid)

			// Tips:
			// tips := extractFloat(singleLine, "Tips:", "Customer:")
			// fmt.Printf("type: %T      tips:            %9.2f \n", tips, tips)

			// Customer:
			// customer := extractString(singleLine, "Customer:", " PriceVar:")
			// fmt.Printf("type: %T       customer:            %v \n", customer, customer)

			// PriceVar:
			// priceVar := extractFloat(singleLine, " PriceVar:", " Payments:")
			// fmt.Printf("type: %T      priceVar:        %9.2f \n", priceVar, priceVar)

			// Deleted:
			a = len(" Deleted:")
			b = strings.Index(singleLine, " Deleted:")
			d = b + a
			c = d + 5
			deletedString := singleLine[d:c]
			deleted, err := strconv.ParseBool(deletedString)
			if err != nil {
				fmt.Println("error with strconv.ParseBool(deleteString)", err)
			}
			fmt.Printf("type: %T         deleted:             %v \n", deleted, deleted)

			// CreatedAt:
			createdAtString := extractString(singleLine, "CreatedAt:", " UpdatedAt:")

			rfc3339 = "2006-01-02 15:04:05 -0700 MST"
			createdAt, err := time.Parse(rfc3339, createdAtString)
			if err != nil {
				fmt.Println(err)
			}
			// createdAt := t1.Add(-4 * time.Hour)
			fmt.Printf("type: %T    createdAt:           %v \n", createdAt, createdAt)

			// UpdatedAt:
			a = len(" UpdatedAt:")
			b = strings.Index(singleLine, " UpdatedAt:")
			d = b + a
			c = d + 29
			updatedAtString := singleLine[d:c]

			rfc3339 = "2006-01-02 15:04:05 -0700 MST"
			updatedAt, err := time.Parse(rfc3339, updatedAtString)
			if err != nil {
				fmt.Println(err)
			}
			// updatedAt := t1.Add(-4 * time.Hour)
			fmt.Printf("type: %T    updatedAt:           %v \n", updatedAt, updatedAt)

			// insert data into Postgres database
			sqlStatement := "INSERT INTO order_events (site_id, notes, table_number, pager_number, order_id, status, total, deleted, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)RETURNING record_number"
			var recordNumber int
			err = db.QueryRow(sqlStatement, siteID, notes, tableNumber, pagerNumber, orderID, status, total, deleted, createdAt, updatedAt).Scan(&recordNumber)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}

			i++
		}
	}
}

func extractString(singleLine string, startTag string, endTag string) string {

	a := len(startTag)
	b := strings.Index(singleLine, startTag)
	c := strings.Index(singleLine, endTag)

	d := b + a
	// fmt.Println("len:", c-d, "start:", d, "end:", c)
	stringData := singleLine[d:c]

	return stringData
}

func extractInt(singleLine string, startTag string, endTag string) int64 {

	a := len(startTag)
	b := strings.Index(singleLine, startTag)
	c := strings.Index(singleLine, endTag)
	d := b + a
	// fmt.Println("len:", c-d, "start:", d, "end:", c)
	stringData := singleLine[d:c]
	intData, err := strconv.ParseInt(stringData, 10, 64)
	if err != nil {
		fmt.Println("error with strconv.ParseInt(stringData)", err)
	}
	return intData
}

func extractFloat(singleLine string, startTag string, endTag string) float64 {

	a := len(startTag)
	b := strings.Index(singleLine, startTag)
	c := strings.Index(singleLine, endTag)
	d := b + a
	// fmt.Println("len:", c-d, "start:", d, "end:", c)
	stringData := singleLine[d:c]
	floatData, err := strconv.ParseFloat(stringData, 64)
	if err != nil {
		fmt.Println("error with strconv.ParseFloat(stringData)", err)
	}
	return floatData
}
