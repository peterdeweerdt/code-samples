package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func (app AppContext) TestInsertSite(t *testing.T, site *Site) {
	err := app.DB.InsertSite(site)
	assert.NoError(t, err)
}

const (
	TestSitePosID      = 123
	TestCategory1PosID = 234
	TestCategory2PosID = 235
	TestCategory3PosID = 236
)

// TestInsertMenu inserts a menu matching pos.MockKounta.GetMenuForSite() into the database
func (app AppContext) TestInsertMenu(t *testing.T) {
	site := Site{PosID: TestSitePosID, Name: "Test Site 1"}
	if err := app.DB.InsertSite(&site); err != nil {
		t.Fatalf("err: %s", err)
	}

	category1 := Category{SiteID: site.ID, PosID: TestCategory1PosID, Name: "Test Category 1"}
	if err := app.DB.UpsertCategory(&category1); err != nil {
		t.Fatalf("err: %s", err)
	}
	app.TestMarkCategoryClientFacing(t, &category1, true)
	app.TestMarkCategoryInstoreOnly(t, &category1, false)

	category2 := Category{SiteID: site.ID, PosID: TestCategory2PosID, Name: "Test Category 2"}
	if err := app.DB.UpsertCategory(&category2); err != nil {
		t.Fatalf("err: %s", err)
	}
	app.TestMarkCategoryClientFacing(t, &category2, true)
	app.TestMarkCategoryInstoreOnly(t, &category2, true)

	category3 := Category{SiteID: site.ID, PosID: TestCategory3PosID, Name: "Test Category 3"}
	if err := app.DB.UpsertCategory(&category3); err != nil {
		t.Fatalf("err: %s", err)
	}
	app.TestMarkCategoryClientFacing(t, &category3, false)
	app.TestMarkCategoryInstoreOnly(t, &category3, false)

	menuItem1 := MenuItem{SiteID: site.ID, PosID: 345, Name: "Test Menu Item 1", CategoryID: category1.ID}
	if err := app.DB.UpsertMenuItem(&menuItem1); err != nil {
		t.Fatalf("err: %s", err)
	}

	menuItem2 := MenuItem{SiteID: site.ID, PosID: 346, Name: "Test Menu Item 2", CategoryID: category1.ID}
	if err := app.DB.UpsertMenuItem(&menuItem2); err != nil {
		t.Fatalf("err: %s", err)
	}

	menuItem3 := MenuItem{SiteID: site.ID, PosID: 347, Name: "Test Menu Item 3", CategoryID: category2.ID}
	if err := app.DB.UpsertMenuItem(&menuItem3); err != nil {
		t.Fatalf("err: %s", err)
	}

	menuItem4 := MenuItem{SiteID: site.ID, PosID: 348, Name: "Test Menu Item 4", CategoryID: category2.ID}
	if err := app.DB.UpsertMenuItem(&menuItem4); err != nil {
		t.Fatalf("err: %s", err)
	}

	menuItem5 := MenuItem{SiteID: site.ID, PosID: 349, Name: "Test Menu Item 5", CategoryID: category2.ID}
	if err := app.DB.UpsertMenuItem(&menuItem5); err != nil {
		t.Fatalf("err: %s", err)
	}

	optionSet1 := OptionSet{PosID: 567, Name: "Test Option Set 1", MinSelection: 1, MaxSelection: 3}
	if err := app.DB.UpsertOptionSet(&menuItem3, &optionSet1); err != nil {
		t.Fatalf("err: %s", err)
	}

	modifier1 := Modifier{SiteID: site.ID, PosID: 456, Name: "Test Modifier 1"}
	if err := app.DB.UpsertMenuItemModifier(&menuItem1, &modifier1); err != nil {
		t.Fatalf("err: %s", err)
	}
	if err := app.DB.UpsertMenuItemModifier(&menuItem2, &modifier1); err != nil {
		t.Fatalf("err: %s", err)
	}

	modifier2 := Modifier{SiteID: site.ID, PosID: 457, Name: "Test Modifier 2"}
	if err := app.DB.UpsertMenuItemModifier(&menuItem4, &modifier2); err != nil {
		t.Fatalf("err: %s", err)
	}
	if err := app.DB.UpsertOptionSetModifier(&optionSet1, &modifier2); err != nil {
		t.Fatalf("err: %s", err)
	}

	modifier3 := Modifier{SiteID: site.ID, PosID: 458, Name: "Test Modifier 3"}
	if err := app.DB.UpsertOptionSetModifier(&optionSet1, &modifier3); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// TestInsertOrder inserts an order into the database
func (app AppContext) TestInsertOrder(t *testing.T, order *Order) {
	if err := app.DB.InsertOrder(order); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// TestMarkCategoryClientFacing will mark a category as client facing (or not) in database
func (app AppContext) TestMarkCategoryClientFacing(t *testing.T, category *Category, clientFacing bool) {
	// if testing with PG DB
	if pg, ok := app.DB.(Postgres); ok {
		_, err := pg.Exec(
			`UPDATE menu_categories
			SET client_facing = $1
			WHERE pos_id = $2`,
			clientFacing,
			category.PosID)
		if err != nil {
			t.Fatal("test mark category client facing")
		}
		return
	}

	if memoryDB, ok := app.DB.(*MemoryDB); ok {
		if categoryToUpdate, ok := memoryDB.Categories[category.ID]; ok {
			categoryToUpdate.ClientFacing = clientFacing
			memoryDB.Categories[category.ID] = categoryToUpdate
		}
		return
	}
}

// TestMarkCategoryInstoreOnly will mark a category as instore-only (or not) in database
func (app AppContext) TestMarkCategoryInstoreOnly(t *testing.T, category *Category, instoreOnly bool) {
	// if testing with PG DB
	if pg, ok := app.DB.(Postgres); ok {
		_, err := pg.Exec(
			`UPDATE menu_categories
			SET instore_only = $1
			WHERE pos_id = $2`,
			instoreOnly,
			category.PosID)
		if err != nil {
			t.Fatal("test mark category instore only")
		}
		return
	}

	if memoryDB, ok := app.DB.(*MemoryDB); ok {
		if categoryToUpdate, ok := memoryDB.Categories[category.ID]; ok {
			categoryToUpdate.InstoreOnly = instoreOnly
			memoryDB.Categories[category.ID] = categoryToUpdate
		}
		return
	}
}

// NewExpectedOrder returns an order matching pos.NewMockKountaOrder()
func NewExpectedOrder() *Order {
	return &Order{
		ID:        1,
		PosID:     789,
		Status:    OrderStatusOnHold,
		TableName: "7",
		Total:     1500,
		TotalTax:  120,
		Lines: []Line{
			{
				ID:          1,
				PosID:       345,
				ProductName: "Test Line 1",
				Notes:       "Test Notes 1",
				AddedModifiers: []Modifier{
					{
						Added:   true,
						PosID:   456,
						LineID:  1,
						OrderID: 1,
					},
				},
			},
			{
				ID:          2,
				PosID:       346,
				ProductName: "Test Line 2",
				Notes:       "Test Notes 2",
				RemovedModifiers: []Modifier{
					{
						Added:   false,
						PosID:   456,
						LineID:  2,
						OrderID: 1,
					},
				},
			},
		},
		PagerNumber: "765",
		SiteID:      TestSitePosID,
	}
}
