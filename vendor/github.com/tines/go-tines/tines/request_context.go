package tines

import (
	"context"
	"io"
	"net/http"
)

func newRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	// fmt.Printf("%+v", body)
	// fmt.Printf("%+v", url)
	return http.NewRequestWithContext(ctx, method, url, body)
}
