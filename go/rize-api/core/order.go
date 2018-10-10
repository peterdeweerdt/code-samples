package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// These are used to indicate the current status of an order (as strings)
const (
	OrderStatusSubmitted = "SUBMITTED"
	OrderStatusAccepted  = "ACCEPTED"
	OrderStatusRejected  = "REJECTED"
	OrderStatusOnHold    = "ON_HOLD"
	OrderStatusComplete  = "COMPLETE"
	OrderStatusPending   = "PENDING"
	OrderStatusDeleted   = "DELETED"
)

// Order represents a Rize order. This is independent of any other backend the
// application uses.
type Order struct {
	ID          DatabaseID    `json:"id"`
	PosID       KountaID      `json:"-"`
	Status      string        `json:"status"`
	TableName   string        `json:"table_name"`
	CustomerID  sql.NullInt64 `json:"customer_id"`
	Total       int           `json:"total"`
	TotalTax    int           `json:"total_tax"`
	Lines       []Line        `json:"lines"`
	PagerNumber string        `json:"puck_id"` //todo: coordinate the rename with the client apps
	SiteID      KountaID      `json:"site_id"`
	CreatedAt   time.Time     `json:"-"`
	PickupTime  *time.Time    `json:"-"`
}

// METHODS

// MarshalJSON provides a custom override for handling a DB nullable customer_id
func (o *Order) MarshalJSON() ([]byte, error) {
	type alias Order

	var id *int64
	if o.CustomerID.Valid {
		id = &o.CustomerID.Int64
	}

	return json.Marshal(&struct {
		*alias
		CustomerID *int64 `json:"customer_id"`
	}{
		alias:      (*alias)(o),
		CustomerID: id,
	})
}

// CreateOrUpdateOrderFromKounta can be called from an external system to update the associated rize.Order
func (app AppContext) CreateOrUpdateOrderFromKounta(kountaOrder KountaOrder) (*Order, error) {
	existingOrder, err := app.DB.GetOrder(kountaOrder.GetPosID())
	if err != nil {
		return nil, errors.Wrapf(err, "create or update order from kounta")
	}

	if existingOrder != nil {
		order := newOrderFromKountaOrder(kountaOrder)
		err := app.UpdateOrder(order)
		if err != nil {
			return nil, errors.Wrap(err, "create or update order from kounta")
		}
		return order, nil
	}

	return app.createOrderFromKountaOrder(kountaOrder)
}

func (app *AppContext) CreateOrderForPager(siteID KountaID, pagerNumber int64) (*Order, error) {
	existingOrder, err := app.DB.GetOrderByPagerID(siteID, pagerNumber)
	if err != nil {
		return nil, errors.Wrapf(err, "order: error getting order for site '%d' and pager '%d'", pagerNumber)
	}
	if existingOrder != nil {
		return nil, errors.Errorf("order: existing order '%d' found for pager '%d' at site '%d'", existingOrder.ID, pagerNumber, siteID)
	}

	kountaOrder, err := app.Kounta.CreateOrderForPager(siteID, pagerNumber)
	if err != nil {
		return nil, err
	}

	createdOrder, err := app.createOrderFromKountaOrder(kountaOrder)
	if err != nil {
		return nil, err
	}

	return createdOrder, nil
}

// CreateNewOrder will create a new Kounta order with menu items and save in database
func (app *AppContext) CreateNewOrder(siteID KountaID, createOrder CreateOrder) (*Order, error) {
	if err := app.addKountaIDsToNewOrder(siteID, &createOrder); err != nil {
		return nil, errors.Wrap(err, "create new order")
	}

	kountaOrder, err := app.Kounta.CreateOrder(siteID, createOrder)
	if err != nil {
		return nil, errors.Wrap(err, "create new order")
	}

	createdOrder, err := app.createOrderFromKountaOrder(kountaOrder)
	if err != nil {
		return nil, errors.Wrap(err, "create new order")
	}

	return createdOrder, nil
}

