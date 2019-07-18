package gofortiweb

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"mime/multipart"
)

type SystemCertificatesLocalRes []struct {
	ID           string `json:"_id"`
	Name         string `json:"name"`
	CanDelete    bool   `json:"can_delete"`
	CanView      bool   `json:"can_view"`
	CanDownload  bool   `json:"can_download"`
	CanConfig    bool   `json:"can_config"`
	Comments     string `json:"comments"`
	Status       string `json:"status"`
	Issuer       string `json:"issuer"`
	ValidFrom    string `json:"validFrom"`
	ValidTo      string `json:"validTo"`
	Version      int    `json:"version"`
	Subject      string `json:"subject"`
	SerialNumber string `json:"serialNumber"`
	Extension    string `json:"extension"`
}

type SystemCertificateLocalRes struct {
	CER  string `json:"cer"`
	Key  string `json:"key"`
	Pass string `json:"pass"`
}

type SystemCertificateLocalImportRes struct {
	ID      string `json:"_id"`
	Message string `json:"msg"`
}

// SystemGetCertificatesLocal returns system local certificates
func (c *Client) SystemGetCertificatesLocal(adom string) (SystemCertificatesLocalRes, error) {

	req, err := c.NewRequest(adom, "GET", "System/Certificates/Local", nil)
	if err != nil {
		return SystemCertificatesLocalRes{}, fmt.Errorf("Failed create http request: %s", err)
	}

	get, err := c.Client.Do(req)
	if err != nil {
		return SystemCertificatesLocalRes{}, err
	}
	defer get.Body.Close()

	if get.StatusCode != 200 {
		return SystemCertificatesLocalRes{}, fmt.Errorf("Failed to get system certificates local endpoint with http code: %d", get.StatusCode)
	}

	body, err := ioutil.ReadAll(get.Body)
	if err != nil {
		return SystemCertificatesLocalRes{}, err
	}

	var systemCertificatesLocalRes SystemCertificatesLocalRes
	err = json.Unmarshal(body, &systemCertificatesLocalRes)
	if err != nil {
		return SystemCertificatesLocalRes{}, err
	}

	return systemCertificatesLocalRes, nil
}

// SystemGetCertificateLocal returns system local certificate by name
func (c *Client) SystemGetCertificateLocal(adom, name string) (SystemCertificateLocalRes, error) {

	req, err := c.NewRequest(adom, "GET", fmt.Sprintf("System/Certificates/LocalDownloadJson/%s", name), nil)
	if err != nil {
		return SystemCertificateLocalRes{}, fmt.Errorf("Failed create http request: %s", err)
	}

	get, err := c.Client.Do(req)
	if err != nil {
		return SystemCertificateLocalRes{}, err
	}
	defer get.Body.Close()

	if get.StatusCode != 200 {
		return SystemCertificateLocalRes{}, fmt.Errorf("Failed to get system certificates local endpoint with http code: %d", get.StatusCode)
	}

	body, err := ioutil.ReadAll(get.Body)
	if err != nil {
		return SystemCertificateLocalRes{}, err
	}

	var systemCertificateLocalRes SystemCertificateLocalRes
	err = json.Unmarshal(body, &systemCertificateLocalRes)
	if err != nil {
		return SystemCertificateLocalRes{}, err
	}

	return systemCertificateLocalRes, nil
}

// SystemCreateCertificateLocal creates a new local certificate
func (c *Client) SystemCreateCertificateLocal(adom string, cert, key []byte, password string) error {

	form, contentType, err := createCertificateForm(cert, key, password)
	if err != nil {
		return err
	}

	req, err := c.NewRequest(adom, "POST", "System/Certificates/Local", form)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("local certificate creation failed with status code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var importRes SystemCertificateLocalImportRes
	err = json.Unmarshal(body, &importRes)
	if err != nil {
		return err
	}

	if importRes.Message != "" {
		return fmt.Errorf("local certificate creation failed: %s", importRes.Message)
	}

	return nil
}

func createCertificateForm(cert, key []byte, password string) (form *bytes.Buffer, contentType string, err error) {

	var b bytes.Buffer

	// Parse certificate
	block, _ := pem.Decode(cert)
	pub, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return &b, contentType, err
	}

	w := multipart.NewWriter(&b)
	defer w.Close()

	w.WriteField("type", "certificate")
	w.WriteField("password", password)

	certPart, err := w.CreateFormFile("certificateFile", pub.Subject.CommonName+".crt")
	if err != nil {
		return &b, contentType, err
	}

	_, err = certPart.Write(cert)
	if err != nil {
		return &b, contentType, err
	}

	keyPart, err := w.CreateFormFile("keyFile", "tls.key")
	if err != nil {
		return &b, contentType, err
	}

	_, err = keyPart.Write(key)

	return &b, w.FormDataContentType(), err
}

// SystemDeleteCertificatesLocal deletes system local certificate by name
func (c *Client) SystemDeleteCertificatesLocal(adom, name string) error {

	req, err := c.NewRequest(adom, "DELETE", fmt.Sprintf("System/Certificates/Local/%s", name), nil)
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
