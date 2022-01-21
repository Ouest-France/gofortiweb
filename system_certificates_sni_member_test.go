package gofortiweb

import (
	"os"
	"testing"
)

func TestSystemGetCertificateSNIMembers(t *testing.T) {

	client, err := NewClientHelper()
	if err != nil {
		t.Error(err)
	}

	_, err = client.SystemGetCertificateSNIMembers(os.Getenv("GOFORTIWEB_VDOM"), os.Getenv("GOFORTIWEB_SNI"))
	if err != nil {
		t.Error(err)
	}
}

func TestSystemGetCertificateSNIMember(t *testing.T) {

	client, err := NewClientHelper()
	if err != nil {
		t.Error(err)
	}

	_, err = client.SystemGetCertificateSNIMember(os.Getenv("GOFORTIWEB_VDOM"), os.Getenv("GOFORTIWEB_SNI"), "1")
	if err != nil {
		t.Error(err)
	}
}

func TestSystemCreateCertificateSNIMember(t *testing.T) {

	client, err := NewClientHelper()
	if err != nil {
		t.Fatal(err)
	}

	err = client.SystemCreateCertificateSNIMember(
		os.Getenv("GOFORTIWEB_VDOM"),
		os.Getenv("GOFORTIWEB_SNI"),
		"^.*\\.foobar\\.fr",
		"gofortiweb.test.com",
		"certificat_letsencrypt",
		"regular")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSystemUpdateCertificateSNIMember(t *testing.T) {

	client, err := NewClientHelper()
	if err != nil {
		t.Fatal(err)
	}

	err = client.SystemUpdateCertificateSNIMember(
		os.Getenv("GOFORTIWEB_VDOM"),
		os.Getenv("GOFORTIWEB_SNI"),
		"1",
		"^.*\\.foobar\\.fr",
		"gofortiweb.test.com",
		"certificat_letsencrypt",
		"regular")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSystemDeleteCertificateSNIMember(t *testing.T) {

	client, err := NewClientHelper()
	if err != nil {
		t.Fatal(err)
	}

	err = client.SystemDeleteCertificateSNIMember(
		os.Getenv("GOFORTIWEB_VDOM"),
		os.Getenv("GOFORTIWEB_SNI"),
		"1")
	if err != nil {
		t.Fatal(err)
	}
}
