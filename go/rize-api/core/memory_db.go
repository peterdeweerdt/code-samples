package core

import (
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type MemoryDB struct {
	Error           error
	TableMaps       map[string]TableMap
	CayanKeyVersion int
	CayanKey        *Key
	Tokens          map[DatabaseID]Token
	Customers       map[DatabaseID]Customer
	Orders          map[DatabaseID]Order
	LineCount       int
	Payments        map[string]Payment
	Sites           map[DatabaseID]Site
	Categories      map[DatabaseID]Category
	MenuItems       map[DatabaseID]MenuItem
	Modifiers       map[DatabaseID]Modifier
	OptionSets      map[DatabaseID]OptionSet
}

func (db *MemoryDB) Init() {
	db.Error = nil
	db.TableMaps = map[string]TableMap{}
	db.CayanKeyVersion = 0
	db.CayanKey = nil
	db.Tokens = map[DatabaseID]Token{}
	db.Customers = map[DatabaseID]Customer{}
	db.Orders = map[DatabaseID]Order{}
	db.LineCount = 0
	db.Payments = map[string]Payment{}
	db.Sites = map[DatabaseID]Site{}
	db.Categories = map[DatabaseID]Category{}
	db.MenuItems = map[DatabaseID]MenuItem{}
	db.Modifiers = map[DatabaseID]Modifier{}
	db.OptionSets = map[DatabaseID]OptionSet{}
}

type databaseIDSlice []DatabaseID

func (a databaseIDSlice) Len() int           { return len(a) }
func (a databaseIDSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a databaseIDSlice) Less(i, j int) bool { return a[i] < a[j] }

// Implement DB interface

func (db *MemoryDB) InsertTableMap(tableMap *TableMap) error {
	if db.Error != nil {
		return db.Error
	}

	tableMap.BeaconID = strings.ToLower(tableMap.BeaconID)
	db.TableMaps[tableMap.BeaconID] = *tableMap
	return nil
}

func (db *MemoryDB) GetTableMapByBeaconID(id string) (*TableMap, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	tableMap, contains := db.TableMaps[id]
	if !contains {
		return nil, nil
	}

	return &tableMap, nil
}

func (db *MemoryDB) UpdateCayanKey(token string) error {
	if db.Error != nil {
		return db.Error
	}

	db.CayanKeyVersion++
	db.CayanKey = &Key{
		Name:    "cayan",
		Version: db.CayanKeyVersion,
		Value:   token,
	}
	return nil
}

func (db *MemoryDB) GetCayanKey() (*Key, error) {
	return db.CayanKey, db.Error
}

func (db *MemoryDB) InsertToken(token *Token) error {
	if db.Error != nil {
		return db.Error
	}

	token.ID = DatabaseID(len(db.Tokens) + 1)
	db.Tokens[token.ID] = *token
	return nil
}

func (db *MemoryDB) GetToken(tokenString string) (*Token, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	for _, token := range db.Tokens {
		if token.Token == tokenString {
			return &token, nil
		}
	}

	return nil, nil
}

func (db *MemoryDB) DeleteToken(id DatabaseID) error {
	if db.Error != nil {
		return db.Error
	}

	delete(db.Tokens, id)
	return nil
}

func (db *MemoryDB) DeleteTokens(customerID DatabaseID) error {
	if db.Error != nil {
		return db.Error
	}

	for id, token := range db.Tokens {
		if token.CustomerID == customerID {
			delete(db.Tokens, id)
		}
	}
	return nil
}

func (db *MemoryDB) InsertCustomer(customer *Customer) error {
	if db.Error != nil {
		return db.Error
	}

	customer.ID = DatabaseID(len(db.Customers) + 1)
	db.Customers[customer.ID] = *customer
	return nil
}

func (db *MemoryDB) UpdateCustomerPassword(id DatabaseID, passwordHash string) (*Customer, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	customer, contains := db.Customers[id]
	if !contains {
		return nil, nil
	}

	customer.Password = passwordHash
	db.Customers[id] = customer
	return &customer, nil
}

func (db *MemoryDB) GetCustomer(id DatabaseID) (*Customer, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	customer, contains := db.Customers[id]
	if !contains {
		return nil, nil
	}
	return &customer, nil
}

func (db *MemoryDB) GetCustomerByExternalID(id string) (*Customer, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	for _, customer := range db.Customers {
		if customer.ExternalID == id {
			return &customer, nil
		}
	}
	return nil, nil
}

func (db *MemoryDB) GetCustomerByEmail(email string) (*Customer, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	for _, customer := range db.Customers {
		if customer.Email == email {
			return &customer, nil
		}
	}
	return nil, nil
}

func (db *MemoryDB) InsertOrder(order *Order) error {
	if db.Error != nil {
		return db.Error
	}

	order.ID = DatabaseID(len(db.Orders) + 1)
	for i := range order.Lines {
		line := &order.Lines[i]
		line.OrderID = order.ID
		db.LineCount++
		line.ID = DatabaseID(db.LineCount)

		for _, modifierID := range line.ModifierIDs {
			absoluteValueModifierID := KountaID(math.Abs(float64(modifierID)))
			modifier, err := db.GetMenuModifierByKountaID(order.SiteID, absoluteValueModifierID)
			if err != nil {
				return err
			}
			if modifier == nil {
				continue
			}

			modifier.SiteID = 0
			modifier.LineID = line.ID
			modifier.OrderID = line.OrderID
			modifier.Added = modifierID >= 0

			if modifier.Added {
				line.AddedModifiers = append(line.AddedModifiers, *modifier)
			} else {
				line.RemovedModifiers = append(line.RemovedModifiers, *modifier)
			}
		}
	}
	db.Orders[order.ID] = *order
	return nil
}

func (db *MemoryDB) UpdateOrder(order *Order) error {
	if db.Error != nil {
		return db.Error
	}

	existingOrder, err := db.GetOrder(order.PosID)
	if err != nil {
		return err
	}
	order.ID = existingOrder.ID

	for i := range order.Lines {
		line := &order.Lines[i]
		line.OrderID = order.ID
		db.LineCount++
		line.ID = DatabaseID(db.LineCount)

		for _, modifierID := range line.ModifierIDs {
			absoluteValueModifierID := KountaID(math.Abs(float64(modifierID)))
			modifier, err := db.GetMenuModifierByKountaID(order.SiteID, absoluteValueModifierID)
			if err != nil {
				return err
			}
			if modifier == nil {
				continue
			}

			modifier.SiteID = 0
			modifier.LineID = line.ID
			modifier.OrderID = line.OrderID
			modifier.Added = modifierID >= 0

			if modifier.Added {
				line.AddedModifiers = append(line.AddedModifiers, *modifier)
			} else {
				line.RemovedModifiers = append(line.RemovedModifiers, *modifier)
			}
		}
	}

	db.Orders[order.ID] = *order
	return nil
}

func (db *MemoryDB) UpdateOrderTableName(order *Order, tableName string) error {
	if db.Error != nil {
		return db.Error
	}

	existingOrder := db.Orders[order.ID]
	existingOrder.TableName = tableName
	db.Orders[order.ID] = existingOrder
	return nil
}

func (db *MemoryDB) UpdateOrderCustomerID(order *Order, customerID DatabaseID) error {
	if db.Error != nil {
		return db.Error
	}

	existingOrder := db.Orders[order.ID]
	existingOrder.CustomerID.Int64 = int64(customerID)
	existingOrder.CustomerID.Valid = customerID > 0
	db.Orders[order.ID] = existingOrder
	return nil
}

func (db *MemoryDB) UpdateOrderPickupTime(order *Order, pickupTime time.Time) error {
	if db.Error != nil {
		return db.Error
	}

	existingOrder := db.Orders[order.ID]
	existingOrder.PickupTime = &pickupTime
	db.Orders[order.ID] = existingOrder
	return nil
}

func (db *MemoryDB) GetOrder(orderID KountaID) (*Order, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	for _, order := range db.Orders {
		if order.PosID == orderID {
			return &order, nil
		}
	}

	return nil, nil
}

func (db *MemoryDB) GetOrderByDatabaseID(orderID DatabaseID) (*Order, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	order, contains := db.Orders[orderID]
	if !contains {
		return nil, nil
	}

	return &order, nil
}

func (db *MemoryDB) GetOrderByPagerID(siteID KountaID, pagerID int64) (*Order, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	pagerString := strconv.FormatInt(pagerID, 10)
	for _, order := range db.Orders {
		if order.SiteID == siteID && order.PagerNumber == pagerString {
			return &order, nil
		}
	}

	return nil, nil
}

func (db *MemoryDB) SelectOrdersByCustomerID(customerID DatabaseID) (*[]Order, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	orders := []Order{}
	for _, order := range db.Orders {
		if order.CustomerID.Int64 == int64(customerID) {
			orders = append(orders, order)
		}
	}

	return &orders, nil
}

func (db *MemoryDB) SelectOnHoldAndPendingOrdersByTable(siteID KountaID, tableName string) (*[]Order, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	orders := []Order{}
	for _, order := range db.Orders {
		if order.SiteID == siteID && order.TableName == tableName &&
			(order.Status == OrderStatusOnHold || order.Status == OrderStatusPending) {
			orders = append(orders, order)
		}
	}

	return &orders, nil
}

func (db *MemoryDB) SelectOnHoldAndPendingOrdersByPagerID(siteID KountaID, pagerID int64) (*[]Order, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	pagerString := strconv.FormatInt(pagerID, 10)
	orders := []Order{}
	for _, order := range db.Orders {
		if order.SiteID == siteID && order.PagerNumber == pagerString &&
			(order.Status == OrderStatusOnHold || order.Status == OrderStatusPending) {
			orders = append(orders, order)
		}
	}

	return &orders, nil
}

func (db *MemoryDB) InsertOrderUpdate(orderUpdate KountaOrderUpdate) error {
	return db.Error
}

func (db *MemoryDB) GetLine(lineID DatabaseID) (*Line, error) {
	for _, order := range db.Orders {
		for _, line := range order.Lines {
			if line.ID == lineID {
				return &line, nil
			}
		}
	}
	return nil, nil
}

func (db *MemoryDB) SelectLines(orderID DatabaseID) (*[]Line, error) {
	lines := db.Orders[orderID].Lines
	return &lines, db.Error
}

func (db *MemoryDB) SelectAddedModifiers(lineID DatabaseID) (*[]Modifier, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	for _, order := range db.Orders {
		for _, line := range order.Lines {
			if line.ID == lineID {
				if line.AddedModifiers == nil {
					return &[]Modifier{}, nil
				}
				return &line.AddedModifiers, nil
			}
		}
	}

	return &[]Modifier{}, nil
}

func (db *MemoryDB) SelectRemovedModifiers(lineID DatabaseID) (*[]Modifier, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	for _, order := range db.Orders {
		for _, line := range order.Lines {
			if line.ID == lineID {
				if line.RemovedModifiers == nil {
					return &[]Modifier{}, nil
				}
				return &line.RemovedModifiers, nil
			}
		}
	}

	return &[]Modifier{}, nil
}

func (db *MemoryDB) InsertPayment(payment *Payment, order *Order) error {
	if db.Error != nil {
		return db.Error
	}

	payment.TransactionID = strconv.FormatInt(rand.Int63(), 10)
	db.Payments[payment.TransactionID] = *payment

	existingOrder := db.Orders[order.ID]
	existingOrder.Status = order.Status
	db.Orders[order.ID] = existingOrder
	return nil
}

func (db *MemoryDB) GetPaymentByOrderID(id DatabaseID) (*Payment, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	for _, payment := range db.Payments {
		if payment.OrderID == id {
			return &payment, nil
		}
	}

	return nil, nil
}

func (db *MemoryDB) InsertSite(site *Site) error {
	if db.Error != nil {
		return db.Error
	}

	site.ID = DatabaseID(len(db.Sites) + 1)
	db.Sites[site.ID] = *site
	return nil
}

func (db *MemoryDB) UpdateSiteMenuHash(site *Site, menuHash string) error {
	if db.Error != nil {
		return db.Error
	}

	existingSite := db.Sites[site.ID]
	existingSite.MenuHash = menuHash
	db.Sites[site.ID] = existingSite
	return nil
}

func (db *MemoryDB) SelectSites() (*[]Site, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	sites := make([]Site, 0, len(db.Sites))
	for _, site := range db.Sites {
		sites = append(sites, site)
	}
	return &sites, nil
}

func (db *MemoryDB) GetSite(id KountaID) (*Site, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	for _, site := range db.Sites {
		if site.PosID == id {
			return &site, nil
		}
	}
	return nil, nil
}

func (db *MemoryDB) UpsertCategory(category *Category) error {
	if db.Error != nil {
		return db.Error
	}

	var existingID DatabaseID
	for _, v := range db.Categories {
		if v.PosID == category.PosID {
			existingID = v.ID
			break
		}
	}

	if existingID == 0 {
		category.ID = DatabaseID(len(db.Categories) + 1)
	} else {
		category.ID = existingID
	}

	db.Categories[category.ID] = *category
	return nil
}

func (db *MemoryDB) SelectCategories() (*[]Category, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	categoryIDs := make(databaseIDSlice, 0, len(db.Categories))
	for _, category := range db.Categories {
		categoryIDs = append(categoryIDs, category.ID)
	}

	sort.Sort(categoryIDs)
	categories := make([]Category, 0, len(categoryIDs))
	for _, categoryID := range categoryIDs {
		categories = append(categories, db.Categories[categoryID])
	}
	return &categories, nil
}

func (db *MemoryDB) SelectCategoriesBySiteID(siteID KountaID) (*[]Category, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	categoryIDs := make(databaseIDSlice, 0, len(db.Categories))
	for _, category := range db.Categories {
		if db.Sites[category.SiteID].PosID == siteID {
			categoryIDs = append(categoryIDs, category.ID)
		}
	}

	sort.Sort(categoryIDs)
	categories := make([]Category, 0, len(categoryIDs))
	for _, categoryID := range categoryIDs {
		categories = append(categories, db.Categories[categoryID])
	}
	return &categories, nil
}

func (db *MemoryDB) DeleteCategory(categoryID KountaID) error {
	if db.Error != nil {
		return db.Error
	}

	for id, category := range db.Categories {
		if category.PosID == categoryID {
			delete(db.Categories, id)

			for mid, item := range db.MenuItems {
				if item.CategoryID == id {
					delete(db.MenuItems, mid)
				}
			}
		}
	}
	return nil
}

func (db *MemoryDB) UpsertMenuItem(item *MenuItem) error {
	if db.Error != nil {
		return db.Error
	}

	var existingID DatabaseID
	for _, mi := range db.MenuItems {
		if mi.PosID == item.PosID {
			existingID = mi.ID

			break
		}
	}

	if existingID == 0 {
		item.ID = DatabaseID(len(db.MenuItems) + 1)
	} else {
		item.ID = existingID
	}

	db.MenuItems[item.ID] = *item
	return nil
}

func (db *MemoryDB) SelectMenuItems() (*[]MenuItem, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	itemIDs := make(databaseIDSlice, 0, len(db.MenuItems))
	for _, item := range db.MenuItems {
		itemIDs = append(itemIDs, item.ID)
	}

	sort.Sort(itemIDs)
	items := make([]MenuItem, 0, len(itemIDs))
	for _, itemID := range itemIDs {
		items = append(items, db.MenuItems[itemID])
	}
	return &items, nil
}

func (db *MemoryDB) SelectMenuItemsByCategoryID(siteID KountaID, categoryID DatabaseID) (*[]MenuItem, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	itemIDs := make(databaseIDSlice, 0, len(db.MenuItems))
	for _, item := range db.MenuItems {
		if db.Sites[item.SiteID].PosID == siteID && db.Categories[item.CategoryID].ID == categoryID {
			itemIDs = append(itemIDs, item.ID)
		}
	}

	sort.Sort(itemIDs)
	items := make([]MenuItem, 0, len(itemIDs))
	for _, itemID := range itemIDs {
		items = append(items, db.MenuItems[itemID])
	}
	return &items, nil
}

func (db *MemoryDB) GetMenuItem(siteID KountaID, menuItemID DatabaseID) (*MenuItem, error) {
	if db.Error != nil {
		return nil, db.Error
	}
	menuItem, contains := db.MenuItems[menuItemID]
	if !contains {
		return nil, nil
	}

	return &menuItem, nil
}

func (db *MemoryDB) DeleteMenuItem(menuItemID KountaID) error {
	if db.Error != nil {
		return db.Error
	}

	for id, menuItem := range db.MenuItems {
		if menuItem.PosID == menuItemID {
			delete(db.MenuItems, id)
		}
	}
	return nil
}

func (db *MemoryDB) UpsertMenuItemModifier(item *MenuItem, modifier *Modifier) error {
	if db.Error != nil {
		return db.Error
	}

	existingItem, contains := db.MenuItems[item.ID]
	if !contains {
		return errors.New("menu item not in database")
	}

	var existingModifierID DatabaseID
	for _, m := range db.Modifiers {
		if m.PosID == modifier.PosID {
			existingModifierID = m.ID
			break
		}
	}

	if existingModifierID == 0 {
		modifier.ID = DatabaseID(len(db.Modifiers) + 1)
	} else {
		modifier.ID = existingModifierID
	}

	db.Modifiers[modifier.ID] = *modifier

	for _, v := range existingItem.Modifiers {
		if v.ID == modifier.ID {
			// Already included, return early
			return nil
		}
	}

	existingItem.Modifiers = append(existingItem.Modifiers, *modifier)
	db.MenuItems[existingItem.ID] = existingItem

	return nil
}

func (db *MemoryDB) UpsertOptionSetModifier(optionSet *OptionSet, modifier *Modifier) error {
	if db.Error != nil {
		return db.Error
	}

	existingOptionSet, contains := db.OptionSets[optionSet.ID]
	if !contains {
		return errors.New("option set not in database")
	}

	var existingModifierID DatabaseID
	for _, v := range db.Modifiers {
		if v.PosID == modifier.PosID {
			existingModifierID = v.ID
			break
		}
	}

	if existingModifierID == 0 {
		modifier.ID = DatabaseID(len(db.Modifiers) + 1)
	} else {
		modifier.ID = existingModifierID
	}

	db.Modifiers[modifier.ID] = *modifier

	for _, o := range existingOptionSet.Options {
		if o.ID == modifier.ID {
			// Already included, return early
			return nil
		}
	}

	existingOptionSet.Options = append(existingOptionSet.Options, *modifier)
	db.OptionSets[existingOptionSet.ID] = existingOptionSet

	return nil
}

func (db *MemoryDB) GetMenuModifier(siteID KountaID, modifierID DatabaseID) (*Modifier, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	for _, modifier := range db.Modifiers {
		if modifier.ID == modifierID && db.Sites[modifier.SiteID].PosID == siteID {
			return &modifier, nil
		}
	}
	return nil, nil
}

func (db *MemoryDB) GetMenuModifierByKountaID(siteID, modifierID KountaID) (*Modifier, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	for _, modifier := range db.Modifiers {
		if modifier.PosID == modifierID && db.Sites[modifier.SiteID].PosID == siteID {
			return &modifier, nil
		}
	}
	return nil, nil
}

func (db *MemoryDB) SelectMenuModifiers() (*[]Modifier, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	modifierIDs := make(databaseIDSlice, 0, len(db.Modifiers))
	for _, modifier := range db.Modifiers {
		modifierIDs = append(modifierIDs, modifier.ID)
	}

	sort.Sort(modifierIDs)
	modifiers := make([]Modifier, 0, len(modifierIDs))
	for _, modifierID := range modifierIDs {
		modifiers = append(modifiers, db.Modifiers[modifierID])
	}
	return &modifiers, nil
}

func (db *MemoryDB) SelectMenuItemModifiers(siteID KountaID, menuItemID DatabaseID) (*[]Modifier, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	var itemModifiers []Modifier
	for _, menuItem := range db.MenuItems {
		if menuItem.ID == menuItemID && db.Sites[menuItem.SiteID].PosID == siteID {
			itemModifiers = menuItem.Modifiers
			break
		}
	}

	if itemModifiers == nil {
		return &[]Modifier{}, nil
	}

	modifiers := make([]Modifier, 0, len(itemModifiers))
	for _, v := range itemModifiers {
		modifier, contains := db.Modifiers[v.ID]
		if contains {
			modifiers = append(modifiers, modifier)
		}
	}

	return &modifiers, nil
}

func (db *MemoryDB) DeleteMenuModifier(modifierID KountaID) error {
	if db.Error != nil {
		return db.Error
	}

	for id, modifier := range db.Modifiers {
		if modifier.PosID == modifierID {
			delete(db.Modifiers, id)
		}
	}

	return nil
}

func (db *MemoryDB) UpsertOptionSet(item *MenuItem, optionSet *OptionSet) error {
	if db.Error != nil {
		return db.Error
	}

	existingItem, contains := db.MenuItems[item.ID]
	if !contains {
		return errors.New("menu item not in database")
	}

	var existingOptionSetID DatabaseID
	for _, v := range db.OptionSets {
		if v.PosID == optionSet.PosID {
			existingOptionSetID = v.ID
			break
		}
	}

	if existingOptionSetID == 0 {
		optionSet.ID = DatabaseID(len(db.OptionSets) + 1)
	} else {
		optionSet.ID = existingOptionSetID
	}

	db.OptionSets[optionSet.ID] = *optionSet

	for _, v := range existingItem.OptionSets {
		if v.ID == optionSet.ID {
			// Already included, return early
			return nil
		}
	}

	existingItem.OptionSets = append(existingItem.OptionSets, *optionSet)
	db.MenuItems[existingItem.ID] = existingItem

	return nil
}

func (db *MemoryDB) GetOptionSet(optionSetID DatabaseID) (*OptionSet, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	for _, optionSet := range db.OptionSets {
		if optionSet.ID == optionSetID {
			return &optionSet, nil
		}
	}
	return nil, nil
}

func (db *MemoryDB) SelectOptionSets() (*[]OptionSet, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	optionSetIDs := make(databaseIDSlice, 0, len(db.OptionSets))
	for _, optionSet := range db.OptionSets {
		optionSetIDs = append(optionSetIDs, optionSet.ID)
	}

	sort.Sort(optionSetIDs)
	optionSets := make([]OptionSet, 0, len(optionSetIDs))
	for _, optionSetID := range optionSetIDs {
		optionSets = append(optionSets, db.OptionSets[optionSetID])
	}
	return &optionSets, nil
}

func (db *MemoryDB) SelectOptionSetsByItemID(menuItemID DatabaseID) (*[]OptionSet, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	var itemOptionSets []OptionSet
	for _, menuItem := range db.MenuItems {
		if menuItem.ID == menuItemID {
			itemOptionSets = menuItem.OptionSets
			break
		}
	}

	if itemOptionSets == nil {
		return &[]OptionSet{}, nil
	}

	optionSets := make([]OptionSet, 0, len(itemOptionSets))
	for _, v := range itemOptionSets {
		optionSet, contains := db.OptionSets[v.ID]
		if contains {
			optionSets = append(optionSets, optionSet)
		}
	}

	return &optionSets, nil
}
func (db *MemoryDB) DeleteOptionSet(optionSetID KountaID) error {
	if db.Error != nil {
		return db.Error
	}

	for id, optionSet := range db.OptionSets {
		if optionSet.PosID == optionSetID {
			delete(db.OptionSets, id)
		}
	}

	return nil
}
