package core

import "time"

type Menu struct {
	Categories []Category `json:"categories"`
	UpdatedAt  time.Time  `json:"updated_at"`
	SiteID     DatabaseID `json:"-"`
	SitePosID  KountaID   `json:"site_id"` // SitePosID is returned to client's instead of the database ID
}
