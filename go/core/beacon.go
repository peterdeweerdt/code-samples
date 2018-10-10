package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"pjd"
)

// Beacon definition
type Beacon struct {
	FactoryID string `json:"factory_id"`
	Name      string `json:"name"`
	ConfigID  int    `json:"config_id"`
}

// BeaconConfiguration definition
type BeaconConfiguration struct {
	ID                         int         `json:"id,omitempty"`
	Name                       string      `json:"name"`
	BeaconType                 string      `json:"beacon_type"`
	TransmissionPower          interface{} `json:"transmission_power"`
	AntennaType                string      `json:"antenna_type"`
	CalibratedPower            interface{} `json:"calibrated_power"`
	EddystoneTransmissionPower interface{} `json:"eddystone_transmission_power"`
	UIDInstanceID              string      `json:"uid_instance_id"`
	UIDNamespaceID             string      `json:"uid_namespace_id"`
	TransmitUIDFrequency       string      `json:"transmit_uid_frequency"`
	TLM                        bool        `json:"tlm"`
}

// SetupNewBeaconConfigurations will create BeaconConfigurations and assign them to the associated Beacon
func SetupNewBeaconConfigurations(pg Postgres, siteID KountaID, version, power, count int) {
	configs := createConfigurations(siteID, version, power, count)
	savedConfigs := saveConfigurations(configs)
	assignConfigurations(siteID, savedConfigs)
	addTableMappings(pg, siteID, configs)

	// configs := getConfigurationsForVersion("v4")
	// deleteConfigurations(configs)
}

func deleteConfigurations(configs []BeaconConfiguration) {
	for _, c := range configs {
		doRequest("DELETE", fmt.Sprintf("https://manager.gimbal.com/api/beacon_configurations/%d", c.ID), nil)
	}
}

func createConfigurations(siteID KountaID, version, power, count int) []BeaconConfiguration {
	log.Println("Creating configurations")

	const (
		beaconType  = "Eddystone"
		antennaType = "Omnidirectional"
		namespace   = "ff000000000000000000"
		frequency   = "high"
	)

	configs := make([]BeaconConfiguration, count)

	for i := 1; i <= count; i++ {
		configs[i-1] = BeaconConfiguration{
			Name:                       fmt.Sprintf("%d_%02d_v%d", siteID, i, version),
			BeaconType:                 beaconType,
			TransmissionPower:          power,
			AntennaType:                antennaType,
			CalibratedPower:            power,
			EddystoneTransmissionPower: power,
			UIDInstanceID:              fmt.Sprintf("%08d0%03d", siteID, i),
			UIDNamespaceID:             namespace,
			TransmitUIDFrequency:       frequency,
			TLM:                        true,
		}
	}

	return configs
}

func saveConfigurations(configs []BeaconConfiguration) []BeaconConfiguration {
	log.Println("Saving configurations in Gimbal")

	savedConfigs := []BeaconConfiguration{}

	for _, c := range configs {
		log.Println("Saving:", c)

		resp := doRequest("POST", "https://manager.gimbal.com/api/beacon_configurations", c)
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(data, &c)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Saved config:", c)

		savedConfigs = append(savedConfigs, c)
	}

	return savedConfigs
}

func getConfigurationsForVersion(version string) []BeaconConfiguration {
	resp := doRequest("GET", "https://manager.gimbal.com/api/beacon_configurations", nil)

	log.Println(pjd.MustDumpResponse(resp))

	configs := []BeaconConfiguration{}
	err := json.NewDecoder(resp.Body).Decode(&configs)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(configs)

	versionConfigs := []BeaconConfiguration{}
	for _, c := range configs {
		if strings.HasSuffix(c.Name, version) {
			versionConfigs = append(versionConfigs, c)
		}
	}

	log.Println(versionConfigs)

	return versionConfigs
}

func assignConfigurations(siteID KountaID, configs []BeaconConfiguration) {
	log.Println("Assigning configurations to beacons in Gimbal")

	resp := doRequest("GET", "https://manager.gimbal.com/api/beacons", nil)

	beacons := []Beacon{}
	err := json.NewDecoder(resp.Body).Decode(&beacons)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(beacons)

	for _, beacon := range beacons {
		if strings.HasPrefix(beacon.Name, fmt.Sprintf("%d_", siteID)) {
			beacon.ConfigID = getConfigForBeacon(configs, &beacon).ID
			resp := doRequest("PUT", "https://manager.gimbal.com/api/beacons/"+beacon.FactoryID, beacon)
			if resp.StatusCode >= 300 {
				log.Println(pjd.MustDumpResponse(resp))
			}
		}
	}
}

func doRequest(method, url string, payload interface{}) *http.Response {
	var body io.Reader
	if payload != nil {
		json, _ := json.Marshal(payload)
		log.Println(string(json))
		body = bytes.NewReader(json)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", "Token token=f92211ccd5a8b2eede3fb917d50efc1d")
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	log.Println(resp.Status)

	return resp
}

func getConfigForBeacon(configs []BeaconConfiguration, b *Beacon) *BeaconConfiguration {
	for _, c := range configs {
		if strings.HasPrefix(c.Name, b.Name) {
			return &c
		}
	}
	return nil
}

func addTableMappings(pg Postgres, siteID KountaID, savedConfigs []BeaconConfiguration) {
	log.Println("Creating table mappings in Rize")

	for i, config := range savedConfigs {
		pg.InsertTableMap(&TableMap{
			BeaconID:  config.UIDNamespaceID + config.UIDInstanceID,
			SiteID:    siteID,
			TableName: strconv.Itoa(i + 1),
		})
	}
}
