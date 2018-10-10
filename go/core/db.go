package core

import "time"

type DatabaseID int64

type DB interface {
	InsertTableMap(tableMap *TableMap) error
	GetTableMapByBeaconID(id string) (*TableMap, error)

	UpdateCayanKey(token string) error
	GetCayanKey() (*Key, error)

	InsertToken(token *Token) error
	GetToken(tokenString string) (*Token, error)
	DeleteToken(id DatabaseID) error
	DeleteTokens(customerID DatabaseID) error

	InsertCustomer(customer *Customer) error
	UpdateCustomerPassword(id DatabaseID, passwordHash string) (*Customer, error)
	GetCustomer(id DatabaseID) (*Customer, error)
	GetCustomerByExternalID(id string) (*Customer, error)
	GetCustomerByEmail(email string) (*Customer, error)

	InsertOrder(order *Order) error
	UpdateOrder(order *Order) error
	UpdateOrderTableName(order *Order, tableName string) error
	UpdateOrderCustomerID(order *Order, customerID DatabaseID) error
	UpdateOrderPickupTime(order *Order, pickupTime time.Time) error
	GetOrder(orderID KountaID) (*Order, error)
	GetOrderByDatabaseID(orderID DatabaseID) (*Order, error)
	GetOrderByPagerID(siteID KountaID, pagerID int64) (*Order, error)
	SelectOrdersByCustomerID(customerID DatabaseID) (*[]Order, error)

	// SelectOnHoldAndPendingOrdersByTable will return all orders that are either 'on hold' or 'pending' for a given table
	SelectOnHoldAndPendingOrdersByTable(siteID KountaID, tableName string) (*[]Order, error)
	// SelectOnHoldAndPendingOrdersByPagerID will return all orders that are either 'on hold' or 'pending' for a pager ID
	SelectOnHoldAndPendingOrdersByPagerID(siteID KountaID, pagerID int64) (*[]Order, error)
	InsertOrderUpdate(orderUpdate KountaOrderUpdate) error
	GetLine(lineID DatabaseID) (*Line, error)
	SelectLines(orderID DatabaseID) (*[]Line, error)
	SelectAddedModifiers(lineID DatabaseID) (*[]Modifier, error)
	SelectRemovedModifiers(lineID DatabaseID) (*[]Modifier, error)

	InsertPayment(payment *Payment, order *Order) error
	// GetPaymentByOrderID will get a payment for a given order ID
	GetPaymentByOrderID(id DatabaseID) (*Payment, error)

	InsertSite(site *Site) error
	UpdateSiteMenuHash(site *Site, menuHash string) error
	SelectSites() (*[]Site, error)
	GetSite(id KountaID) (*Site, error)

	// UpsertCategory will either update or insert core.Category into database
	UpsertCategory(category *Category) error
	SelectCategories() (*[]Category, error)
	SelectCategoriesBySiteID(siteID KountaID) (*[]Category, error)
	DeleteCategory(categoryID KountaID) error

	UpsertMenuItem(item *MenuItem) error
	SelectMenuItems() (*[]MenuItem, error)
	SelectMenuItemsByCategoryID(siteID KountaID, categoryID DatabaseID) (*[]MenuItem, error)
	GetMenuItem(siteID KountaID, menuItemID DatabaseID) (*MenuItem, error)
	DeleteMenuItem(menuItemID KountaID) error

	UpsertMenuItemModifier(item *MenuItem, modifier *Modifier) error
	UpsertOptionSetModifier(optionSet *OptionSet, modifier *Modifier) error
	GetMenuModifier(siteID KountaID, modifierID DatabaseID) (*Modifier, error)
	GetMenuModifierByKountaID(siteID, modifierID KountaID) (*Modifier, error)
	SelectMenuModifiers() (*[]Modifier, error)
	SelectMenuItemModifiers(siteID KountaID, menuItemID DatabaseID) (*[]Modifier, error)
	DeleteMenuModifier(modifierID KountaID) error

	UpsertOptionSet(item *MenuItem, optionSet *OptionSet) error
	GetOptionSet(optionSetID DatabaseID) (*OptionSet, error)
	SelectOptionSets() (*[]OptionSet, error)
	SelectOptionSetsByItemID(menuItemID DatabaseID) (*[]OptionSet, error)
	DeleteOptionSet(optionSetID KountaID) error
}
