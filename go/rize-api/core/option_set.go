package core

type OptionSet struct {
	ID           DatabaseID `json:"id"`
	PosID        KountaID   `json:"-"`
	Name         string     `json:"name"`
	MinSelection int        `json:"min_selection"`
	MaxSelection int        `json:"max_selection"`
	Options      []Modifier `json:"options"`
}
