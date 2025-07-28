package cyberark

import (
	"log/slog"
	"os"
	"testing"
)

const (
	subdomainAllstate = "allstate-ng"
)

func newTestClient() *Client {
	return NewClient(subdomainAllstate)
}

func TestSessionLifeRadius(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	username := os.Getenv("TEST_CYBERARK_USERNAME")
	password := os.Getenv("TEST_CYBERARK_PASSWORD")

	if username == "" || password == "" {
		t.Skip("Skipping test: TEST_CYBERARK_USERNAME or TEST_CYBERARK_PASSWORD environment variables not set")
	}

	client := newTestClient()
	session, err := client.Logon(t.Context(), LogonTypeRADIUS, username, password)
	if err != nil {
		t.Error(err)
		return
	}

	if session == nil {
		t.Error("session is nil")
		return
	}

	err = client.Logoff(t.Context(), *session)
	if err != nil {
		t.Error(err)
		return
	}
}
