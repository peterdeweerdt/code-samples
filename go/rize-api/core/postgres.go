package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Used to register the Postgres driver with database/sql
	"github.com/pkg/errors"
	"pjd"
)

type Postgres struct {
	*sqlx.DB
}

// InitPG connects to a Postgres database
func InitPG(url string) (Postgres, error) {
	return initPG(url, false)
}

// CleanPG connects to a Postgres database and clears all tables
func CleanPG(url string) (Postgres, error) {
	return initPG(url, true)
}

func initPG(url string, clean bool) (pg Postgres, err error) {
	if url == "" {
		log.Fatalln("DB URL was empty")
	}

	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return pg, err
	}
	db.MapperFunc(pjd.ToSnakeCase)

	pg = Postgres{db}
	if err != nil {
		return pg, err
	}

	if clean {
		err = pg.dropTables()
		if err != nil {
			return pg, err
		}
	}

	err = pg.runMigrations()
	if err != nil {
		return pg, err
	}

	return pg, err
}

func (pg Postgres) dropTables() error {
	tableNames := []string{
		"customers",
		"keys",
		"kounta_log",
		"lines",
		"menu_categories",
		"menu_item_modifiers_mapping",
		"menu_item_option_sets_mapping",
		"menu_items",
		"menu_modifiers",
		"menu_option_set_modifiers_mapping",
		"menu_option_sets",
		"migrations",
		"modifiers",
		"orders",
		"payments",
		"site_menu_categories_mapping",
		"site_menu_items_pricing",
		"site_menu_modifiers_pricing",
		"sites",
		"table_mapping",
		"tokens",
	}

	var err error
	for _, table := range tableNames {
		_, err = pg.Exec("DROP TABLE IF EXISTS " + string(table) + " CASCADE")
		if err != nil {
			log.Println("Error cleaning:", err)
		}
	}
	return err
}

