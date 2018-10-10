package core

import (
	"crypto/md5"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
)

// Site is an individual restaurant
type Site struct {
	ID          DatabaseID `json:"-"`
	PosID       KountaID   `json:"id"` // PosID is the only identifier given to clients
	MenuHash    string     `json:"-"`
	Name        string     `json:"name"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Address     *string    `json:"address"`
	PhoneNumber string     `json:"phone_number"`
}

// UpdateAllMenus will update the menu for each Rize site and store it in the database
func (app AppContext) UpdateAllMenus() error {
	sites, err := app.Kounta.GetAllSites()
	if err != nil {
		return errors.Wrap(err, "update all menus")
	}

	// Get existing items, to set difference and know which to delete from database
	existingCategories, existingItems, existingModifiers, existingOptionSets, err := app.getExistingMenu()
	if err != nil {
		return errors.Wrap(err, "update all menus")
	}

	for _, site := range sites {
		menu, err := app.Kounta.GetMenuForSite(site.PosID)
		if err != nil {
			return errors.Wrap(err, "update all menus")
		}

		// compute hash BEFORE assigning DatabaseIDs
		site.MenuHash = computeHash(menu.Categories)

		// get site from DB
		existingSite, err := app.DB.GetSite(site.PosID)
		if err != nil {
			return errors.Wrap(err, "update all menus")
		}

		var siteID DatabaseID

		// save site to database
		if existingSite == nil {
			if err := app.DB.InsertSite(&site); err != nil {
				return errors.Wrap(err, "update all menus")
			}
			siteID = site.ID
		} else {
			if existingSite.MenuHash != site.MenuHash {
				if err := app.DB.UpdateSiteMenuHash(existingSite, site.MenuHash); err != nil {
					return errors.Wrap(err, "update all menus")
				}
			}
			siteID = existingSite.ID
		}

		// Store the menu
		for _, category := range menu.Categories {
			delete(existingCategories, category.PosID)

			category.SiteID = siteID
			if err := app.DB.UpsertCategory(&category); err != nil {
				return errors.Wrap(err, "add menu to database")
			}

			for _, menuItem := range category.MenuItems {
				delete(existingItems, menuItem.PosID)

				menuItem.SiteID = siteID
				menuItem.CategoryID = category.ID
				if err := app.DB.UpsertMenuItem(&menuItem); err != nil {
					return errors.Wrap(err, "add menu to database")
				}

				for _, modifier := range menuItem.Modifiers {
					delete(existingModifiers, modifier.PosID)

					modifier.SiteID = siteID
					if err := app.DB.UpsertMenuItemModifier(&menuItem, &modifier); err != nil {
						return errors.Wrap(err, "add menu to database")
					}
				}

				for _, optionSet := range menuItem.OptionSets {
					delete(existingOptionSets, optionSet.PosID)

					if err := app.DB.UpsertOptionSet(&menuItem, &optionSet); err != nil {
						return errors.Wrap(err, "add menu to database")
					}

					for _, modifier := range optionSet.Options {
						delete(existingModifiers, modifier.PosID)

						modifier.SiteID = siteID
						if err := app.DB.UpsertOptionSetModifier(&optionSet, &modifier); err != nil {
							return errors.Wrap(err, "add menu to database")
						}
					}
				}
			}
		}

	}

	// Now that we've updated the categories, items, option sets, and modifiers we retrieved from Kounta,
	// delete any remaining categories, items, option sets, and modifiers in the database.
	for id := range existingCategories {
		if err := app.DB.DeleteCategory(id); err != nil {
			return errors.Wrap(err, "add menu to database")
		}
	}
	for id := range existingItems {
		if err := app.DB.DeleteMenuItem(id); err != nil {
			return errors.Wrap(err, "add menu to database")
		}
	}
	for id := range existingModifiers {
		if err := app.DB.DeleteMenuModifier(id); err != nil {
			return errors.Wrap(err, "add menu to database")
		}
	}
	for id := range existingOptionSets {
		if err := app.DB.DeleteOptionSet(id); err != nil {
			return errors.Wrap(err, "add menu to database")
		}
	}

	log.Println("Done updating all menus")
	return nil
}

// GetMenuForSite will return a Menu struct for a given siteID, Omitting any Category that is not client facing
func (app AppContext) GetMenuForSite(siteID KountaID) (*Menu, error) {
	categories, err := app.GetCategoriesForSite(siteID)
	if err != nil {
		return nil, errors.Wrap(err, "get menu for site")
	}

	clientFacingCategories := []Category{}
	for i, v := range categories {
		if v.ClientFacing {
			clientFacingCategories = append(clientFacingCategories, categories[i])
		}
	}

	site, err := app.DB.GetSite(siteID)
	if err != nil {
		return nil, errors.Wrap(err, "get menu for site")
	}

	menu := Menu{Categories: clientFacingCategories, UpdatedAt: site.UpdatedAt, SitePosID: siteID}
	return &menu, nil
}

func (app AppContext) getExistingMenu() (categoryIDs map[KountaID]bool, menuItemIDs map[KountaID]bool, modifierIDs map[KountaID]bool, optionSetIDs map[KountaID]bool, err error) {
	existingCategories, err := app.DB.SelectCategories()
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "get existing menu")
	}
	categoryIDs = make(map[KountaID]bool, len(*existingCategories))
	for _, c := range *existingCategories {
		categoryIDs[c.PosID] = true
	}

	existingItems, err := app.DB.SelectMenuItems()
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "get existing menu")
	}
	menuItemIDs = make(map[KountaID]bool, len(*existingItems))
	for _, mi := range *existingItems {
		menuItemIDs[mi.PosID] = true
	}

	existingModifiers, err := app.DB.SelectMenuModifiers()
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "get existing menu")
	}
	modifierIDs = make(map[KountaID]bool, len(*existingModifiers))
	for _, m := range *existingModifiers {
		modifierIDs[m.PosID] = true
	}

	existingOptionSets, err := app.DB.SelectOptionSets()
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "get existing menu")
	}
	optionSetIDs = make(map[KountaID]bool, len(*existingOptionSets))
	for _, o := range *existingOptionSets {
		optionSetIDs[o.PosID] = true
	}

	return
}

func computeHash(categories []Category) string {
	humanReadableString := fmt.Sprintf("%v", categories)
	return fmt.Sprintf("%x", md5.Sum([]byte(humanReadableString)))
}
