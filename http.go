package psn

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var baseURL = must(url.Parse("https://np.community.playstation.net"))

func (auth *Authenticator) request(ctx context.Context, method string, body io.Reader, url *url.URL, header http.Header, value any) error {
	req, err := http.NewRequestWithContext(ctx, method, url.String(), body)
	if err != nil {
		return fmt.Errorf("can't create new %s request: %v", method, err)
	}

	req.Header = header

	resp, err := auth.client.Do(req)
	if err != nil {
		return fmt.Errorf("can't execute %s request: %v", method, err)
	}

	defer func() {
		err = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad request")
	}

	err = json.NewDecoder(resp.Body).Decode(&value)
	if err != nil {
		return fmt.Errorf("can't decode %s request: %v", method, err)
	}

	return nil
}

func (auth *Authenticator) post(ctx context.Context, formData url.Values, url *url.URL, header http.Header, value any) error {
	return auth.request(ctx, http.MethodPost, strings.NewReader(formData.Encode()), url, header, value)
}

func (auth *Authenticator) get(ctx context.Context, url *url.URL, header http.Header, value any) error {
	return auth.request(ctx, http.MethodGet, nil, url, header, value)
}