func (pg Postgres) runMigrations() error {
	log.Println("Running migrations")

	migrationsDir := "migrations"
	allMigrations, err := ioutil.ReadDir(migrationsDir) //when running from cmd/rize/
	if err != nil {
		migrationsDir = "../migrations"
		m, err := ioutil.ReadDir(migrationsDir) //when running from rize/ (tests)
		if err != nil {
			return err
		}
		allMigrations = m
	}
	if len(allMigrations) == 0 {
		return errors.Errorf("Could not find any migrations")
	}

	var current int64
	err = pg.Get(&current, "SELECT max(version) FROM migrations")
	if err != nil {
		current = 0
	}

	log.Println("Current DB version:", current)

	err = pg.transact(func(tx *sqlx.Tx) error {
		for _, file := range allMigrations {
			parts := strings.Split(file.Name(), "_")
			version, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return err
			}

			if version > current {
				log.Println("Running migration:", file.Name())

				bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", migrationsDir, file.Name()))
				if err != nil {
					return err
				}

				_, err = tx.Exec(string(bytes))
				if err != nil {
					return err
				}

				_, err = tx.Exec("INSERT INTO migrations (version, file_name, ran) VALUES ($1, $2, current_timestamp)", version, file.Name())
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	return err
}

func (pg Postgres) transact(exec func(tx *sqlx.Tx) error) (err error) {
	tx, err := pg.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = exec(tx)

	return err
}

/// DB Interface

// Table Mapping

func (pg Postgres) InsertTableMap(tableMap *TableMap) error {
	_, err := pg.Exec(`INSERT INTO table_mapping (beacon_id, site_id, table_name)
		VALUES($1, $2, $3)`,
		strings.ToLower(tableMap.BeaconID), tableMap.SiteID, tableMap.TableName)
	return err
}

func (pg Postgres) GetTableMapByBeaconID(id string) (*TableMap, error) {
	tableMap := TableMap{}
	err := pg.Get(&tableMap, "SELECT * FROM table_mapping WHERE beacon_id = $1 LIMIT 1", strings.ToLower(id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &tableMap, err
}

// Cayan

func (pg Postgres) UpdateCayanKey(token string) error {
	_, err := pg.Exec(`INSERT INTO keys (name, version, value)
			   VALUES ('cayan', 1, $1)
			   ON CONFLICT (name)
			   DO UPDATE SET version = keys.version + 1, value = EXCLUDED.value;`, token)
	return err
}

func (pg Postgres) GetCayanKey() (*Key, error) {
	key := Key{}
	err := pg.Get(&key, `SELECT * FROM keys WHERE name = 'cayan' ORDER BY version DESC LIMIT 1`)
	return &key, err
}

// Token

func (pg Postgres) InsertToken(token *Token) error {
	return pg.QueryRow(`INSERT INTO tokens (service, name, token, customer_id, expiry)
			    VALUES($1, $2, $3, $4, $5)
			    RETURNING id`,
		token.Service, token.Name, token.Token, token.CustomerID, token.Expiry).Scan(&token.ID)
}

func (pg Postgres) GetToken(tokenString string) (*Token, error) {
	token := Token{}
	err := pg.Get(&token, `SELECT id, service, name, token, customer_id, expiry FROM tokens WHERE token = $1`, tokenString)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &token, err
}

func (pg Postgres) DeleteToken(id DatabaseID) error {
	_, err := pg.Exec(`DELETE FROM tokens WHERE id = $1`, id)
	return err
}

func (pg Postgres) DeleteTokens(customerID DatabaseID) error {
	_, err := pg.Exec(`DELETE FROM tokens WHERE customer_id = $1`, customerID)
	return err
}

// Customer

func (pg Postgres) InsertCustomer(c *Customer) error {
	return pg.Get(c,
		`INSERT INTO customers(first_name, last_name, email, password, image_url, external_id, phone, pos_id)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`,
		c.FirstName, c.LastName, c.Email, c.Password, c.ImageURL, c.ExternalID, c.Phone, c.PosID)
}

func (pg Postgres) UpdateCustomerPassword(id DatabaseID, passwordHash string) (*Customer, error) {
	customer := Customer{}
	err := pg.Get(&customer, `UPDATE customers SET password = $1 WHERE id = $2 RETURNING *`, passwordHash, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &customer, err
}

func (pg Postgres) GetCustomer(id DatabaseID) (*Customer, error) {
	customer := Customer{}
	err := pg.Get(&customer, `SELECT * FROM customers WHERE id = $1 LIMIT 1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &customer, err
}

func (pg Postgres) GetCustomerByExternalID(id string) (*Customer, error) {
	customer := Customer{}
	err := pg.Get(&customer, `SELECT * FROM customers WHERE external_id = $1 LIMIT 1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &customer, err
}

func (pg Postgres) GetCustomerByEmail(email string) (*Customer, error) {
	customer := Customer{}
	err := pg.Get(&customer, `SELECT * FROM customers WHERE email = $1 LIMIT 1`, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &customer, err
}

// Order

func (pg Postgres) InsertOrder(order *Order) error {
	return pg.transact(func(tx *sqlx.Tx) error {
		err := tx.QueryRow(
			`INSERT INTO orders (pos_id, status, table_name, total, total_tax, pager_number, site_id, customer_id, created_at)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id`,
			order.PosID,
			order.Status,
			order.TableName,
			order.Total,
			order.TotalTax,
			order.PagerNumber,
			order.SiteID,
			order.CustomerID,
			time.Now()).
			Scan(&order.ID)
		if err != nil {
			return err
		}

		if err = pg.insertLines(tx, order); err != nil {
			return err
		}

		return nil
	})
}

func (pg Postgres) UpdateOrder(order *Order) error {
	return pg.transact(func(tx *sqlx.Tx) error {
		err := tx.QueryRow(
			`UPDATE orders
			SET status = $1, customer_id = $2, total = $3, total_tax = $4, pager_number = $5
			WHERE pos_id = $6
			RETURNING id`,
			order.Status, order.CustomerID, order.Total, order.TotalTax, order.PagerNumber, order.PosID).
			Scan(&order.ID)
		if err != nil {
			return errors.Wrap(err, "update order")
		}

		// this was split out from the UPDATE above to prevent Kounta from clearing out the table value from LRS
		if order.TableName != "" {
			_, err = tx.Exec(`UPDATE orders SET table_name = $1 WHERE pos_id = $2`, order.TableName, order.PosID)
			if err != nil {
				return errors.Wrap(err, "update order")
			}
		}

		if err = pg.insertLines(tx, order); err != nil {
			return errors.Wrap(err, "update order")
		}

		return nil
	})
}

func (pg Postgres) insertLines(tx *sqlx.Tx, o *Order) error {
	// we have to remove all existing lines because there's not way to uniquely identify them in Kounta to update them

	if _, err := tx.Exec("DELETE FROM modifiers WHERE order_id = $1", o.ID); err != nil {
		return errors.Wrapf(err, "insertLines: error removing modifiers for order '%d'", o.ID)
	}

	if _, err := tx.Exec("DELETE FROM lines WHERE order_id = $1", o.ID); err != nil {
		return errors.Wrapf(err, "insertLines: error deleting lines for order %d", o.ID)
	}

	lineInsert, err := tx.Preparex(`INSERT INTO lines (price, product_name, notes, quantity, total, total_tax, order_id, pos_id)
									VALUES($1, $2, $3, $4, $5, $6, $7, $8)
									RETURNING id`)
	if err != nil {
		return errors.Wrap(err, "insertLines: error preparing lines insert")
	}
	defer lineInsert.Close()

	modInsert, err := tx.Preparex(`INSERT INTO modifiers (name, added, price, price_ex_tax, pos_id, line_id, order_id)
								   VALUES($1, $2, $3, $4, $5, $6, $7)`)
	if err != nil {
		return errors.Wrapf(err, "insertLines: preparing to insert modifiers for lines on order %d", o.ID)
	}
	defer modInsert.Close()

	for i := range o.Lines {
		line := &o.Lines[i] // update the original instead of a copy

		line.OrderID = o.ID

		err = lineInsert.Get(&line.ID, line.Price, line.ProductName, line.Notes, line.Quantity, line.Total, line.TotalTax, line.OrderID, line.PosID)
		if err != nil {
			return errors.Wrapf(err, "error adding line '%s' to order", line.ProductName)
		}

		for _, modifierID := range line.ModifierIDs {
			// when we get a line from Kounta, its modifiers may be positive or negative (depending on added or removed)
			// however when we get modifiers for a product, they are always positive
			// so to compare them we have to take absolute value of the modifier that we are looking for
			absoluteValueModifierID := KountaID(math.Abs(float64(modifierID)))
			modifier, err := pg.GetMenuModifierByKountaID(o.SiteID, absoluteValueModifierID)
			if err != nil {
				return errors.Wrapf(err, "Could not get modifier %d for line %d", modifier, line.ID)
			}
			if modifier == nil {
				//this must be an Option Set, not a Modifier...ignoring
				continue
			}

			modifier.SiteID = 0
			modifier.LineID = line.ID
			modifier.OrderID = line.OrderID
			modifier.Added = modifierID >= 0

			_, err = modInsert.Exec(modifier.Name, modifier.Added, modifier.PriceWithTax, modifier.Price, modifier.PosID, modifier.LineID, modifier.OrderID)
			if err != nil {
				return errors.Wrapf(err, "error adding modifier '%s' to line %d", modifier.Name, line.ID)
			}

			if modifier.Added {
				line.AddedModifiers = append(line.AddedModifiers, *modifier)
			} else {
				line.RemovedModifiers = append(line.RemovedModifiers, *modifier)
			}
		}
	}

	return nil
}

func (pg Postgres) UpdateOrderCustomerID(order *Order, customerID DatabaseID) error {
	_, err := pg.Exec(`UPDATE orders SET customer_id = $1 WHERE id = $2`, customerID, order.ID)
	return err
}

func (pg Postgres) UpdateOrderPickupTime(order *Order, pickupTime time.Time) error {
	_, err := pg.Exec(
		`UPDATE orders
		SET pickup_time = $1
		WHERE id = $2`,
		pickupTime.Format(time.RFC3339), order.ID)
	return err
}

func (pg Postgres) UpdateOrderTableName(order *Order, tableName string) error {
	_, err := pg.Exec(`UPDATE orders
			      SET table_name = $1
			      WHERE site_id = $2 AND pager_number = $3`, tableName, order.SiteID, order.PagerNumber)
	return err
}

func (pg Postgres) GetOrder(orderID KountaID) (*Order, error) {
	order := Order{}
	err := pg.Get(&order, `SELECT * FROM orders WHERE pos_id = $1`, orderID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &order, err
}

func (pg Postgres) GetOrderByDatabaseID(orderID DatabaseID) (*Order, error) {
	order := Order{}
	err := pg.Get(&order, `SELECT * FROM orders WHERE id = $1`, orderID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &order, err
}

func (pg Postgres) GetOrderByPagerID(siteID KountaID, pagerID int64) (*Order, error) {
	pagerString := strconv.FormatInt(pagerID, 10)
	order := Order{}
	err := pg.Get(&order, `SELECT * FROM orders WHERE site_id = $1 AND pager_number = $2`, siteID, pagerString)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &order, err
}

func (pg Postgres) SelectOrdersByCustomerID(customerID DatabaseID) (*[]Order, error) {
	orders := []Order{}
	err := pg.Select(&orders, `SELECT * FROM orders WHERE customer_id = $1`, customerID)
	return &orders, err
}

func (pg Postgres) SelectOnHoldAndPendingOrdersByTable(siteID KountaID, tableName string) (*[]Order, error) {
	orders := []Order{}
	err := pg.Select(&orders, `
		SELECT * FROM orders
		WHERE table_name = $1
		AND site_id = $2
		AND (status = $3 OR status = $4)`,
		tableName, siteID, OrderStatusOnHold, OrderStatusPending)
	return &orders, err
}

func (pg Postgres) SelectOnHoldAndPendingOrdersByPagerID(siteID KountaID, pagerID int64) (*[]Order, error) {
	pagerString := strconv.FormatInt(pagerID, 10)
	orders := []Order{}
	err := pg.Select(&orders, `
		SELECT * FROM orders
		WHERE site_id = $1
		AND pager_number = $2
		AND (status = $3 OR status = $4)`,
		siteID, pagerString, OrderStatusOnHold, OrderStatusPending)
	return &orders, err
}

func (pg Postgres) InsertOrderUpdate(orderUpdate KountaOrderUpdate) error {
	lines, err := json.Marshal(orderUpdate.GetLines())
	if err != nil {
		return errors.Wrap(err, "LogOrderUpdate")
	}

	payments, err := json.Marshal(orderUpdate.GetPayments())
	if err != nil {
		return errors.Wrap(err, "LogOrderUpdate")
	}

	lock, err := json.Marshal(orderUpdate.GetLock())
	if err != nil {
		return errors.Wrap(err, "LogOrderUpdate")
	}

	_, err = pg.DB.Exec(`INSERT INTO kounta_log (order_id, sale_number, created_at, updated_at, deleted, status, notes, total, paid, tips, register_id, site_id, lines, price_variation, payments, lock, staff_member_id)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`,
		orderUpdate.GetOrderID(),
		orderUpdate.GetSaleNumber(),
		orderUpdate.GetCreatedAt(),
		orderUpdate.GetUpdatedAt(),
		orderUpdate.GetDeleted(),
		orderUpdate.GetStatus(),
		orderUpdate.GetNotes(),
		orderUpdate.GetTotal(),
		orderUpdate.GetPaid(),
		orderUpdate.GetTips(),
		orderUpdate.GetRegisterID(),
		orderUpdate.GetSiteID(),
		lines,
		orderUpdate.GetPriceVariation(),
		payments,
		lock,
		orderUpdate.GetStaffMemberID())
	if err != nil {
		return errors.Wrap(err, "LogOrderUpdate")
	}

	return nil
}

func (pg Postgres) GetLine(lineID DatabaseID) (*Line, error) {
	line := Line{}
	err := pg.Get(&line, `SELECT * FROM lines WHERE id = $1`, lineID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &line, err
}

func (pg Postgres) SelectLines(orderID DatabaseID) (*[]Line, error) {
	lines := []Line{}
	err := pg.Select(&lines, `SELECT * FROM lines WHERE order_id = $1`, orderID)
	return &lines, err
}

func (pg Postgres) SelectAddedModifiers(lineID DatabaseID) (*[]Modifier, error) {
	modifiers := []Modifier{}
	err := pg.Select(&modifiers, `SELECT * FROM modifiers WHERE added = TRUE AND line_id = $1`, lineID)
	return &modifiers, err
}

func (pg Postgres) SelectRemovedModifiers(lineID DatabaseID) (*[]Modifier, error) {
	modifiers := []Modifier{}
	err := pg.Select(&modifiers, `SELECT * FROM modifiers WHERE added = FALSE AND line_id = $1`, lineID)
	return &modifiers, err
}

// Payment

func (pg Postgres) InsertPayment(payment *Payment, order *Order) error {
	return pg.transact(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(
			`INSERT INTO payments (amount, tip, transaction_id, date, customer_id, order_id)
			VALUES($1, $2, $3, $4, $5, $6)`,
			payment.Amount,
			payment.Tip,
			payment.TransactionID,
			payment.Date,
			payment.CustomerID,
			payment.OrderID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`UPDATE orders SET status = $1, pager_number = '' WHERE id = $2`, order.Status, order.ID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (pg Postgres) GetPaymentByOrderID(id DatabaseID) (*Payment, error) {
	payment := Payment{}
	err := pg.Get(&payment, "SELECT * FROM payments WHERE order_id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &payment, err
}

// Menu

func (pg Postgres) InsertSite(site *Site) error {
	return pg.QueryRow(
		`INSERT INTO sites(pos_id, menu_hash, name, updated_at)
			    VALUES($1, $2, $3, $4)
			    RETURNING id`,
		site.PosID,
		site.MenuHash,
		site.Name,
		time.Now()).Scan(&site.ID)
}

func (pg Postgres) UpdateSiteMenuHash(site *Site, menuHash string) error {
	_, err := pg.Exec(
		`UPDATE sites
		SET menu_hash = $1, updated_at = $2
		WHERE pos_id = $3`,
		menuHash,
		time.Now(),
		site.PosID)
	return err
}

func (pg Postgres) SelectSites() (*[]Site, error) {
	sites := []Site{}
	err := pg.Select(&sites, `SELECT * FROM sites`)
	return &sites, err
}

func (pg Postgres) GetSite(id KountaID) (*Site, error) {
	site := Site{}
	err := pg.Get(&site, `SELECT * FROM sites WHERE pos_id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &site, err
}

func (pg Postgres) UpsertCategory(category *Category) error {
	return pg.transact(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(
			`INSERT INTO menu_categories (name, pos_id)
			VALUES ($1, $2)
			ON CONFLICT (pos_id) DO UPDATE SET name = EXCLUDED.name`,
			category.Name, category.PosID)
		if err != nil {
			return err
		}

		if err := tx.Get(category, `SELECT * FROM menu_categories WHERE pos_id = $1`, category.PosID); err != nil {
			return err
		}

		_, err = tx.Exec(
			`INSERT INTO site_menu_categories_mapping (site_id, menu_category_id)
			VALUES ($1, $2)
			ON CONFLICT (site_id, menu_category_id) DO NOTHING`,
			category.SiteID, category.ID)

		return err
	})
}

func (pg Postgres) SelectCategories() (*[]Category, error) {
	categories := []Category{}
	if err := pg.Select(&categories, `SELECT * FROM menu_categories`); err != nil {
		return nil, errors.Wrap(err, "select categories")
	}
	return &categories, nil
}

func (pg Postgres) SelectCategoriesBySiteID(siteID KountaID) (*[]Category, error) {
	categories := []Category{}
	err := pg.Select(&categories,
		`SELECT category.*
		FROM menu_categories category
		JOIN site_menu_categories_mapping m ON category.id = m.menu_category_id
		JOIN sites site ON m.site_id = site.id
		WHERE site.pos_id = $1`,
		siteID)
	return &categories, err
}

func (pg Postgres) DeleteCategory(categoryID KountaID) error {
	_, err := pg.Exec(`DELETE FROM menu_categories WHERE pos_id = $1`, categoryID)
	return err
}

func (pg Postgres) UpsertMenuItem(item *MenuItem) error {
	return pg.transact(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(
			`INSERT INTO menu_items (pos_id, name, category_id, description)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (pos_id) DO UPDATE SET (name, category_id, description) = (EXCLUDED.name, EXCLUDED.category_id, EXCLUDED.description)`,
			item.PosID, item.Name, item.CategoryID, item.Description)
		if err != nil {
			return err
		}

		if err := tx.Get(item, `SELECT * FROM menu_items WHERE pos_id = $1`, item.PosID); err != nil {
			return err
		}

		_, err = tx.Exec(
			`INSERT INTO site_menu_items_pricing (site_id, menu_item_id, price)
			VALUES ($1, $2, $3)
			ON CONFLICT (site_id, menu_item_id) DO UPDATE SET price = EXCLUDED.price`,
			item.SiteID, item.ID, item.Price)

		return err
	})
}

func (pg Postgres) SelectMenuItems() (*[]MenuItem, error) {
	menuItems := []MenuItem{}
	err := pg.Select(&menuItems, `SELECT * FROM menu_items`)
	return &menuItems, err
}

// SelectMenuItems looks up menu items by Kounta site and category ids. These aren't stored with the menu items
// so a join with the categories and site tables are necessary to match by the Kounta ids
func (pg Postgres) SelectMenuItemsByCategoryID(siteID KountaID, categoryID DatabaseID) (*[]MenuItem, error) {
	menuItems := []MenuItem{}
	err := pg.Select(&menuItems,
		`SELECT m.*, p.price
		FROM menu_items m
		JOIN menu_categories c ON m.category_id = c.id
		JOIN site_menu_items_pricing p ON p.menu_item_id = m.id
		JOIN sites s ON p.site_id = s.id
		WHERE c.id = $1 AND s.pos_id = $2`,
		categoryID, siteID)
	return &menuItems, err
}

// GetMenuItem gets a single menu item by database ID
func (pg Postgres) GetMenuItem(siteID KountaID, menuItemID DatabaseID) (*MenuItem, error) {
	m := MenuItem{}
	err := pg.Get(&m,
		`SELECT m.*, p.price
		FROM menu_items m
		JOIN site_menu_items_pricing p ON m.id = p.menu_item_id
		JOIN sites s ON p.site_id = s.id
		WHERE m.id = $1 AND s.pos_id = $2`, menuItemID, siteID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &m, err
}

func (pg Postgres) DeleteMenuItem(menuItemID KountaID) error {
	_, err := pg.Exec(`DELETE FROM menu_items WHERE pos_id = $1`, menuItemID)
	return err
}

func (pg Postgres) UpsertMenuItemModifier(item *MenuItem, modifier *Modifier) error {
	return pg.transact(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(
			`INSERT INTO menu_modifiers (name, pos_id)
			VALUES ($1, $2)
			ON CONFLICT (pos_id) DO UPDATE SET name = EXCLUDED.name`,
			modifier.Name, modifier.PosID)
		if err != nil {
			return err
		}

		if err := tx.Get(modifier, `SELECT * FROM menu_modifiers WHERE pos_id = $1`, modifier.PosID); err != nil {
			return err
		}

		_, err = tx.Exec(
			`INSERT INTO menu_item_modifiers_mapping (modifier_id, menu_item_id)
			VALUES ($1, $2)
			ON CONFLICT (modifier_id, menu_item_id) DO NOTHING`,
			modifier.ID, item.ID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(
			`INSERT INTO site_menu_modifiers_pricing (site_id, menu_modifier_id, price)
			VALUES ($1, $2, $3)
			ON CONFLICT (site_id, menu_modifier_id) DO UPDATE SET price = EXCLUDED.price`,
			modifier.SiteID, modifier.ID, modifier.Price)

		return err
	})
}

func (pg Postgres) UpsertOptionSetModifier(optionSet *OptionSet, modifier *Modifier) error {
	return pg.transact(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(
			`INSERT INTO menu_modifiers (name, pos_id)
			VALUES ($1, $2)
			ON CONFLICT (pos_id) DO UPDATE SET name = EXCLUDED.name`,
			modifier.Name, modifier.PosID)
		if err != nil {
			return err
		}

		if err := tx.Get(modifier, `SELECT * FROM menu_modifiers WHERE pos_id = $1`, modifier.PosID); err != nil {
			return err
		}

		_, err = tx.Exec(
			`INSERT INTO menu_option_set_modifiers_mapping (option_set_id, modifier_id, price)
			VALUES ($1, $2, $3)
			ON CONFLICT (option_set_id, modifier_id) DO UPDATE SET price = EXCLUDED.price`,
			optionSet.ID, modifier.ID, modifier.Price)
		if err != nil {
			return err
		}

		_, err = tx.Exec(
			`INSERT INTO site_menu_modifiers_pricing (site_id, menu_modifier_id, price)
			VALUES ($1, $2, $3)
			ON CONFLICT (site_id, menu_modifier_id) DO UPDATE SET price = EXCLUDED.price`,
			modifier.SiteID, modifier.ID, modifier.Price)

		return err
	})
}

func (pg Postgres) GetMenuModifier(siteID KountaID, modifierID DatabaseID) (*Modifier, error) {
	m := Modifier{}
	err := pg.Get(&m,
		`SELECT m.*, p.price
		FROM menu_modifiers m
		JOIN site_menu_modifiers_pricing p ON m.id = p.menu_modifier_id
		JOIN sites s ON s.id = p.site_id
		WHERE m.id = $1 AND s.pos_id = $2`, modifierID, siteID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get modifier %d", modifierID))
	}

	return &m, nil
}

func (pg Postgres) GetMenuModifierByKountaID(siteID, modifierID KountaID) (*Modifier, error) {
	m := Modifier{}
	err := pg.Get(&m,
		`SELECT m.*, p.price
		FROM menu_modifiers m
		JOIN site_menu_modifiers_pricing p ON m.id = p.menu_modifier_id
		JOIN sites s ON s.id = p.site_id
		WHERE m.pos_id = $1 AND s.pos_id = $2`, modifierID, siteID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get modifier %d", modifierID))
	}

	return &m, nil
}

func (pg Postgres) SelectMenuModifiers() (*[]Modifier, error) {
	modifiers := []Modifier{}
	err := pg.Select(&modifiers, `SELECT * FROM menu_modifiers`)
	return &modifiers, err
}

func (pg Postgres) SelectMenuItemModifiers(siteID KountaID, menuItemID DatabaseID) (*[]Modifier, error) {
	modifiers := []Modifier{}
	err := pg.Select(&modifiers,
		`SELECT m.*
		FROM menu_items i
		JOIN menu_item_modifiers_mapping mim ON i.id = mim.menu_item_id
		JOIN menu_modifiers m ON m.id = mim.modifier_id
		JOIN site_menu_modifiers_pricing p ON m.id = p.menu_modifier_id
		JOIN sites s ON p.site_id = s.id
		WHERE s.pos_id = $1 AND i.id = $2`, siteID, menuItemID)
	return &modifiers, err
}

func (pg Postgres) DeleteMenuModifier(modifierID KountaID) error {
	_, err := pg.Exec(`DELETE FROM menu_modifiers WHERE pos_id = $1`, modifierID)
	return err
}

func (pg Postgres) UpsertOptionSet(item *MenuItem, optionSet *OptionSet) error {
	return pg.transact(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(
			`INSERT INTO menu_option_sets (name, pos_id, min_selection, max_selection)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (pos_id) DO UPDATE SET (name, min_selection, max_selection) = (EXCLUDED.name, EXCLUDED.min_selection, EXCLUDED.max_selection)`,
			optionSet.Name, optionSet.PosID, optionSet.MinSelection, optionSet.MaxSelection)
		if err != nil {
			return err
		}

		if err := tx.Get(optionSet, `SELECT * FROM menu_option_sets WHERE pos_id = $1`, optionSet.PosID); err != nil {
			return err
		}

		_, err = tx.Exec(
			`INSERT INTO menu_item_option_sets_mapping (menu_item_id, option_set_id)
			VALUES ($1, $2)
			ON CONFLICT (menu_item_id, option_set_id) DO NOTHING`,
			item.ID, optionSet.ID)
		if err != nil {
			return err
		}

		return err
	})
}

func (pg Postgres) GetOptionSet(optionSetID DatabaseID) (*OptionSet, error) {
	optionSet := OptionSet{}
	err := pg.Get(&optionSet, `SELECT * FROM menu_option_sets WHERE id = $1`, optionSetID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &optionSet, err
}

func (pg Postgres) SelectOptionSets() (*[]OptionSet, error) {
	optionSets := []OptionSet{}
	err := pg.Select(&optionSets, `SELECT * FROM menu_option_sets`)
	return &optionSets, err
}

func (pg Postgres) SelectOptionSetsByItemID(menuItemID DatabaseID) (*[]OptionSet, error) {
	optionSets := []OptionSet{}
	err := pg.Select(&optionSets,
		`SELECT o.*
		FROM menu_option_sets o
		JOIN menu_item_option_sets_mapping mm ON o.id = mm.option_set_id
		WHERE mm.menu_item_id = $1`, menuItemID)
	if err != nil {
		return nil, err
	}
	if len(optionSets) == 0 {
		return &optionSets, nil
	}

	for i, o := range optionSets {
		options := []Modifier{}
		err := pg.Select(&options,
			// Using price_ex_tax because old tables have price as price_ex_tax and reusing Modifier
			`SELECT m.*, mom.price AS price_ex_tax
			FROM menu_modifiers m
			JOIN menu_option_set_modifiers_mapping mom ON m.id = mom.modifier_id
			WHERE mom.option_set_id = $1`, o.ID)
		if err != nil {
			return nil, err
		}
		o.Options = options
		optionSets[i] = o
	}

	return &optionSets, nil
}
func (pg Postgres) DeleteOptionSet(optionSetID KountaID) error {
	_, err := pg.Exec(`DELETE FROM menu_option_sets WHERE pos_id = $1`, optionSetID)
	return err
}
