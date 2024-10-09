package v1

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type GetOrganizationsByNameResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      []Organization
}

// ParseGetOrganizationsByNameResponse parses an HTTP response from a GetOrganizationsByNameWithResponse call
func ParseGetOrganizationsByNameResponse(rsp *http.Response) (*GetOrganizationsByNameResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetOrganizationsByNameResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Organization
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		if len(dest) == 1 {
			response.JSON200 = dest
		}
	}

	return response, nil
}

// GetOrganizationsByNameWithResponse request returning *GetOrganizationsByNameResponse
func (c *ClientWithResponses) GetOrganizationsByNameWithResponse(ctx context.Context, name string, reqEditors ...RequestEditorFn) (*GetOrganizationsByNameResponse, error) {
	params := GetOrganizationsParams{
		Name: &name,
	}
	rsp, err := c.GetOrganizations(ctx, &params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetOrganizationsByNameResponse(rsp)
}
