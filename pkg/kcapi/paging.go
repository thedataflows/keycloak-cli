package kcapi

import (
	"context"
	"encoding/json"
	"maps"
	"strconv"
)

// DefaultPageSize is the page size ListAll requests from list operations via
// the first/max query parameters.
const DefaultPageSize = 100

// ListAll pages through a list operation, following first/max, until a
// short page is returned. Call must target a collection operation.
//
// Each page request carries a clone of call.Params with first and max filled
// in; the caller's Params are never mutated. Invoke errors propagate as-is.
// An empty or non-array page body ends the walk, returning whatever has
// accumulated so far (nil when nothing did).
func (c *Client) ListAll(ctx context.Context, call Call) ([]json.RawMessage, error) {
	var all []json.RawMessage
	for first := 0; ; first += DefaultPageSize {
		page := call
		params := make(P, len(call.Params)+2)
		maps.Copy(params, call.Params)
		params["first"] = strconv.Itoa(first)
		params["max"] = strconv.Itoa(DefaultPageSize)
		page.Params = params

		raw, err := c.Invoke(ctx, page)
		if err != nil {
			return nil, err
		}
		var items []json.RawMessage
		if raw == nil || json.Unmarshal(raw, &items) != nil || items == nil {
			return all, nil // empty body or non-array result: nothing more to page
		}
		all = append(all, items...)
		if len(items) < DefaultPageSize {
			return all, nil
		}
	}
}
