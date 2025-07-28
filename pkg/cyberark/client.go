package cyberark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/errors"
)

const (
	cyberArkHostname     = "https://%s.privilegecloud.cyberark.com/"
	passwordVaultBaseURL = cyberArkHostname + "PasswordVault/API/"
	authBaseURL          = passwordVaultBaseURL + "Auth/"
	accountsBaseURL      = cyberArkHostname + passwordVaultBaseURL + "Accounts/"
	safesBaseURL         = cyberArkHostname + passwordVaultBaseURL + "Safes/"

	defaultTimeout = 30 * time.Second
)

type Session string

type LogonType string

const (
	LogonTypeCyberArk = "CyberArk"
	LogonTypeLDAP     = "LDAP"
	LogonTypeRADIUS   = "RADIUS"
)

type errorResponse struct {
	ErrorCode    string `json:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage"`
	Message      string `json:"Message"`
}

type ClientError struct {
	Code    string
	Message string
}

func NewClientErrorFromResponse(jsonError errorResponse) *ClientError {
	msg := jsonError.ErrorMessage
	if msg == "" {
		msg = jsonError.Message
	}

	return NewClientError(jsonError.ErrorCode, msg)
}

func NewClientError(code, message string) *ClientError {
	return &ClientError{
		Code:    code,
		Message: message,
	}
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("code %s: %s", e.Code, e.Message)
}

type Client struct {
	Subdomain  string
	HTTPClient *http.Client
}

func NewClient(subdomain string) *Client {
	return &Client{
		Subdomain:  subdomain,
		HTTPClient: &http.Client{Timeout: defaultTimeout},
	}
}

func (c *Client) apiURL(target string) string {
	return fmt.Sprintf(target, c.Subdomain)
}

func (c *Client) newAuthorizedRequestWithContext(ctx context.Context, method string, url string, body io.Reader, token Session) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err == nil {
		req.Header.Add("Authorization", string(token))
	}
	return req, err
}

func (c *Client) Logon(ctx context.Context, logonMethod, username, password string) (*Session, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	reqBody := map[string]string{
		"username": username,
		"password": password,
	}

	jsonBody, _ := json.Marshal(reqBody)

	url := fmt.Sprintf("%s%s/Logon/", c.apiURL(authBaseURL), logonMethod)

	slog.Debug("attempting logon", "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create logon request, error: %+w", err)
	}
	req.Header.Add("content-type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request, error: %+w", err)
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body, error: %+w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		sessionToken := Session(strings.Trim(string(responseBody), "\""))
		return &sessionToken, nil

	default:
		slog.Debug("logon error", "response", responseBody)
		var jsonError errorResponse
		if err := json.Unmarshal(responseBody, &jsonError); err != nil {
			return nil, NewClientError(strconv.Itoa(resp.StatusCode), fmt.Sprintf("Logon failed with response from api: %s, %+v", resp.Status, err))
		}
		return nil, NewClientErrorFromResponse(jsonError)
	}
}

func (c *Client) Logoff(ctx context.Context, session Session) error {
	if session == "" {
		return nil
	}

	url := c.apiURL(authBaseURL) + "/Logoff/"

	req, err := c.newAuthorizedRequestWithContext(ctx, http.MethodPost, url, nil, session)
	if err != nil {
		return fmt.Errorf("failed to create logoff request, error: %+w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request, error: %+w", err)
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body, error: %+w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	default:
		slog.Debug("logoff error", "response", responseBody)
		var jsonError errorResponse
		if err := json.Unmarshal(responseBody, &jsonError); err != nil {
			return NewClientError(strconv.Itoa(resp.StatusCode), fmt.Sprintf("Logon failed with response from api: %s. %+v", resp.Status, err))
		}
		return NewClientErrorFromResponse(jsonError)
	}
}
