package tines_cli

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTinesClient(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate that the client is sending expected request values
		assert.Equal("application/json", r.Header.Get("Content-Type"), "client should send JSON data")
		assert.Equal("application/json", r.Header.Get("Accept"))
		assert.Equal("foo", r.Header.Get("x-user-token"))
		assert.Equal("tines-terraform-client", r.Header.Get("User-Agent"))
		assert.Equal("tines-terraform-provider-test", r.Header.Get("x-tines-client-version"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok")) //nolint:errcheck
	}))
	defer ts.Close()

	// Validate that the Tines CLI gets instantiated correctly
	tenant := ts.URL
	apiKey := "foo"
	version := "test"

	client, err := NewClient(tenant, apiKey, version)

	assert.Nil(err, "the Tines CLI client should instantiate successfully")

	// Validate that a generic HTTP request fires as expected
	res, err := client.doRequest("GET", "/", nil)

	assert.Equal([]byte("ok"), res, "the server should return an expected response body")
	assert.Nil(err, "the request should not return an error")
}
