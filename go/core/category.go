package core

import (
	"github.com/pkg/errors"
)

// Rize model of a category of menu items
type Category struct {
	ID           DatabaseID `json:"id"`
	SiteID       DatabaseID `json:"-"`       // we only return SitePosID to clients
	SitePosID    KountaID   `json:"site_id"` // computed field, don't save in DB
	PosID        KountaID   `json:"-"`
	Name         string     `json:"name"`
	MenuItems    []MenuItem `json:"menu_items"`
	ClientFacing bool       `json:"-"`
	InstoreOnly  bool       `json:"instore_only"`
}

func (app AppContext) GetCategoriesForSite(siteID KountaID) ([]Category, error) {
	categories, err := app.DB.SelectCategoriesBySiteID(siteID)
	if err != nil {
		return nil, errors.Wrap(err, "get categories for site")
	}
	if categories == nil {
		return []Category{}, nil
	}

	for i, category := range *categories {
		(*categories)[i].SitePosID = siteID

		menuItems, err := app.DB.SelectMenuItemsByCategoryID(siteID, category.ID)
		if err != nil {
			return nil, errors.Wrap(err, "get categories for site")
		}

		for i := range *menuItems {
			(*menuItems)[i].SitePosID = siteID
		}

		(*categories)[i].MenuItems = *menuItems
	}
	return *categories, nil
}
