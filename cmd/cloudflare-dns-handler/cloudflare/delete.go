package cloudflare

import (
	"fmt"
	"net/http"
)

// DELETE request
const deleteDNSRecord = "https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s"

func (c *cloudflare) prepareDeleteRequest(zoneId string, recordId string) (*http.Request, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf(deleteDNSRecord, zoneId, recordId), nil)
	if err != nil {
		return req, err
	}

	c.applyHeaders(req)
	return req, nil
}

func (c *cloudflare) delete(zoneId string, recordId string) error {
	req, err := c.prepareDeleteRequest(zoneId, recordId)
	if err != nil {
		return err
	}

	body, err := do(req)
	if err != nil {
		return err
	}

	_, err = handleGenericResult(body)
	if err != nil {
		return err
	}

	return nil
}