// AddMenuItemsToOrder will add a a list of new menu items to an existing order
func (app *AppContext) AddMenuItemsToOrder(orderID DatabaseID, menuItems []CreateOrderMenuItem) (*Order, error) {
	// find existing order in db
	order, err := app.FindOrderByID(DatabaseID(orderID))
	if err != nil {
		return nil, errors.Wrap(err, "add menu items to order")
	}
	if order == nil {
		return nil, errors.New(fmt.Sprintf("add menu items to order: order %d not found", orderID))
	}

	for i := range menuItems {
		if err = app.addKountaIDsToNewMenuItem(order.SiteID, &menuItems[i]); err != nil {
			return nil, errors.Wrap(err, "add menu items to order")
		}
	}

	updatedOrder, err := app.Kounta.AddMenuItemsToOrder(order.PosID, menuItems)
	if err != nil {
		return nil, errors.Wrap(err, "add menu items to order")
	}

	order, err = app.CreateOrUpdateOrderFromKounta(updatedOrder)
	if err != nil {
		return nil, errors.Wrap(err, "add menu items to order")
	}

	return order, nil
}

func (app AppContext) createOrderFromKountaOrder(kountaOrder KountaOrder) (*Order, error) {
	order := newOrderFromKountaOrder(kountaOrder)

	if err := app.DB.InsertOrder(order); err != nil {
		return nil, errors.Wrapf(err, "createOrderFromKountaOrder(%d)", kountaOrder.GetPosID())
	}

	return order, nil
}

func newOrderFromKountaOrder(kountaOrder KountaOrder) *Order {
	order := &Order{
		PosID:       kountaOrder.GetPosID(),
		Status:      kountaOrder.GetStatus(),
		TableName:   kountaOrder.GetTable(),
		Total:       kountaOrder.GetTotal(),
		TotalTax:    kountaOrder.GetTotalTax(),
		Lines:       kountaOrder.GetLines(),
		PagerNumber: kountaOrder.GetPagerNumber(),
		SiteID:      kountaOrder.GetSiteID(),
	}

	return order
}

// FindOrderByPosID will return a core.Order from the database whose PosID matches id
func (app AppContext) FindOrderByPosID(id KountaID) (*Order, error) {
	order, err := app.DB.GetOrder(id)
	if err != nil {
		return nil, errors.Wrap(err, "find order by POS id")
	}

	error := app.loadLines(order)
	if error != nil {
		return nil, errors.Wrap(error, "find order by POS id")
	}

	return order, nil
}

func (app AppContext) FindOrderByID(id DatabaseID) (*Order, error) {
	order, err := app.DB.GetOrderByDatabaseID(id)
	if err != nil {
		return nil, errors.Wrapf(err, "find order by id %d", id)
	}
	if order == nil {
		return nil, nil
	}

	app.loadLines(order)

	return order, nil
}

// FindOrdersByTableName will find all "payable" orders for a given table name
func (app AppContext) FindPayableOrdersByTableName(siteID KountaID, tableName string) ([]Order, error) {
	orders, err := app.DB.SelectOnHoldAndPendingOrdersByTable(siteID, tableName)
	if err != nil {
		return nil, errors.Wrapf(err, "find payable orders by table %s", tableName)
	}

	payableOrders, err := app.getPayableOrders(*orders)
	if err != nil {
		return nil, errors.Wrapf(err, "find payable orders by table %s", tableName)
	}

	for i := range payableOrders {
		app.loadLines(&(payableOrders)[i])
	}

	return payableOrders, nil
}

