package tines_cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client - Tines API client.
type Client struct {
	TenantUrl       string
	ApiKey          string
	ProviderVersion string
	HTTPClient      *http.Client
}

// NewClient - Creates a new Tines API client.
func NewClient(tenant, apiKey, providerVersion string) (*Client, error) {
	if tenant == "" {
		return nil, Error{
			Type: ErrorTypeAuthentication,
			Errors: []ErrorMessage{
				{
					Message: "host error",
					Details: errEmptyTenant,
				},
			},
		}
	}

	if apiKey == "" {
		return nil, Error{
			Type: ErrorTypeAuthentication,
			Errors: []ErrorMessage{
				{
					Message: "credential error",
					Details: errEmptyApiKey,
				},
			},
		}
	}
	c := Client{
		TenantUrl:       tenant,
		ApiKey:          apiKey,
		ProviderVersion: providerVersion,
		HTTPClient:      &http.Client{},
	}
	return &c, nil
}

func (c *Client) doRequest(method, path string, data []byte) ([]byte, error) {
	tenant, err := url.Parse(c.TenantUrl)
	if err != nil {
		return nil, Error{
			Type: ErrorTypeRequest,
			Errors: []ErrorMessage{
				{
					Message: errParseError,
					Details: err.Error(),
				},
			},
		}
	}

	fullUrl := tenant.JoinPath(path).String()
	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer(data))
	if err != nil {
		return nil, Error{
			Type: ErrorTypeRequest,
			Errors: []ErrorMessage{
				{
					Message: errDoRequestError,
					Details: err.Error(),
				},
			},
		}
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Set("User-Agent", "tines-terraform-client")
	req.Header.Set("x-tines-client-version", fmt.Sprintf("tines-terraform-provider-%s", c.ProviderVersion))
	req.Header.Set("x-user-token", c.ApiKey)

	resp, respErr := c.HTTPClient.Do(req)
	if respErr != nil {
		return nil, Error{
			Type: ErrorTypeRequest,
			Errors: []ErrorMessage{
				{
					Message: errDoRequestError,
					Details: respErr.Error(),
				},
			},
		}
	}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, Error{
			Type:       ErrorTypeServer,
			StatusCode: resp.StatusCode,
			Errors: []ErrorMessage{
				{
					Message: errReadBodyError,
					Details: readErr.Error(),
				},
			},
		}
	}

	// Return a server error for 5XX responses
	if resp.StatusCode >= http.StatusInternalServerError {
		errMsgs := c.getErrorMessages(body)

		return nil, Error{
			Type:       ErrorTypeServer,
			StatusCode: resp.StatusCode,
			Errors:     errMsgs,
		}
	}

	// Return a request error for 4XX responses
	if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError {
		errMsgs := c.getErrorMessages(body)

		return nil, Error{
			Type:       ErrorTypeRequest,
			StatusCode: resp.StatusCode,
			Errors:     errMsgs,
		}
	}

	return body, nil
}

func (c *Client) getErrorMessages(body []byte) []ErrorMessage {
	var errorInfo Error
	var errorMsgs []ErrorMessage

	// The structure of an error response body can be inconsistent between API endpoints,
	// so we try a couple techniques to capture the error messages.
	jsonErr := json.Unmarshal(body, &errorInfo)
	if jsonErr != nil {
		jsonErr := json.Unmarshal(body, &errorMsgs)
		if jsonErr != nil && body != nil {
			errorMsgs = []ErrorMessage{{Message: "message", Details: string(body)}}
		}
		return errorMsgs
	}

	return errorInfo.Errors
}
