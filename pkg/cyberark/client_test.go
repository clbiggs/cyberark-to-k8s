package cyberark

import (
	"log/slog"
	"os"
	"runtime"
	"strings"
	"testing"
)

const (
	subdomainAllstate = "allstate-ng"
)

var (
	globalSession *Session
	testsToRun    = os.Getenv("TEST_CYBERARK_TESTS")
)

func skipTest(testname string) bool {
	if testsToRun == "*" || strings.Contains(strings.ToLower(testsToRun), strings.ToLower(testname)) {
		return false
	}
	return true
}

func functionName() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	shortname := frame.Function[strings.Index(frame.Function, ".")+1:]
	return shortname
}

func newTestClient() *Client {
	return NewClient(subdomainAllstate)
}

func TestSessionLifeRadius(t *testing.T) {
	fname := functionName()
	if skipTest(fname) {
		t.Skipf("Skipping test: %s not in list to run", fname)
	}

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

func TestGetAccountList(t *testing.T) {
	fname := functionName()
	if skipTest(fname) {
		t.Skipf("Skipping test: %s not in list to run", fname)
	}

	slog.SetLogLoggerLevel(slog.LevelDebug)

	client := newTestClient()

	if globalSession == nil {
		var err error
		username := os.Getenv("TEST_CYBERARK_USERNAME")
		password := os.Getenv("TEST_CYBERARK_PASSWORD")

		if username == "" || password == "" {
			t.Skip("Skipping test: TEST_CYBERARK_USERNAME or TEST_CYBERARK_PASSWORD environment variables not set")
		}
		globalSession, err = client.Logon(t.Context(), LogonTypeRADIUS, username, password)
		if err != nil {
			t.Error(err)
			return
		}
	}

	if globalSession == nil {
		t.Error("session is nil")
		return
	}

	accountList, err := client.RetrieveAccounts(t.Context(), *globalSession, nil, "", "", nil, nil, "", "")
	if err != nil {
		t.Error(err)
		return
	}

	if len(accountList) == 0 {
		t.Error("no accounts found")
		return
	}
}

func TestGetAccount(t *testing.T) {
	fname := functionName()
	if skipTest(fname) {
		t.Skipf("Skipping test: %s not in list to run", fname)
	}

	slog.SetLogLoggerLevel(slog.LevelDebug)
	accountID := os.Getenv("TEST_CYBERARK_ACCOUNTID")
	if accountID == "" {
		t.Skip("Skipping test: TEST_CYBERARK_ACCOUNTID environment variables not set")
	}

	client := newTestClient()

	if globalSession == nil {
		var err error
		username := os.Getenv("TEST_CYBERARK_USERNAME")
		password := os.Getenv("TEST_CYBERARK_PASSWORD")

		if username == "" || password == "" {
			t.Skip("Skipping test: TEST_CYBERARK_USERNAME or TEST_CYBERARK_PASSWORD environment variables not set")
		}
		globalSession, err = client.Logon(t.Context(), LogonTypeRADIUS, username, password)
		if err != nil {
			t.Error(err)
			return
		}
	}

	if globalSession == nil {
		t.Error("session is nil")
		return
	}

	account, err := client.RetrieveAccount(t.Context(), *globalSession, accountID)
	if err != nil {
		t.Error(err)
		return
	}

	if account == nil {
		t.Error("no account found")
		return
	}
}

func TestGetPassword(t *testing.T) {
	fname := functionName()
	if skipTest(fname) {
		t.Skipf("Skipping test: %s not in list to run", fname)
	}

	slog.SetLogLoggerLevel(slog.LevelDebug)
	accountID := os.Getenv("TEST_CYBERARK_ACCOUNTID")
	if accountID == "" {
		t.Skip("Skipping test: TEST_CYBERARK_ACCOUNTID environment variables not set")
	}

	client := newTestClient()

	if globalSession == nil {
		var err error
		username := os.Getenv("TEST_CYBERARK_USERNAME")
		password := os.Getenv("TEST_CYBERARK_PASSWORD")

		if username == "" || password == "" {
			t.Skip("Skipping test: TEST_CYBERARK_USERNAME or TEST_CYBERARK_PASSWORD environment variables not set")
		}
		globalSession, err = client.Logon(t.Context(), LogonTypeRADIUS, username, password)
		if err != nil {
			t.Error(err)
			return
		}
	}

	if globalSession == nil {
		t.Error("session is nil")
		return
	}

	password, err := client.RetrievePassword(t.Context(), *globalSession, accountID, "Unit Test")
	if err != nil {
		t.Error(err)
		return
	}

	if password == nil || password.Password == "" {
		t.Error("no account found")
		return
	}
}
