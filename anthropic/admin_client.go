package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	adminBaseURL     = "https://api.anthropic.com"
	anthropicVersion = "2023-06-01"
)

// adminClient is a minimal client for the Anthropic Admin API
// (https://api.anthropic.com/v1/organizations/...), which is not covered by
// the official Go SDK.
type adminClient struct {
	apiKey     string
	httpClient *http.Client
}

func newAdminClient(apiKey string) *adminClient {
	return &adminClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

type adminListPage struct {
	Data    []json.RawMessage `json:"data"`
	HasMore bool              `json:"has_more"`
	FirstID *string           `json:"first_id"`
	LastID  *string           `json:"last_id"`
}

// getJSON performs a GET request against the Admin API and decodes into out.
func (c *adminClient) getJSON(ctx context.Context, path string, params url.Values, out interface{}) error {
	u := adminBaseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", anthropicVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("anthropic admin api error: %s %s: status %d: %s", http.MethodGet, path, resp.StatusCode, string(body))
	}
	return json.Unmarshal(body, out)
}

// pagedListPage is the response envelope for endpoints that paginate with an
// opaque next_page token (Compliance directory, projects, attachments).
type pagedListPage struct {
	Data []json.RawMessage `json:"data"`
	// HasMore is a pointer because some endpoints (e.g. analytics user
	// activity) omit it and signal the end of data via next_page alone.
	HasMore  *bool   `json:"has_more"`
	NextPage *string `json:"next_page"`
}

// listAllPaged paginates through a next_page-token list endpoint, invoking
// handle for each raw item. handle returns false to stop early.
func (c *adminClient) listAllPaged(ctx context.Context, path string, extraParams url.Values, limit int64, handle func(item json.RawMessage) (bool, error)) error {
	params := url.Values{}
	for k, vs := range extraParams {
		for _, v := range vs {
			params.Add(k, v)
		}
	}
	pageSize := int64(100)
	if limit > 0 && limit < pageSize {
		pageSize = limit
	}
	params.Set("limit", fmt.Sprintf("%d", pageSize))

	for {
		var page pagedListPage
		if err := c.getJSON(ctx, path, params, &page); err != nil {
			return err
		}
		for _, item := range page.Data {
			cont, err := handle(item)
			if err != nil {
				return err
			}
			if !cont {
				return nil
			}
		}
		if page.HasMore != nil && !*page.HasMore {
			return nil
		}
		if page.NextPage == nil || *page.NextPage == "" {
			return nil
		}
		params.Set("page", *page.NextPage)
	}
}

// analyticsBucket is one time bucket from the analytics usage/cost report
// endpoints, whose rows are nested under data[].results[].
type analyticsBucket struct {
	StartingAt string            `json:"starting_at"`
	EndingAt   string            `json:"ending_at"`
	Results    []json.RawMessage `json:"results"`
}

type analyticsBucketedPage struct {
	Data     []analyticsBucket `json:"data"`
	HasMore  bool              `json:"has_more"`
	NextPage *string           `json:"next_page"`
}

// listAnalyticsBuckets paginates a bucketed analytics report endpoint,
// invoking handle for each result row along with its bucket boundaries.
func (c *adminClient) listAnalyticsBuckets(ctx context.Context, path string, extraParams url.Values, handle func(bucket analyticsBucket, item json.RawMessage) (bool, error)) error {
	params := url.Values{}
	for k, vs := range extraParams {
		for _, v := range vs {
			params.Add(k, v)
		}
	}

	for {
		var page analyticsBucketedPage
		if err := c.getJSON(ctx, path, params, &page); err != nil {
			return err
		}
		for _, bucket := range page.Data {
			for _, item := range bucket.Results {
				cont, err := handle(bucket, item)
				if err != nil {
					return err
				}
				if !cont {
					return nil
				}
			}
		}
		if !page.HasMore || page.NextPage == nil || *page.NextPage == "" {
			return nil
		}
		params.Set("page", *page.NextPage)
	}
}

// listAll paginates through an Admin API list endpoint, invoking handle for
// each raw item. handle returns false to stop pagination early (e.g. when the
// query row limit has been reached).
func (c *adminClient) listAll(ctx context.Context, path string, limit int64, handle func(item json.RawMessage) (bool, error)) error {
	params := url.Values{}
	pageSize := int64(100)
	if limit > 0 && limit < pageSize {
		pageSize = limit
	}
	params.Set("limit", fmt.Sprintf("%d", pageSize))

	for {
		var page adminListPage
		if err := c.getJSON(ctx, path, params, &page); err != nil {
			return err
		}
		for _, item := range page.Data {
			cont, err := handle(item)
			if err != nil {
				return err
			}
			if !cont {
				return nil
			}
		}
		if !page.HasMore || page.LastID == nil {
			return nil
		}
		params.Set("after_id", *page.LastID)
	}
}
