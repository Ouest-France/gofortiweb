package gofortiweb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type SystemCertificatesSNIMemberSlice struct {
	Results []SystemCertificatesSNIMember `json:"results"`
}

type SystemCertificatesSNIMember struct {
	ID              string `json:"id,omitempty"`
	Domain          string `json:"domain"`
	DomainType      string `json:"domain-type"`
	MultiLocalCert  string `json:"multi-local-cert"`
	LocalCert       string `json:"local-cert"`
	CertificateType string `json:"certificate-type"`
	InterGroup      string `json:"inter-group"`
}

// SystemGetCertificateSNIMembers returns SNI members certificates
func (c *Client) SystemGetCertificateSNIMembers(vdom, sni string) ([]SystemCertificatesSNIMember, error) {

	req, err := c.NewRequest(vdom, "GET", fmt.Sprintf("cmdb/system/certificate.sni/members?mkey=%s", sni), nil)
	if err != nil {
		return []SystemCertificatesSNIMember{}, fmt.Errorf("failed create http request: %s", err)
	}

	get, err := c.Client.Do(req)
	if err != nil {
		return []SystemCertificatesSNIMember{}, err
	}
	defer get.Body.Close()

	if get.StatusCode != 200 {
		return []SystemCertificatesSNIMember{}, fmt.Errorf("failed to get SNI members certificates with http code: %d", get.StatusCode)
	}

	body, err := ioutil.ReadAll(get.Body)
	if err != nil {
		return []SystemCertificatesSNIMember{}, err
	}

	var systemCertificatesSNIMembers SystemCertificatesSNIMemberSlice
	err = json.Unmarshal(body, &systemCertificatesSNIMembers)
	if err != nil {
		return []SystemCertificatesSNIMember{}, err
	}

	return systemCertificatesSNIMembers.Results, nil
}

// SystemGetCertificateSNIMember returns system SNI certificate member by id
func (c *Client) SystemGetCertificateSNIMember(vdom, sni, id string) (SystemCertificatesSNIMember, error) {

	members, err := c.SystemGetCertificateSNIMembers(vdom, sni)
	if err != nil {
		return SystemCertificatesSNIMember{}, fmt.Errorf("failed to get members list: %s", err)
	}

	for _, member := range members {
		if member.ID == id {
			return member, nil
		}
	}

	return SystemCertificatesSNIMember{}, fmt.Errorf("member %q not found in SNI %q", id, sni)
}

// SystemCreateCertificateSNIMember add a new SNI certificate member
func (c *Client) SystemCreateCertificateSNIMember(vdom, sni, domain, cert, intermediate string, domaineType string) error {

	createJSON := struct {
		Data SystemCertificatesSNIMember `json:"data"`
	}{
		Data: SystemCertificatesSNIMember{
			Domain:          domain,
			DomainType:      domaineType,
			MultiLocalCert:  "disable",
			CertificateType: "disable",
			LocalCert:       cert,
			InterGroup:      intermediate,
		},
	}

	payloadJSON, err := json.Marshal(createJSON)
	if err != nil {
		return err
	}

	payloadStr := string(payloadJSON)
	fmt.Println(payloadStr)

	req, err := c.NewRequest(vdom, "POST", fmt.Sprintf("cmdb/system/certificate.sni/members?mkey=%s", sni), bytes.NewReader(payloadJSON))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("certificate SNI member creation failed with status code: %d", res.StatusCode)
	}

	return nil
}

// SystemUpdateCertificateSNIMember update an existing SNI certificate member
func (c *Client) SystemUpdateCertificateSNIMember(vdom, sni, id, domain, cert, intermediate string, domaineType string) error {

	createJSON := struct {
		Data SystemCertificatesSNIMember `json:"data"`
	}{
		Data: SystemCertificatesSNIMember{
			Domain:          domain,
			DomainType:      domaineType,
			MultiLocalCert:  "disable",
			CertificateType: "disable",
			LocalCert:       cert,
			InterGroup:      intermediate,
		},
	}

	payloadJSON, err := json.Marshal(createJSON)
	if err != nil {
		return err
	}

	req, err := c.NewRequest(vdom, "PUT", fmt.Sprintf("cmdb/system/certificate.sni/members?mkey=%s&sub_mkey=%s", sni, id), bytes.NewReader(payloadJSON))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("certificate SNI member update failed with status code: %d", res.StatusCode)
	}

	return nil
}

// SystemDeleteCertificateSNIMember deletes a SNI certificate member
func (c *Client) SystemDeleteCertificateSNIMember(vdom, sni, id string) error {

	req, err := c.NewRequest(vdom, "DELETE", fmt.Sprintf("cmdb/system/certificate.sni/members?mkey=%s&sub_mkey=%s", sni, id), nil)
	if err != nil {
		return fmt.Errorf("failed delete http request: %s", err)
	}

	get, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer get.Body.Close()

	if get.StatusCode != 200 {
		return fmt.Errorf("failed to delete system local certificate local with http code: %d", get.StatusCode)
	}

	return nil
}
