package core

type TableMap struct {
	BeaconID  string   `json:"beacon_id"`
	SiteID    KountaID `json:"site_id"`
	TableName string   `json:"table_id"` //todo: coordinate this rename with client apps
}
