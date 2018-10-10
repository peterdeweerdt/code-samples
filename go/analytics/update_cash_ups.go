package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/peterdeweerdt/rize_analytics/core"
	"github.com/peterdeweerdt/rize_analytics/kounta"

	_ "github.com/lib/pq"
)

var (
	timeStart    time.Time
	sqlStatement string
)

func main() {

	fmt.Println("... processing:                           get all cash up")

	// validating connection arguments with postgres database
	timeStart = time.Now()
	fmt.Println("... checking database:                     validating the connection arguments")
	db, err := sql.Open("postgres", "user=postgres dbname=rize_analytics port=5432 password=April2016! sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("... responseTime - sql.Open:              ", time.Since(timeStart))

	// validating connection with postgres database
	timeStart = time.Now()
	fmt.Println("... checking database:                     validating the connection")
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("... responseTime - db.Ping:               ", time.Since(timeStart))

	requestURL := core.GetCompanyLevelURL(20246, "cashups")
	fmt.Println("... request URL                           ", requestURL)

	totalPages := core.GetTotalNumberOfPages(requestURL)
	fmt.Println("... request URL                           ", requestURL)
	fmt.Println("... total number of pages:                ", totalPages)

	countCashUps := 0
	i := 1
	for i < totalPages+1 {

		fmt.Println("... request URL                           ", requestURL)
		response := core.CreateGetRequest(requestURL)
		fmt.Println("... status from kounta:                   ", response.Status)

		timeStart = time.Now()
		requestURL = response.Header.Get("X-Next-Page")
		fmt.Println("... response.Header(X-Next-Page):         ", requestURL)
		lengthRequestURL := len(requestURL)
		fmt.Println("... length of request URL:                ", lengthRequestURL)
		fmt.Println("... response time - response.Header.Get   ", time.Since(timeStart))

		// proceed only if we receive a valid response from the API
		okToProceed := response.StatusCode == 200
		fmt.Println("... OK to Proceed:                        ", okToProceed)
		if okToProceed {

			// decode API json response using core.OrderDetail struct
			timeStart = time.Now()
			r := kounta.CashUps{}
			err := json.NewDecoder(response.Body).Decode(&r)
			if err == io.EOF {
				fmt.Println("... end of file: continue")
			}
			if err != nil {
				fmt.Println("... failure: ", err)
			}
			fmt.Println("... response time - json.NewDecoder       ", time.Since(timeStart))

			totalCashUpPerPage := len(r)
			fmt.Println("... total number of cashups per page is:  ", totalCashUpPerPage)

			j := 0
			for j < totalCashUpPerPage {
				fmt.Println("... working on page number:   ", i, "out of ", totalPages)
				fmt.Println("... working on line number:   ", j+1, "out of ", totalCashUpPerPage)
				percentComplete := (float64(countCashUps) / float64(totalCashUpPerPage*totalPages)) * 100.0
				fmt.Printf("... percent complete:     %9.1f%%\n", percentComplete)

				// assign decoded json response to variables
				cashUpID := r[j].ID
				siteID := r[j].SiteID
				processed := r[j].Processed
				createdAt := r[j].CreatedAt

				// print variables that will be loaded into the database to the screen
				fmt.Println("-----------------------------------------------------------------------")
				fmt.Println("... cash_up_id:               ", cashUpID)
				fmt.Println("... site_id:                  ", siteID)
				fmt.Println("... processed:                ", processed)
				fmt.Println("... createdAt:                ", createdAt)

				// check if the category is already in the database
				timeStart = time.Now()
				sqlStatement = "SELECT cash_up_id FROM cash_ups WHERE cash_up_id=$1"
				row := db.QueryRow(sqlStatement, cashUpID)
				if err != nil {
					fmt.Println("... error with db.QueryRow    ", err)
				}
				fmt.Println("... responseTime - db.QueryRow", time.Since(timeStart))

				timeStart = time.Now()
				err := row.Scan(&cashUpID)
				if err != nil {
					fmt.Println("... error with row.Scan       ", err)
				}
				fmt.Println("... responseTime - row.Scan:  ", time.Since(timeStart))

				switch err {

				case sql.ErrNoRows:
					// if category is not in the database, insert a new record
					timeStart = time.Now()
					fmt.Println("... database action:           adding a new category to database")
					sqlStatement = "INSERT INTO cash_ups (cash_up_id, site_id, processed, created_at) VALUES ($1, $2, $3, $4)RETURNING cash_up_id"
					err = db.QueryRow(sqlStatement, cashUpID, siteID, processed, createdAt).Scan(&cashUpID)
					if err != nil {
						fmt.Println("... error with db.QueryRow    ", err)
						panic(err)
					}
					fmt.Println("... responseTime - db.QueryRow ", time.Since(timeStart))

				case nil:
					// if category ID is already in the database, update the record
					timeStart = time.Now()
					fmt.Println("... database action:           updating order")
					sqlStatement = "UPDATE cash_ups SET site_id=$2, processed=$3, created_at=$4 WHERE cash_up_id=$1"
					_, err = db.Exec(sqlStatement, cashUpID, siteID, processed, createdAt)
					if err != nil {
						fmt.Println("... error with db.Exec        ", err)
						panic(err)
					}
					fmt.Println("... responseTime - db.Exec:   ", time.Since(timeStart))

				default:
					fmt.Println("Error with db.QueryRow  ")
					panic(err)
				}
				countCashUps++
				j++
			}
		}
		i++

	}

}
