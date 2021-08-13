package cloudflare

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/letsencrypt/sre-tools/cloudflare-dns-handler/config"
)

var (
	ServerErr                       = errors.New("cloudflare api returned a 500 status")
	FailedToParseGeneralResponseErr = errors.New("failed to parse general response")
	ApiReturnedErrorsErr            = errors.New("cloudflare api returned errors")
)

type GenericResult struct {
	Result  json.RawMessage `json:"result"`
	Success bool            `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

type Record struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Ttl     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

type cloudflare struct {
	apiToken string
}

func Init(apiToken string) cloudflare {
	return cloudflare{apiToken: apiToken}
}

// Apply applies the configuration specified in the config supplied by the
// caller. This function will ensure that certain DNS records exist, and certain
// DNS records don't exist.
func (c *cloudflare) Apply(conf config.Configuration) error {
	listResult, err := c.list(conf.ZoneIdentifier)
	if err != nil {
		return err
	}

	for _, record := range conf.Records {
		result := find(listResult, record)
		fmt.Printf("Checking if %s with value of %s exists\n", record.Name, record.Content)

		if record.Exist { // We want the record to exist
			if result == nil { // Record doesn't already exist
				fmt.Printf("\tAttempting to create the record.\n")
				err = c.handleCreation(conf.ZoneIdentifier, record)
				if err != nil {
					log.Fatalf("\tError when creating the DNS record: %v\n", err)

					// Unreachable
					return err
				}
			} else {
				fmt.Printf("\tRecord exists, continuing to the next one.\n")
				continue
			}
		} else { // We don't want the record to exist
			if result != nil { // Record does exist
				fmt.Printf("\tAttempting to delete the record.\n")
				err = c.handleDeletion(conf.ZoneIdentifier, result)
				if err != nil {
					log.Fatalf("\tError when deleting the DNS record: %v\n", err)

					// Unreachable
					return err
				}
			} else {
				fmt.Printf("\tRecord does not exist, continuing to the next one.\n")
			}
		}
	}

	return nil
}

func (c *cloudflare) handleDeletion(zoneId string, record *ListResult) error {
	return c.delete(zoneId, record.Id)
}

func (c *cloudflare) handleCreation(zoneId string, record config.Record) error {
	return c.create(zoneId, Record{
		Type:    record.Type,
		Name:    record.Name,
		Content: record.Content,
		Ttl:     record.Ttl,
	})
}

func find(existing []ListResult, want config.Record) *ListResult {
	for _, record := range existing {
		if record.Name == want.Name &&
			record.Content == want.Content &&
			record.Ttl == want.Ttl &&
			record.Type == want.Type &&
			record.Proxied == want.Proxied {
			return &record
		}
	}
	return nil
}

func handleGenericResult(body []byte) (GenericResult, error) {
	var result GenericResult
	err := json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("%w: %v", FailedToParseGeneralResponseErr, err)
	}

	if result.Success {
		return result, nil
	}

	var errorContent strings.Builder
	for _, err := range result.Errors {
		_, e := fmt.Fprintf(&errorContent, "error code: %d, error message: %s\n", err.Code, err.Message)

		// This should NEVER happen
		if e != nil {
			return result, e
		}
	}

	return result, fmt.Errorf("%s: %w", errorContent.String(), ApiReturnedErrorsErr)
}

func do(r *http.Request) ([]byte, error) {
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 500 {
		return nil, ServerErr
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *cloudflare) applyHeaders(r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))
}
