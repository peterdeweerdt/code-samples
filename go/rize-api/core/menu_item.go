package core

type MenuItem struct {
	ID          DatabaseID  `json:"id"`
	SiteID      DatabaseID  `json:"-"`
	PosID       KountaID    `json:"-"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Price       int         `json:"price"`
	CategoryID  DatabaseID  `json:"-"`
	Modifiers   []Modifier  `json:"-"`
	OptionSets  []OptionSet `json:"-"`
	SitePosID   KountaID    `json:"site_id"`
}
