package cloudflare

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GET request
const listDNSRecords = "https://api.cloudflare.com/client/v4/zones/%s/dns_records"

type ListResult struct {
	Record
	Id string `json:"id,omitempty"`
}

func (c *cloudflare) prepareListRequest(zoneId string) (*http.Request, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(listDNSRecords, zoneId), nil)
	if err != nil {
		return req, err
	}

	c.applyHeaders(req)
	return req, nil
}

func (c *cloudflare) list(zoneId string) ([]ListResult, error) {
	req, err := c.prepareListRequest(zoneId)
	if err != nil {
		return nil, err
	}

	body, err := do(req)
	if err != nil {
		return nil, err
	}

	genericResult, err := handleGenericResult(body)
	if err != nil {
		return nil, err
	}

	var listResult []ListResult
	err = json.Unmarshal(genericResult.Result, &listResult)
	if err != nil {
		return nil, err
	}

	return listResult, nil
}