// FindOrdersByPagerID will find all "payable" orders for a given pager ID
func (app AppContext) FindPayableOrdersByPagerID(siteID KountaID, pagerID int64) ([]Order, error) {
	orders, err := app.DB.SelectOnHoldAndPendingOrdersByPagerID(siteID, pagerID)
	if err != nil {
		return nil, errors.Wrapf(err, "find payable orders by pager %d", pagerID)
	}

	payableOrders, err := app.getPayableOrders(*orders)
	if err != nil {
		return nil, errors.Wrapf(err, "find payable orders by pager %d", pagerID)
	}

	for i := range payableOrders {
		err = app.loadLines(&(payableOrders)[i])
		if err != nil {
			return nil, errors.Wrapf(err, "find payable orders by pager %d", pagerID)
		}
	}

	return payableOrders, nil
}

// getPayableOrders will return a slice of only the payable orders from the given slice
func (app AppContext) getPayableOrders(orders []Order) ([]Order, error) {
	// filter out any paid orders
	payableOrders := []Order{}
	for _, order := range orders {

		payable, err := app.IsOrderPayable(order)
		if err != nil {
			return nil, errors.Wrapf(err, "filter out payable orders")
		}
		if payable {
			payableOrders = append(payableOrders, order)
		}
	}

	return payableOrders, nil
}

// IsOrderPayable will return boolean on whether order can be paid
func (app AppContext) IsOrderPayable(order Order) (bool, error) {
	payment, err := app.DB.GetPaymentByOrderID(order.ID)
	if err != nil {
		return false, errors.Wrapf(err, "filter out payable orders")
	}
	if payment != nil {
		return false, nil
	}

	return order.Status == OrderStatusPending || order.Status == OrderStatusOnHold, nil
}

// TODO: Move this into database layer
// loadLines will find and attach lines to this order
func (app AppContext) loadLines(o *Order) error {
	lines, err := app.DB.SelectLines(o.ID)
	if err != nil {
		return errors.Wrapf(err, "Error loading lines for Order '%v'", o.ID)
	}
	o.Lines = *lines

	for i, l := range o.Lines {
		added, err := app.DB.SelectAddedModifiers(l.ID)
		if err != nil {
			return errors.Wrapf(err, "Error loading added modifiers for Line '%v'", l.ID)
		}
		if len(*added) > 0 {
			o.Lines[i].AddedModifiers = *added
		}

		removed, err := app.DB.SelectRemovedModifiers(l.ID)
		if err != nil {
			return errors.Wrapf(err, "Error loading removed modifiers for Line '%v'", l.ID)
		}
		if len(*removed) > 0 {
			o.Lines[i].RemovedModifiers = *removed
		}
	}
	return nil
}

func (app AppContext) DeleteLine(orderID, lineID DatabaseID) error {
	order, err := app.DB.GetOrderByDatabaseID(orderID)
	if err != nil {
		return errors.Wrap(err, "delete line")
	}
	if order == nil {
		return errors.New("delete line: order not found")
	}
	if order.Status != OrderStatusSubmitted {
		return errors.New("delete line: order not submitted status")
	}
	line, err := app.DB.GetLine(lineID)
	if err != nil {
		return errors.Wrap(err, "delete line")
	}
	if line == nil {
		return errors.New("delete line: line not found")
	}

	if err := app.Kounta.DeleteLineItem(order.PosID, line.PosID); err != nil {
		return errors.Wrap(err, "delete line")
	}

	// get updated order from kounta
	kountaOrder, err := app.Kounta.GetOrderByID(order.PosID)
	if err != nil {
		return errors.Wrap(err, "delete line")
	}
	_, err = app.CreateOrUpdateOrderFromKounta(kountaOrder)
	if err != nil {
		return errors.Wrap(err, "delete line")
	}

	return nil
}

