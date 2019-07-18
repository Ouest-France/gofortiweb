package gofortiweb

import (
	"crypto/tls"
	"encoding/base64"
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

func (c *Client) NewRequest(adom string, method string, path string, body io.Reader) (*http.Request, error) {

	req, err := http.NewRequest(method, fmt.Sprintf("%s/api/v1.0/%s", c.Address, path), body)

	if err != nil {
		return nil, err
	}

	creds := fmt.Sprintf("%s:%s:%s", c.Username, c.Password, adom)
	encodedCreds := base64.StdEncoding.EncodeToString([]byte(creds))

	req.Header.Add("Authorization", encodedCreds)

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
