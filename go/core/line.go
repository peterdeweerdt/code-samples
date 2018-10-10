package core

// Line represents a line item on an order.
type Line struct {
	ID               DatabaseID `json:"id"`
	PosID            KountaID   `json:"-"`
	OrderID          DatabaseID `json:"-"`
	ModifierIDs      []KountaID `json:"-"`
	Price            int        `json:"price"`
	Total            int        `json:"total"`
	TotalTax         int        `json:"total_tax"`
	ProductName      string     `json:"product_name"`
	Notes            string     `json:"notes"`
	Quantity         int        `json:"quantity"`
	AddedModifiers   []Modifier `json:"added_modifiers"`
	RemovedModifiers []Modifier `json:"removed_modifiers"`
}