func (app AppContext) LinkOrderWithTable(siteID KountaID, pagerNumber int64, tableName string) error {
	order, err := app.DB.GetOrderByPagerID(siteID, pagerNumber)
	if err != nil {
		return errors.Wrapf(err, "order: error linking order with site '%d' and pager '%s' with table '%s'", siteID, pagerNumber, tableName)
	}
	if order == nil {
		return errors.New(fmt.Sprintf("link order with table: no orders for pager %d", pagerNumber))
	}

	err = app.Kounta.LinkOrderWithTable(order.PosID, tableName)
	if err != nil {
		return errors.Wrapf(err, "order: error linking order '%d' with table '%s' in pos", siteID, tableName)
	}

	err = app.DB.UpdateOrderTableName(order, tableName)
	if err != nil {
		return errors.Wrap(err, "order: error setting table name on order")
	}

	return err
}

func (app AppContext) UpdateOrder(order *Order) error {
	if order.Status == OrderStatusComplete || order.Status == OrderStatusRejected || order.Status == OrderStatusDeleted {
		order.PagerNumber = "" // clear pager number for paid orders so that pager can be reused
	}
	return app.DB.UpdateOrder(order)
}

func (app AppContext) UpdateOrderWithCustomer(orderID, customerID DatabaseID) error {
	order, err := app.DB.GetOrderByDatabaseID(orderID)
	if err != nil {
		return errors.Wrapf(err, "UpdateOrderWithCustomer(%d, %d)", orderID, customerID)
	}
	if order == nil {
		return errors.New(fmt.Sprintf("update order: no order for id %d", orderID))
	}

	customer, err := app.DB.GetCustomer(customerID)
	if err != nil {
		return errors.Wrapf(err, "UpdateOrderWithCustomer(%d, %d)", orderID, customerID)
	}

	err = app.Kounta.AddCustomerToOrder(order.PosID, customer.PosID)
	if err != nil {
		return errors.Wrapf(err, "UpdateOrderWithCustomer(%d, %d)", orderID, customerID)
	}

	err = app.DB.UpdateOrderCustomerID(order, customerID)
	if err != nil {
		return errors.Wrapf(err, "UpdateOrderWithCustomer(%d, %d)", orderID, customerID)
	}

	return nil
}

func (app AppContext) CompleteOrder(order *Order) error {
	order.Status = OrderStatusComplete

	err := app.Kounta.CompleteOrder(order.PosID)
	if err != nil {
		return errors.Wrapf(err, "CompleteOrder(%d)", order.ID)
	}

	err = app.UpdateOrder(order)
	if err != nil {
		return errors.Wrapf(err, "CompleteOrder(%d)", order.ID)
	}

	return err
}

// RejectOrder will mark an order as rejected in kounta and database
func (app AppContext) RejectOrder(orderID KountaID) error {
	kountaOrder, err := app.Kounta.RejectOrder(orderID)
	if err != nil {
		return errors.Wrap(err, "reject order")
	}

	_, err = app.CreateOrUpdateOrderFromKounta(kountaOrder)
	if err != nil {
		return errors.Wrap(err, "reject order")
	}

	return nil
}

// UpdatePickupDetails will update any pickup details on an order, or add new details if none set, and move to ON HOLD
func (app AppContext) UpdatePickupDetails(orderID DatabaseID, pickupDetails PickupDetails) error {
	order, err := app.DB.GetOrderByDatabaseID(orderID)
	if err != nil {
		return errors.Wrap(err, "update pickup details")
	}
	if order == nil {
		return errors.New("update pickup details: order not found")
	}

	pickupTime := pickupDetails.PickupTime
	if err != nil {
		errors.Wrap(err, "update pickup details")
	}
	const layout = "3:04 pm"
	pickupTimeReadableString := pickupTime.Format(layout)
	notes := fmt.Sprintf("TO-GO APP - PAID\n\n%s\n\n%s\n\n%s", pickupTimeReadableString, pickupDetails.CustomerName, pickupDetails.PhoneNumber)

	err = app.Kounta.SetOrderNotes(order.PosID, notes)
	if err != nil {
		return errors.Wrap(err, "update pickup details")
	}

	err = app.Kounta.PutOrderOnHold(order.PosID)
	if err != nil {
		return errors.Wrap(err, "update pickup details")
	}

	// turn around and get order from kounta now that we've updated
	kountaOrder, err := app.Kounta.GetOrderByID(order.PosID)
	if err != nil {
		return errors.Wrap(err, "update pickup details")
	}

	_, err = app.CreateOrUpdateOrderFromKounta(kountaOrder)
	if err != nil {
		return errors.Wrap(err, "update pickup details")
	}

	if err = app.DB.UpdateOrderPickupTime(order, *pickupTime); err != nil {
		return errors.Wrap(err, "update pickup details")
	}

	return nil
}

