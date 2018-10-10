package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var (
	a, b, c, d int
)

type pagerEvents struct {
	timeStamp        time.Time
	siteID           int
	pagerElapsedTime int
	pagerUUID        string
	orderType        string
	tableName        string
	pagerState       string
	pagerID          int
	pagerPaged       string
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
	file, err := os.Open("lrs_events.log")
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

			fmt.Println("-----------------------------------------------------------")
			fmt.Println("Line number: ", i)
			// fmt.Println(singleLine)

			// TimeStamp:
			timeStamp := extractStringData(singleLine, "<190>1 ", " host")

			// 2017-04-02T20:16:28.810638+00:00
			rfc3339Micro := "2006-01-02T15:04:05.999999-07:00"
			t1, err := time.Parse(rfc3339Micro, timeStamp)
			if err != nil {
				fmt.Println(err)
			}
			t2 := t1.Add(-4 * time.Hour)
			fmt.Println("timeStamp UTC:", t1)
			fmt.Println("timeStamp EST:", t2)

			// Site:
			siteID := extractIntData(singleLine, "Site: ", " {ElapsedTime:")
			fmt.Println("siteID: ", siteID)

			// ElapsedTime:
			pagerElapsedTime := extractIntData(singleLine, "ElapsedTime:", " UUID:")
			fmt.Println("elapsedTimeLRS: ", pagerElapsedTime)

			// UUID:
			pagerUUID := extractStringData(singleLine, "UUID:", " OrderType:")
			fmt.Println("UUID:", pagerUUID)

			// OrderType:
			orderType := extractStringData(singleLine, "OrderType:", " TableName:")
			fmt.Println("orderType:", orderType)

			// TableName:
			tableNumber := extractStringData(singleLine, "TableName:", " State:")
			fmt.Println("tableName:", tableNumber)

			// State:
			pagerState := extractStringData(singleLine, "State:", " PagerNumber:")
			fmt.Println("state:", pagerState)

			// PagerNumber:
			pagerNumber := extractIntData(singleLine, "PagerNumber:", " Paged:")
			fmt.Println("pagerID: ", pagerNumber)

			// Paged:
			pagerPaged := extractStringData(singleLine, "Paged:", "}")
			fmt.Println("paged:", pagerPaged)

			timeStart = time.Now()
			sqlStatement := "INSERT INTO pager_events (time_stamp, site_id, elapsed_time, UUID, order_type, table_number, pager_state, pager_number, pager_paged) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)RETURNING record_number"
			var recordNumber int
			err = db.QueryRow(sqlStatement, timeStamp, siteID, pagerElapsedTime, pagerUUID, orderType, tableNumber, pagerState, pagerNumber, pagerPaged).Scan(&recordNumber)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			fmt.Println("... responseTime - db.QueryRow(): ", time.Since(timeStart))

			i++
		}
	}
}

func extractStringData(singleLine string, startTag string, endTag string) string {

	a := len(startTag)
	b := strings.Index(singleLine, startTag)
	c := strings.Index(singleLine, endTag)
	d := b + a
	stringData := singleLine[d:c]

	return stringData
}

func extractIntData(singleLine string, startTag string, endTag string) int {

	a := len(startTag)
	b := strings.Index(singleLine, startTag)
	c := strings.Index(singleLine, endTag)
	d := b + a
	stringData := singleLine[d:c]

	intData, err := strconv.Atoi(stringData)
	if err != nil {
		fmt.Println(err)
	}
	return intData
}
