package core

// Modifier represents a modification to a line item. It is used for both line modifiers and menu item modifiers
type Modifier struct {
	ID           DatabaseID `json:"id"`
	Name         string     `json:"name"`
	PriceWithTax int        `json:"price" db:"price"`
	Price        int        `json:"price_ex_tax" db:"price_ex_tax"`
	Added        bool       `json:"is_added"` // Added is true if it is an added modifier and false if removed
	SiteID       DatabaseID `json:"-"`
	PosID        KountaID   `json:"-"`
	LineID       DatabaseID `json:"-"`
	OrderID      DatabaseID `json:"-"`
}

func (modifier Modifier) Equals(other Modifier) bool {
	return (modifier.SiteID == other.SiteID) && // For when comparing menu modifiers
		(modifier.LineID == other.LineID) && // For when comparing item modifiers
		(modifier.PosID == other.PosID) &&
		(modifier.OrderID == other.OrderID) &&
		(modifier.Added == other.Added)
}
