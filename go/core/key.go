package core

// Key is a data store of versioned keys
type Key struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
	Value   string `json:"key"`
}
