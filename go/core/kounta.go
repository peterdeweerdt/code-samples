package core

import "time"

type KountaID int64

type Kounta interface {
	CreateOrder(siteID KountaID, newOrder CreateOrder) (KountaOrder, error)
	CreateOrderForPager(siteID KountaID, pagerNumber int64) (KountaOrder, error)
	GetOrderByID(posOrderID KountaID) (KountaOrder, error)
	AddMenuItemsToOrder(orderID KountaID, menuItems []CreateOrderMenuItem) (KountaOrder, error)
	LinkOrderWithTable(orderID KountaID, tableName string) error
	// SetOrderNotes will set the notes field on an order
	SetOrderNotes(orderID KountaID, notes string) error
	RejectOrder(orderID KountaID) (KountaOrder, error)
	// PutOrderOnHold will move Order to ON_HOLD status
	PutOrderOnHold(posOrderID KountaID) error
	// CompleteOrder will mark the Order complete, preventing any further modifications
	CompleteOrder(posOrderID KountaID) error
	CompleteAllPendingOrders(siteID KountaID) error
	ParseOrder(buffer []byte) (KountaOrder, error)
	ParseOrderUpdate(buffer []byte) (KountaOrderUpdate, error)

	RecordPayment(payment Payment, posOrderID KountaID) error

	CreateCustomer(email, firstName, lastName, phone string, rizeID DatabaseID) (KountaCustomer, error)
	GetCustomerByEmail(email string) (KountaCustomer, error)
	AddCustomerToOrder(posOrderID, customerID KountaID) error

	GetAllSites() ([]Site, error)
	GetMenuForSite(siteID KountaID) (*Menu, error)
	DeleteLineItem(orderID KountaID, lineID KountaID) error
}

type KountaOrder interface {
	GetPosID() KountaID
	GetStatus() string
	GetTable() string
	GetTotal() int
	GetTotalTax() int
	GetLines() []Line
	GetSiteID() KountaID
	GetPagerNumber() string
	GetNotes() string
}

type KountaOrderUpdate interface {
	GetOrderID() KountaID
	GetSaleNumber() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetDeleted() bool
	GetStatus() string
	GetNotes() string
	GetTotal() float64
	GetPaid() float64
	GetTips() float64
	GetRegisterID() KountaID
	GetSiteID() KountaID
	GetLines() []map[string]interface{}
	GetPriceVariation() float64
	GetPayments() []map[string]interface{}
	GetLock() []string
	GetStaffMemberID() KountaID
}

type KountaCustomer interface {
	GetPosID() KountaID
}
