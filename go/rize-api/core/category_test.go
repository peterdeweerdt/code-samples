package core_test

import (
	"testing"

	"core"
	"github.com/stretchr/testify/assert"
)

func TestGetCategoriesForSite(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()
	app.TestInsertMenu(t)

	// act
	categories, err := app.GetCategoriesForSite(core.TestSitePosID)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, 3, len(categories))
	assert.Equal(t, "Test Category 1", categories[0].Name)
	assert.Equal(t, "Test Category 2", categories[1].Name)
	assert.False(t, categories[0].InstoreOnly)
	assert.True(t, categories[1].InstoreOnly)

	menuItems1 := categories[0].MenuItems
	assert.Equal(t, 2, len(menuItems1))
	assert.Equal(t, "Test Menu Item 1", menuItems1[0].Name)
	assert.Equal(t, "Test Menu Item 2", menuItems1[1].Name)
}