func (app *AppContext) addKountaIDsToNewOrder(siteID KountaID, createOrder *CreateOrder) error {
	for i := range createOrder.MenuItems {
		if err := app.addKountaIDsToNewMenuItem(siteID, &createOrder.MenuItems[i]); err != nil {
			return errors.Wrap(err, "add kounta ids to new order")
		}
	}
	return nil
}

func (app *AppContext) addKountaIDsToNewMenuItem(siteID KountaID, menuItem *CreateOrderMenuItem) error {
	existingMenuItem, err := app.DB.GetMenuItem(siteID, menuItem.ID)
	if err != nil {
		return errors.Wrap(err, "add kounta ids to new menu item")
	}
	if existingMenuItem == nil {
		return errors.New("add kounta ids to new menu item")
	}

	menuItem.PosID = existingMenuItem.PosID
	kountaModifiers := make([]KountaID, len(menuItem.SelectedModifierIDs))
	for i, modifierID := range menuItem.SelectedModifierIDs {
		existingModifier, err := app.DB.GetMenuModifier(siteID, modifierID)
		if err != nil {
			return errors.Wrap(err, "add kounta ids to new menu item")
		}
		if existingModifier == nil {
			return errors.New(fmt.Sprintf("add kounta ids to new menu item: modifier %d not found", modifierID))
		}
		kountaModifiers[i] = existingModifier.PosID
	}
	menuItem.PosModifierIDs = kountaModifiers

	kountaOptions := make([]MenuItemKountaOption, len(menuItem.SelectedOptions))
	for i, option := range menuItem.SelectedOptions {
		optionSet, err := app.DB.GetOptionSet(option.OptionSetID)
		if err != nil {
			return errors.Wrap(err, "add kounta ids to new menu item")
		}
		if optionSet == nil {
			return errors.New(fmt.Sprintf("add kounta ids to new menu item: option set %d not found", option.OptionSetID))
		}

		modifier, err := app.DB.GetMenuModifier(siteID, option.ModifierID)
		if err != nil {
			return errors.Wrap(err, "add kounta ids to new menu item")
		}
		if modifier == nil {
			return errors.New(fmt.Sprintf("add kounta ids to new menu item: modifier %d not found", option.ModifierID))
		}

		kountaOptions[i] = MenuItemKountaOption{
			OptionSetID: optionSet.PosID,
			ModifierID:  modifier.PosID,
		}
	}
	menuItem.PosSelectedOptions = kountaOptions

	return nil
}

// AppendUniqueOrders will combine initialOrders with any orders in newOrders that are NOT already in initialOrders.
// It will NOT add duplicates from newOrders.
func (app *AppContext) AppendUniqueOrders(initialOrders []Order, newOrders []Order) []Order {

	var toAdd []Order

	for _, newOrder := range newOrders {
		found := false
		for _, existingOrder := range initialOrders {
			if existingOrder.ID == newOrder.ID {
				found = true
			}
		}
		if !found {
			// make sure this new order isn't already in list of orders to add
			alreadyAdded := false
			for _, addOrder := range toAdd {
				if addOrder.ID == newOrder.ID {
					alreadyAdded = true
				}
			}
			if !alreadyAdded {
				toAdd = append(toAdd, newOrder)
			}
		}
	}

	initialOrders = append(initialOrders, toAdd...)
	return initialOrders
}
