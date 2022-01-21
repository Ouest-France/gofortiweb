package gofortiweb

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"testing"
	"time"
)

func TestSystemGetCertificatesLocal(t *testing.T) {

	client, err := NewClientHelper()
	if err != nil {
		t.Error(err)
	}

	_, err = client.SystemGetCertificatesLocal(os.Getenv("GOFORTIWEB_VDOM"))
	if err != nil {
		t.Error(err)
	}
}

func TestSystemCreateCertificateLocal(t *testing.T) {

	client, err := NewClientHelper()
	if err != nil {
		t.Fatal(err)
	}

	cert, key, err := generateCertificate()
	if err != nil {
		t.Fatal(err)
	}

	err = client.SystemCreateCertificateLocal(os.Getenv("GOFORTIWEB_VDOM"), "gofortiweb.test.com", cert, key, "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSystemDeleteCertificateLocal(t *testing.T) {

	client, err := NewClientHelper()
	if err != nil {
		t.Fatal(err)
	}

	err = client.SystemDeleteCertificatesLocal(os.Getenv("GOFORTIWEB_VDOM"), "gofortiweb.test.com")
	if err != nil {
		t.Fatal(err)
	}
}

func generateCertificate() (cert, key []byte, err error) {

	rawKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return cert, key, err
	}

	template := x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{"Gofortiweb Test"},
			CommonName:   "gofortiweb.test.com",
		},
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Minute * 30),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"gofortiweb.test.com"},
	}

	rawCert, err := x509.CreateCertificate(rand.Reader, &template, &template, &rawKey.PublicKey, rawKey)
	if err != nil {
		return cert, key, err
	}

	cert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rawCert})
	key = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rawKey)})

	os.WriteFile("/home/pablo/test.crt", cert, 0x777)
	os.WriteFile("/home/pablo/test.key", key, 0x777)

	return cert, key, nil
}
