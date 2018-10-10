package core

import (
	"database/sql"
	"time"
)

// Payment represents the metadata regarding a payment for an order
type Payment struct {
	Amount        int
	Tip           int
	OrderID       DatabaseID
	TransactionID string
	Date          time.Time
	CustomerID    sql.NullInt64
}
