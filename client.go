package gofortiweb

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
)

// Client represents a Forti API client instance
type Client struct {
	Client   *http.Client
	Address  string
	Username string
	Password string
}

func (c *Client) NewRequest(vdom string, method string, path string, body io.Reader) (*http.Request, error) {

	req, err := http.NewRequest(method, fmt.Sprintf("%s/api/v2.0/%s", c.Address, path), body)
	if err != nil {
		return nil, err
	}

	creds := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Vdom     string `json:"vdom"`
	}{
		Username: c.Username,
		Password: c.Password,
		Vdom:     vdom,
	}
	jsonCreds, err := json.Marshal(creds)
	if err != nil {
		return nil, err
	}
	b64Creds := base64.StdEncoding.EncodeToString(jsonCreds)

	req.Header.Add("Authorization", b64Creds)
	req.Header.Add("Accept", "application/json")

	return req, err
}

func NewClientHelper() (*Client, error) {
	cookieJar, _ := cookiejar.New(nil)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := &http.Client{
		Jar:       cookieJar,
		Transport: tr,
	}

	fortiClient := &Client{
		Client:   httpClient,
		Address:  os.Getenv("GOFORTIWEB_ADDRESS"),
		Username: os.Getenv("GOFORTIWEB_USERNAME"),
		Password: os.Getenv("GOFORTIWEB_PASSWORD"),
	}

	return fortiClient, nil
}
