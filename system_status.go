package gofortiweb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type SystemStatusRes struct {
	Results struct {
		Cluster        string `json:"cluster"`
		HaMultiGroup   bool   `json:"haMultiGroup"`
		ClusterMembers []struct {
			Hostname string `json:"hostname"`
			DevSn    string `json:"dev_sn"`
			Role     string `json:"role"`
		} `json:"cluster_members"`
		SerialNumber         string `json:"serialNumber"`
		OperationMode        string `json:"operationMode"`
		HaStatus             string `json:"haStatus"`
		SystemTime           string `json:"systemTime"`
		FirmwareVersion      string `json:"firmwareVersion"`
		UpDays               string `json:"up_days"`
		UpHrs                string `json:"up_hrs"`
		UpMins               string `json:"up_mins"`
		FirmwarePartition    int    `json:"firmware_partition"`
		AdministrativeDomain string `json:"administrativeDomain"`
		VMLicense            string `json:"vmLicense"`
		Registration         struct {
			Label string `json:"label"`
			URL   string `json:"url"`
			Text  string `json:"text"`
		} `json:"registration"`
		Readonly           bool `json:"readonly"`
		BufferSizeMax      int  `json:"bufferSizeMax"`
		FileUploadLimitMax int  `json:"fileUploadLimitMax"`
	} `json:"results"`
}

// SystemStatus returns system global status
func (c *Client) SystemStatus() (SystemStatusRes, error) {

	req, err := c.NewRequest("root", "GET", "system/status.systemstatus", nil)
	if err != nil {
		return SystemStatusRes{}, fmt.Errorf("failed create http request: %s", err)
	}

	get, err := c.Client.Do(req)
	if err != nil {
		return SystemStatusRes{}, err
	}
	defer get.Body.Close()

	if get.StatusCode != 200 {
		return SystemStatusRes{}, fmt.Errorf("failed to get system status endpoint with http code: %d", get.StatusCode)
	}

	body, err := ioutil.ReadAll(get.Body)
	if err != nil {
		return SystemStatusRes{}, err
	}

	var systemStatusRes SystemStatusRes
	err = json.Unmarshal(body, &systemStatusRes)
	if err != nil {
		return SystemStatusRes{}, err
	}

	return systemStatusRes, nil
}
