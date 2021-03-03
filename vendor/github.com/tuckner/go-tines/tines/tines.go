package tines

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

// httpClient defines an interface for an http.Client implementation so that alternative
// http Clients can be passed in for making requests
type httpClient interface {
	Do(request *http.Request) (response *http.Response, err error)
}

// A Client manages communication with the Tines API.
type Client struct {
	// HTTP client used to communicate with the API.
	client httpClient

	// Base URL for API requests.
	baseURL        *url.URL
	userToken      string
	userEmail      string
	Agent          *AgentService
	GlobalResource *GlobalResourceService
	Story          *StoryService
	Note           *NoteService
}

// NewClient returns a new Tines API client.
// If a nil httpClient is provided, http.DefaultClient will be used.
// baseURL is the HTTP endpoint of your Tines instance and should always be specified with a trailing slash.
func NewClient(httpClient httpClient, baseURL string, userEmail string, userToken string) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	// ensure the baseURL contains a trailing slash so that all paths are preserved in later calls
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	baseURL += "api/v1/"
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		client:    httpClient,
		baseURL:   parsedBaseURL,
		userEmail: userEmail,
		userToken: userToken,
	}
	c.Agent = &AgentService{client: c}
	c.GlobalResource = &GlobalResourceService{client: c}
	c.Story = &StoryService{client: c}
	c.Note = &NoteService{client: c}

	return c, nil
}

// NewRequestWithContext creates an API request.
// A relative URL can be provided in urlStr, in which case it is resolved relative to the baseURL of the Client.
// If specified, the value pointed to by body is JSON encoded and included as the request body.
func (c *Client) NewRequestWithContext(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	// Relative URLs should be specified without a preceding slash since baseURL will have the trailing slash
	rel.Path = strings.TrimLeft(rel.Path, "/")

	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := newRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-user-email", c.userEmail)
	req.Header.Set("x-user-token", c.userToken)

	return req, nil
}

// NewRequest wraps NewRequestWithContext using the background context.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	return c.NewRequestWithContext(context.Background(), method, urlStr, body)
}

// addOptions adds the parameters in opt as URL query parameters to s.  opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// Do sends an API request and returns the API response.
// The API response is JSON decoded and stored in the value pointed to by v, or returned as an error if an API error has occurred.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	httpResp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	err = CheckResponse(httpResp, req)
	if err != nil {
		// Even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return newResponse(httpResp, nil), err
	}

	if v != nil {
		// Open a NewDecoder and defer closing the reader only if there is a provided interface to decode to
		defer httpResp.Body.Close()
		err = json.NewDecoder(httpResp.Body).Decode(v)
	}

	resp := newResponse(httpResp, v)
	return resp, err
}

// CheckResponse checks the API response for errors, and returns them if present.
// A response is considered an error if it has a status code outside the 200 range.
// The caller is responsible to analyze the response body.
// The body can contain JSON (if the error is intended) or xml (sometimes Jira just failes).
func CheckResponse(r *http.Response, req *http.Request) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	err := fmt.Errorf("request failed. Please analyze the request body for more details. Status code: %d. Resp Body: %d. Req URL: %d", r.StatusCode, r.Body, req.URL)
	return err
}

// GetBaseURL will return you the Base URL.
// This is the same URL as in the NewClient constructor
func (c *Client) GetBaseURL() url.URL {
	return *c.baseURL
}

// Response represents Tines API response. It wraps http.Response returned from
// API and provides information about paging.
type Response struct {
	*http.Response

	StartAt    int
	MaxResults int
	Total      int
}

func newResponse(r *http.Response, v interface{}) *Response {
	resp := &Response{Response: r}
	// resp.populatePageValues(v)
	return resp
}

// Sets paging values if response json was parsed to searchResult type
// (can be extended with other types if they also need paging info)
// func (r *Response) populatePageValues(v interface{}) {
// 	switch value := v.(type) {
// 	case *searchResult:
// 		r.StartAt = value.StartAt
// 		r.MaxResults = value.MaxResults
// 		r.Total = value.Total
// 	case *groupMembersResult:
// 		r.StartAt = value.StartAt
// 		r.MaxResults = value.MaxResults
// 		r.Total = value.Total
// 	}
// }
