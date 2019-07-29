package gofortiweb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type SystemCertificatesSNIMember struct {
	ID                  string `json:"id"`
	Domain              string `json:"domain"`
	DomainType          int    `json:"domainType"`
	LocalCertificate    string `json:"localCertificate"`
	CertificateVerify   string `json:"certificateVerify"`
	IntermediateCAGroup string `json:"intermediateCAGroup"`
}

// SystemGetCertificateSNIMembers returns SNI members certificates
func (c *Client) SystemGetCertificateSNIMembers(adom, sni string) ([]SystemCertificatesSNIMember, error) {

	req, err := c.NewRequest(adom, "GET", fmt.Sprintf("System/Certificates/SNI/%s/SniServerNameIndicationMember", sni), nil)
	if err != nil {
		return []SystemCertificatesSNIMember{}, fmt.Errorf("Failed create http request: %s", err)
	}

	get, err := c.Client.Do(req)
	if err != nil {
		return []SystemCertificatesSNIMember{}, err
	}
	defer get.Body.Close()

	if get.StatusCode != 200 {
		return []SystemCertificatesSNIMember{}, fmt.Errorf("Failed to get SNI members certificates with http code: %d", get.StatusCode)
	}

	body, err := ioutil.ReadAll(get.Body)
	if err != nil {
		return []SystemCertificatesSNIMember{}, err
	}

	var systemCertificatesSNIMembers []SystemCertificatesSNIMember
	err = json.Unmarshal(body, &systemCertificatesSNIMembers)
	if err != nil {
		return []SystemCertificatesSNIMember{}, err
	}

	return systemCertificatesSNIMembers, nil
}

// SystemGetCertificateSNIMember returns system SNI certificate member by id
func (c *Client) SystemGetCertificateSNIMember(adom, sni, id string) (SystemCertificatesSNIMember, error) {

	members, err := c.SystemGetCertificateSNIMembers(adom, sni)
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
func (c *Client) SystemCreateCertificateSNIMember(adom, sni, domain, cert, intermediate string, domaineType int) error {

	createJSON := struct {
		DomainType          int    `json:"domainType"`
		LocalCertificate    string `json:"localCertificate"`
		IntermediateCAGroup string `json:"intermediateCAGroup,omitempty"`
		CertificateVerify   string `json:"certificateVerify,omitempty"`
		Domain              string `json:"domain"`
	}{
		DomainType:          domaineType,
		Domain:              domain,
		LocalCertificate:    cert,
		IntermediateCAGroup: intermediate,
	}

	payloadJSON, err := json.Marshal(createJSON)
	if err != nil {
		return err
	}

	req, err := c.NewRequest(adom, "POST", fmt.Sprintf("System/Certificates/SNI/%s/SniServerNameIndicationMember", sni), bytes.NewReader(payloadJSON))
	if err != nil {
		return err
	}

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
func (c *Client) SystemUpdateCertificateSNIMember(adom, sni, id, domain, cert, intermediate string, domaineType int) error {

	createJSON := struct {
		DomainType          int    `json:"domainType"`
		LocalCertificate    string `json:"localCertificate"`
		IntermediateCAGroup string `json:"intermediateCAGroup,omitempty"`
		CertificateVerify   string `json:"certificateVerify,omitempty"`
		Domain              string `json:"domain"`
	}{
		DomainType:          domaineType,
		Domain:              domain,
		LocalCertificate:    cert,
		IntermediateCAGroup: intermediate,
	}

	payloadJSON, err := json.Marshal(createJSON)
	if err != nil {
		return err
	}

	req, err := c.NewRequest(adom, "PUT", fmt.Sprintf("System/Certificates/SNI/%s/SniServerNameIndicationMember/%s", sni, id), bytes.NewReader(payloadJSON))
	if err != nil {
		return err
	}

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
func (c *Client) SystemDeleteCertificateSNIMember(adom, sni, id string) error {

	req, err := c.NewRequest(adom, "DELETE", fmt.Sprintf("System/Certificates/SNI/%s/SniServerNameIndicationMember/%s", sni, id), nil)
	if err != nil {
		return fmt.Errorf("Failed delete http request: %s", err)
	}

	get, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer get.Body.Close()

	if get.StatusCode != 200 {
		return fmt.Errorf("Failed to delete system local certificate local with http code: %d", get.StatusCode)
	}

	return nil
}
