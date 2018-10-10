package core

import (
	"encoding/json"
	"time"
)

// PickupDetails represents a collection of information needed when creating a pickup order.
type PickupDetails struct {
	CustomerName string     `json:"customer_name"`
	PhoneNumber  string     `json:"phone_number"`
	PickupTime   *time.Time `json:"-"`
}

//
//// PickupTime computes a time.Time from supplied json string PickupTimeString
//func (p PickupDetails) PickupTime() (time.Time, error) {
//	return time.Parse(time.RFC3339, p.PickupTimeString)
//}

func (p *PickupDetails) UnmarshalJSON(data []byte) error {
	type Alias PickupDetails
	auxiliary := struct {
		PickupTimeString string `json:"pickup_time"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &auxiliary); err != nil {
		return err
	}
	pickupTime, err := time.Parse(time.RFC3339, auxiliary.PickupTimeString)
	if err != nil {
		return err
	}
	p.PickupTime = &pickupTime
	return nil
}
