package tines_cli

import (
	"bytes"
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
func NewClient(tenant, apiKey, providerVersion *string) (*Client, error) {
	c := Client{
		TenantUrl:       *tenant,
		ApiKey:          *apiKey,
		ProviderVersion: *providerVersion,
		HTTPClient:      &http.Client{},
	}
	return &c, nil
}

func (c *Client) doRequest(method, path string, data []byte) (int, []byte, error) {
	tenant, err := url.Parse(c.TenantUrl)
	if err != nil {
		return 0, nil, err
	}

	fullUrl := tenant.JoinPath(path).String()
	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer(data))
	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Set("User-Agent", "tines-terraform-client")
	req.Header.Set("x-tines-client-version", fmt.Sprintf("tines-terraform-provider-%s", c.ProviderVersion))
	req.Header.Set("x-user-token", c.ApiKey)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, err
	}

	if res.StatusCode > 499 {
		return res.StatusCode, nil, fmt.Errorf("HTTP %d response: %s", res.StatusCode, body)
	}

	return res.StatusCode, body, err
}
