package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// POST request
const createDNSRecord = "https://api.cloudflare.com/client/v4/zones/%s/dns_records"

type DNSRecord struct {
}

func (c *cloudflare) prepareCreateRequest(zoneId string, r Record) (*http.Request, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err :=  http.NewRequest("POST", fmt.Sprintf(createDNSRecord, zoneId), bytes.NewBuffer(data))
	if err != nil {
		return req, err
	}

	c.applyHeaders(req)
	return req, nil
}
func (c *cloudflare) create(zoneId string, r Record)  error {
	req, err := c.prepareCreateRequest(zoneId, r)
	if err != nil {
		return  err
	}

	body, err := do(req)
	if err != nil {
		return  err
	}

	_, err = handleGenericResult(body)
	if err != nil {
		return  err
	}

	return  nil
}
