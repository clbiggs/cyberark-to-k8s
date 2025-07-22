package cyberark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type errorResponse struct {
	ErrorCode    string `json:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage"`
}

type ClientError struct {
	Code    string
	Message string
}

func NewClientError(jsonError errorResponse) ClientError {
	return ClientError{
		Code:    jsonError.ErrorCode,
		Message: jsonError.ErrorMessage,
	}
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

func (e *ClientError) Error() string {
	return fmt.Sprintf("Error Code %s: %s", e.Code, e.Message)
}

func (c *Client) apiURL(target string) string {
	return fmt.Sprintf(target, c.Subdomain)
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

	url := fmt.Sprintf("%s%s/Logon/", c.apiURL(passwordVaultBaseURL), logonMethod)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create logon request, error: %+w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request, error: %+w", err)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body, error: %+w", err)
		}

		sessionToken := Session(strings.Trim(string(responseBody), "\""))
		return &sessionToken, nil

	default:
	}

	return nil, nil
}
