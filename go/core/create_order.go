package core

// CreateOrder represents a new order that has just been created by client
type CreateOrder struct {
	SiteID    KountaID              `json:"site_id"`
	MenuItems []CreateOrderMenuItem `json:"menu_items"`
}

// CreateOrderMenuItem represents a single item on a CreateOrder
type CreateOrderMenuItem struct {
	ID                  DatabaseID               `json:"id"`
	PosID               KountaID                 `json:"-"`
	Quantity            int                      `json:"quantity"`
	SelectedModifierIDs []DatabaseID             `json:"modifiers"`
	PosModifierIDs      []KountaID               `json:"-"`
	SelectedOptions     []MenuItemSelectedOption `json:"options"`
	PosSelectedOptions  []MenuItemKountaOption   `json:"-"`
}

// MenuItemSelectedOption is a selected option set modifier for a line item
type MenuItemSelectedOption struct {
	OptionSetID DatabaseID `json:"option_set_id"`
	ModifierID  DatabaseID `json:"modifier_id"`
}

// MenuItemKountaOption is an option set and modifier pair of KountaIDs
type MenuItemKountaOption struct {
	OptionSetID KountaID
	ModifierID  KountaID
}
