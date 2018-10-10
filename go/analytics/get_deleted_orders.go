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

func main() {
	companyID := 20246
	startDate := "2017-06-28"
	fmt.Println("... process:                  getting deleted orders")
	fmt.Println("... companyID:                ", companyID)
	fmt.Println("... startDate:                ", startDate)

	// validate the connection arguments to rize analytics postgres database
	var timeStart time.Time
	timeStart = time.Now()
	fmt.Println("... checking database:        validating the connection arguments")

	db, err := sql.Open("postgres", core.PostgresInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("... responseTime - sql.Open:  ", time.Since(timeStart))

	// validate the connection to the rize analytics postgress database
	timeStart = time.Now()
	fmt.Println("... checking database:        validating the connection")
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("... responseTime - db.Ping:   ", time.Since(timeStart))

	var requestURL string
	fmt.Println("... calling function:         core.GetAllOrdersByDataURL")
	requestURL = core.GetDeletedOrdersURL(companyID)

	fmt.Println("... calling function:         core.GetTotalNumberOfPages")
	totalPages := core.GetTotalNumberOfPages(requestURL)
	totalOrders := totalPages * 25
	fmt.Println("... total number of pages:    ", totalPages)
	fmt.Println("... total number of orders:   ", totalOrders)

	countOrders := 0
	countPages := 0
	for countPages < totalPages {

		fmt.Println("... calling function:         core.CreateGetRequest")
		response := core.CreateGetRequest(requestURL)

		fmt.Println("... response.StatusCode:      ", response.StatusCode)

		okToProceed := response.StatusCode == 200
		if okToProceed {

			// unmarshal json data from API response
			timeStart := time.Now()
			r := kounta.Order{}
			err := json.NewDecoder(response.Body).Decode(&r)
			if err == io.EOF {
				fmt.Println("... error:                    end of file error with json.NewDecoder()")
				continue
			}
			if err != nil {
				fmt.Println("... error:                    error with json.NewDecoder(): ", err)
				continue
			}
			fmt.Println("... responseTime -            ", time.Since(timeStart))

			totalOrdersPerPage := len(r)
			totalOrders := totalPages * totalOrdersPerPage

			fmt.Println("... total # of orders / page: ", totalOrdersPerPage)

			countOrdersPerPage := 0
			fmt.Println("... countOrdersPerPage        ", countOrdersPerPage+1)
			for countOrdersPerPage < totalOrdersPerPage {

				fmt.Println("-----------------------------------------------------------------------")

				// calcuate percent complete as a percent of total orders to be processed
				percentCompleteOrders := (float64(countOrders) / float64(totalOrders)) * 100.0
				fmt.Println("... orders:                   ", countOrders, "out of ", totalOrders)
				fmt.Printf("... percent complete:         %9.1f%% \n", percentCompleteOrders)

				// calcualte percent complete as a percent of total API calls
				fmt.Println("... pages:                                ", countPages, " out of", totalPages)
				percentCompletePages := (float64(countPages) / float64(totalPages)) * 100.0
				fmt.Printf("... percent complete:                 %9.1f%% \n", percentCompletePages)

				// assign decoded json response to variables
				orderID := r[countOrdersPerPage].OrderID
				saleNumber := r[countOrdersPerPage].SaleNumber
				status := r[countOrdersPerPage].Status
				siteID := r[countOrdersPerPage].SiteID
				registerID := r[countOrdersPerPage].RegisterID
				staffID := r[countOrdersPerPage].StaffID
				total := r[countOrdersPerPage].Total
				totalTax := r[countOrdersPerPage].TotalTax
				paid := r[countOrdersPerPage].Paid
				deleted := r[countOrdersPerPage].Deleted
				createdAtString := r[countOrdersPerPage].CreatedAt
				updatedAtString := r[countOrdersPerPage].UpdatedAt
				createdAt := core.ExtractTimeStamp(createdAtString)
				updatedAt := core.ExtractTimeStamp(updatedAtString)
				weekDay := updatedAt.Weekday()
				hour := updatedAt.Hour()
				date := core.ExtractDate(updatedAtString)
				time := core.ExtractTime(updatedAtString)
				mealPeriod := core.ExtractMealPeriod(updatedAtString)

				// print variables that will be loaded into the database
				fmt.Println("-----------------------------------------------------------------------")
				fmt.Println("... pages:                    ", countOrdersPerPage, "out of ", totalOrdersPerPage)
				fmt.Println("... order_id:                 ", orderID)
				fmt.Println("... sale_number:              ", saleNumber)
				fmt.Println("... status:                   ", status)
				fmt.Println("... site_id:                  ", siteID)
				fmt.Println("... register_id:              ", registerID)
				fmt.Println("... staff_id:                 ", staffID)
				fmt.Println("... total:                    ", total)
				fmt.Println("... total_tax:                ", totalTax)
				fmt.Println("... paid:                     ", paid)
				fmt.Println("... deleted:                  ", deleted)
				fmt.Println("... created_at:               ", createdAt)
				fmt.Println("... updated_at:               ", updatedAt)
				fmt.Println("... week_day:                 ", weekDay)
				fmt.Println("... hour:                     ", hour)
				fmt.Println("... date:                     ", date)
				fmt.Println("... time:                     ", time)
				fmt.Println("... mealPeriod:               ", mealPeriod)
				fmt.Println("-----------------------------------------------------------------------")

				// check if order ID is already in the database
				var sqlStatement string
				sqlStatement = "SELECT order_id FROM deleted_orders WHERE order_id=$1"
				row := db.QueryRow(sqlStatement, orderID)
				if err != nil {
					fmt.Println("... error with db.QueryRow    ", err)
				}
				err := row.Scan(&orderID)
				if err != nil {
					fmt.Println("... error with row.Scan        ", err)
				}

				switch err {
				// if order ID is not in the database, insert a new record
				case sql.ErrNoRows:
					fmt.Println("... database action:          adding order")
					sqlStatement = "INSERT INTO deleted_orders (order_id, sale_number, status,	site_id, register_id, staff_id, total, total_tax, paid, deleted, created_at, updated_at, week_day, hour, date, time, meal_period) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING order_id"
					err = db.QueryRow(sqlStatement, orderID, saleNumber, status, siteID, registerID, staffID, total, totalTax, paid, deleted, createdAt, updatedAt, weekDay, hour, date, time, mealPeriod).Scan(&orderID)
					if err != nil {
						fmt.Println("... error with db.QueryRow:   ", err)
						panic(err)
					}
				// if order ID is already in the database, update the record
				case nil:
					fmt.Println("... database action:          updating order")
					sqlStatement = "UPDATE deleted_orders SET sale_number=$2, status=$3, site_id=$4, register_id=$5, staff_id=$6, total=$7, total_tax=$8, paid=$9, deleted=$10, created_at=$11, updated_at=$12, week_day=$13, hour=$14, date=$15, time=$16, meal_period=$17 WHERE order_id=$1"
					_, err = db.Exec(sqlStatement, orderID, saleNumber, status, siteID, registerID, staffID, total, totalTax, paid, deleted, createdAt, updatedAt, weekDay, hour, date, time, mealPeriod)
					if err != nil {
						fmt.Println("... error with db.Exec        ", err)
						panic(err)
					}

				default:
					fmt.Println("Error with db.QueryRow()")
					panic(err)
				}
				// check if order ID is already in the database

				countOrders++
				fmt.Println("... countOrders:              ", countOrders)

				countOrdersPerPage++
				fmt.Println("... countOrdersPerPage:       ", countOrdersPerPage)

			}
			requestURL := core.GetNextPageURL(&requestURL)
			fmt.Println("... requestURL:               ", requestURL)

			countPages++
			fmt.Println("... countPages:               ", countPages)
		}
	}
}
