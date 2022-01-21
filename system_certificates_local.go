package gofortiweb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
)

type SystemCertificatesLocalSliceRes struct {
	Results []SystemCertificatesLocalRes `json:"results"`
}

type SystemCertificatesLocalRes struct {
	ID           string `json:"_id"`
	Name         string `json:"name"`
	QRef         int    `json:"q_ref"`
	CertType     string `json:"cert_type"`
	PkeyType     int    `json:"pkey_type"`
	CanDelete    bool   `json:"can_delete"`
	CanView      bool   `json:"can_view"`
	CanDownload  bool   `json:"can_download"`
	CanConfig    bool   `json:"can_config"`
	IsDefault    bool   `json:"is_default"`
	Hsm          string `json:"hsm"`
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

type SystemCertificateLocalImportRes struct {
	ID      string `json:"_id"`
	Message string `json:"msg"`
}

// SystemGetCertificatesLocal returns system local certificates
func (c *Client) SystemGetCertificatesLocal(vdom string) ([]SystemCertificatesLocalRes, error) {

	req, err := c.NewRequest(vdom, "GET", "system/certificate.local", nil)
	if err != nil {
		return []SystemCertificatesLocalRes{}, fmt.Errorf("failed create http request: %s", err)
	}

	get, err := c.Client.Do(req)
	if err != nil {
		return []SystemCertificatesLocalRes{}, err
	}
	defer get.Body.Close()

	if get.StatusCode != 200 {
		return []SystemCertificatesLocalRes{}, fmt.Errorf("failed to get system certificates local endpoint with http code: %d", get.StatusCode)
	}

	body, err := ioutil.ReadAll(get.Body)
	if err != nil {
		return []SystemCertificatesLocalRes{}, err
	}

	var systemCertificatesLocalSliceRes SystemCertificatesLocalSliceRes
	err = json.Unmarshal(body, &systemCertificatesLocalSliceRes)
	if err != nil {
		return []SystemCertificatesLocalRes{}, err
	}

	return systemCertificatesLocalSliceRes.Results, nil
}

// SystemCreateCertificateLocal creates a new local certificate
func (c *Client) SystemCreateCertificateLocal(vdom string, name string, cert, key []byte, password string) error {

	form, contentType, err := createCertificateForm(name, cert, key, password)
	if err != nil {
		return err
	}

	req, err := c.NewRequest(vdom, "POST", "system/certificate.local.import_certificate", form)
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

func createCertificateForm(name string, cert, key []byte, password string) (form *bytes.Buffer, contentType string, err error) {

	var b bytes.Buffer

	w := multipart.NewWriter(&b)
	defer w.Close()

	w.WriteField("type", "certificate")
	w.WriteField("password", password)

	certPart, err := w.CreateFormFile("certificateFile", name+".crt")
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
func (c *Client) SystemDeleteCertificatesLocal(vdom, name string) error {

	req, err := c.NewRequest(vdom, "DELETE", fmt.Sprintf("cmdb/system/certificate.local?mkey=%s", name), nil)
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
