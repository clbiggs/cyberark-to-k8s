package cyberark

import (
	"os"
	"testing"
)

const (
	subdomainAllstate = "allstate"
)

func newTestClient() *Client {
	return NewClient(subdomainAllstate)
}

func TestLogonRadius(t *testing.T) {
	username := os.Getenv("TEST_CYBERARK_USERNAME")
	password := os.Getenv("TEST_CYBERARK_PASSWORD")

	if username == "" || password == "" {
		t.Skip("Skipping test: TEST_CYBERARK_USERNAME or TEST_CYBERARK_PASSWORD environment variables not set")
	}

	client := newTestClient()
	session, err := client.Logon(t.Context(), LogonTypeRADIUS, username, password)
	if err != nil {
		t.Error(err)
	}

	if session == nil {
		t.Error("session is nil")
	}
}
