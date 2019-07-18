package gofortiweb

import (
	"testing"
)

func TestSystemStatus(t *testing.T) {

	client, err := NewClientHelper()
	if err != nil {
		t.Error(err)
	}

	_, err = client.SystemStatus()
	if err != nil {
		t.Error(err)
	}
}
