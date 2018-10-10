package core_test

import (
	"testing"

	"core"
	"github.com/stretchr/testify/assert"
)

func TestGetSitesShouldReturnSites(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()
	site1 := core.Site{PosID: 1, MenuHash: "some_long_hash", Name: "Test Site"}
	site2 := core.Site{PosID: 2, MenuHash: "some_long_hash", Name: "Other Test Site"}
	app.TestInsertSite(t, &site1)
	app.TestInsertSite(t, &site2)

	// act
	sites, err := app.DB.SelectSites()

	// assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(*sites))
}

func TestUpdateAllMenusShouldSaveSites(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()

	// act
	app.UpdateAllMenus()

	// assert
	sites, _ := app.DB.SelectSites()
	assert.Equal(t, 1, len(*sites))
}

func TestUpdateAllMenusShouldSaveCategories(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()

	// act
	app.UpdateAllMenus()

	// assert
	categories, _ := app.GetCategoriesForSite(123)
	assert.Equal(t, 2, len(categories))
}

func TestUpdateAllMenusShouldSaveMenuItems(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()

	// act
	app.UpdateAllMenus()

	// assert
	categories, _ := app.GetCategoriesForSite(123)
	assert.Equal(t, 2, len(categories[0].MenuItems))
	assert.Equal(t, 1, len(categories[1].MenuItems))
}

func TestUpdateAllMenusIfMenusMatchShouldNotUpdateSite(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()

	app.UpdateAllMenus()
	sites, _ := app.DB.SelectSites()
	existingSiteCount := len(*sites)
	previousMenuHash := (*sites)[0].MenuHash

	// act
	app.UpdateAllMenus() // update again with identical information

	// assert
	updatedSites, _ := app.DB.SelectSites()
	assert.Equal(t, existingSiteCount, len(*updatedSites))
	assert.Equal(t, previousMenuHash, (*updatedSites)[0].MenuHash)
}

func TestUpdateAllMenusIfHashesDoNotMatchShouldUpdatesSite(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()

	initialHash := "Some old hash"
	site := core.Site{PosID: 123, MenuHash: initialHash, Name: "Test Site 1"}
	app.TestInsertSite(t, &site)

	// act
	app.UpdateAllMenus()

	// assert
	updatedSites, _ := app.DB.SelectSites()
	assert.Equal(t, 1, len(*updatedSites))
	assert.NotEqual(t, initialHash, (*updatedSites)[0].MenuHash)
}

func TestUpdateAllMenusIfHashesDoNotMatchShouldUpdatesMenu(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()

	const (
		sitePosID        = 123
		oldCategoryPosID = 222
	)

	// add a single site and menu to DB
	initialHash := "Some old hash"
	site := core.Site{PosID: sitePosID, MenuHash: initialHash, Name: "Test Site 1"}
	app.TestInsertSite(t, &site)

	oldCategory := core.Category{SiteID: site.ID, SitePosID: sitePosID, PosID: oldCategoryPosID, Name: "Test Category 1"}
	if err := app.DB.UpsertCategory(&oldCategory); err != nil {
		t.Fatalf("err: %s", err)
	}

	menuItem1 := core.MenuItem{SiteID: site.ID, PosID: 2, Name: "Test Menu Item 1", CategoryID: oldCategory.ID}
	if err := app.DB.UpsertMenuItem(&menuItem1); err != nil {
		t.Fatalf("err: %s", err)
	}

	menuItem2 := core.MenuItem{SiteID: site.ID, PosID: 3, Name: "Test Menu Item 2", CategoryID: oldCategory.ID}
	if err := app.DB.UpsertMenuItem(&menuItem2); err != nil {
		t.Fatalf("err: %s", err)
	}

	modifier1 := core.Modifier{SiteID: site.ID, PosID: 4, Name: "Test Modifier 1"}
	if err := app.DB.UpsertMenuItemModifier(&menuItem1, &modifier1); err != nil {
		t.Fatalf("err: %s", err)
	}

	modifier2 := core.Modifier{SiteID: site.ID, PosID: 5, Name: "Test Modifier 2"}
	if err := app.DB.UpsertMenuItemModifier(&menuItem1, &modifier2); err != nil {
		t.Fatalf("err: %s", err)
	}

	modifier3 := core.Modifier{SiteID: site.ID, PosID: 6, Name: "Test Modifier 3"}
	optionSet1 := core.OptionSet{PosID: 7, Name: "Test Option Set 1", MinSelection: 3, MaxSelection: 5}
	if err := app.DB.UpsertOptionSet(&menuItem2, &optionSet1); err != nil {
		t.Fatalf("err: %s", err)
	}
	if err := app.DB.UpsertOptionSetModifier(&optionSet1, &modifier3); err != nil {
		t.Fatalf("err: %s", err)
	}

	// act
	app.UpdateAllMenus() // will generate new menu hash

	// assert
	updatedCategories, _ := app.GetCategoriesForSite(sitePosID)
	assert.Equal(t, 2, len(updatedCategories))
	assert.Equal(t, 2, len(updatedCategories[0].MenuItems))
	assert.Equal(t, 1, len(updatedCategories[1].MenuItems))

	for _, c := range updatedCategories {
		if c.PosID == oldCategoryPosID {
			assert.Fail(t, "Original category should be deleted")
		}
	}

	oldItem, _ := app.DB.GetMenuItem(site.PosID, menuItem1.ID)
	assert.Nil(t, oldItem)

	oldModifier, _ := app.DB.GetMenuModifierByKountaID(sitePosID, 4)
	assert.Nil(t, oldModifier)

	newModifier, _ := app.DB.GetMenuModifierByKountaID(sitePosID, 456)
	assert.NotNil(t, newModifier)

	updatedOptionSets, _ := app.DB.SelectOptionSets()
	var oldOptionSet, newOptionSet *core.OptionSet
	for _, o := range *updatedOptionSets {
		if o.PosID == optionSet1.PosID {
			oldOptionSet = &o
		} else if o.PosID == 567 {
			newOptionSet = &o
		}
	}
	assert.Nil(t, oldOptionSet)
	assert.NotNil(t, newOptionSet)
}

func TestGetSiteByPosIDShouldReturnCorrectSite(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()

	// add a single initialSite and menu to DB
	expectedHash := "some hash"
	initialSite := core.Site{PosID: 123, MenuHash: expectedHash, Name: "Test Site 1"}
	app.TestInsertSite(t, &initialSite)

	// act
	updatedSite, err := app.DB.GetSite(123)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, initialSite.PosID, updatedSite.PosID)
	assert.Equal(t, initialSite.Name, updatedSite.Name)
	assert.Equal(t, expectedHash, updatedSite.MenuHash)
}
