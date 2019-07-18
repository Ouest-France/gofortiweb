package gofortiweb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type SystemStatusRes struct {
	Cluster        string `json:"cluster"`
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
	SystemUptime         string `json:"systemUptime"`
	FirmwarePartition    int    `json:"firmware_partition"`
	AdministrativeDomain string `json:"administrativeDomain"`
	Fipcc                string `json:"fipcc"`
	VMLicense            string `json:"vmLicense"`
	Registration         struct {
		Label string `json:"label"`
		URL   string `json:"url"`
		Text  string `json:"text"`
	} `json:"registration"`
	SecurityService struct {
		Expired          string `json:"expired"`
		LastUpdateTime   string `json:"lastUpdateTime"`
		LastUpdateMethod string `json:"lastUpdateMethod"`
		UpdateURL        string `json:"update_url"`
		UpdateText       string `json:"update_text"`
		BuildNumber      string `json:"buildNumber"`
	} `json:"securityService"`
	AntivirusService struct {
		AntiExpired                 string `json:"anti_expired"`
		AntivirusLastUpdateTime     string `json:"antivirusLastUpdateTime"`
		AntivirusLastUpdateMethod   string `json:"antivirusLastUpdateMethod"`
		AntiUpdateURL               string `json:"anti_update_url"`
		AntiUpdateText              string `json:"anti_update_text"`
		RegularVirusDatabaseVersion string `json:"regularVirusDatabaseVersion"`
		ExVirusDatabaseVersion      string `json:"exVirusDatabaseVersion"`
	} `json:"antivirusService"`
	ReputationService struct {
		ReputationExpired          string `json:"reputation_expired"`
		ReputationLastUpdateTime   string `json:"reputationLastUpdateTime"`
		ReputationLastUpdateMethod string `json:"reputationLastUpdateMethod"`
		ReputationUpdateURL        string `json:"reputation_update_url"`
		ReputationUpdateText       string `json:"reputation_update_text"`
		ReputationBuildNumber      string `json:"reputationBuildNumber"`
	} `json:"reputationService"`
	CredentialStuffingDefense struct {
		Expired         string `json:"expired"`
		ExpiredText     string `json:"expired_text"`
		ExpiredURL      string `json:"expired_url"`
		LastUpdateTime  string `json:"lastUpdateTime"`
		DatabaseVersion string `json:"databaseVersion"`
	} `json:"credentialStuffingDefense"`
	Readonly           bool   `json:"readonly"`
	LogDisk            string `json:"logDisk"`
	BufferSizeMax      int    `json:"bufferSizeMax"`
	FileUploadLimitMax int    `json:"fileUploadLimitMax"`
}

// SystemStatus returns system global status
func (c *Client) SystemStatus() (SystemStatusRes, error) {

	req, err := c.NewRequest("root", "GET", "System/Status/Status", nil)
	if err != nil {
		return SystemStatusRes{}, fmt.Errorf("Failed create http request: %s", err)
	}

	get, err := c.Client.Do(req)
	if err != nil {
		return SystemStatusRes{}, err
	}
	defer get.Body.Close()

	if get.StatusCode != 200 {
		return SystemStatusRes{}, fmt.Errorf("Failed to get system status endpoint with http code: %d", get.StatusCode)
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
