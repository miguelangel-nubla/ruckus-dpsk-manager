package client

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/data/dpsk"
)

type DpskService struct {
	Client *Client
}

func (rc *Client) Dpsk() *DpskService {
	return &DpskService{Client: rc}
}

func (d *DpskService) List() (dpsk.Entries, error) {
	// Create the request URL
	url := d.Client.server + "/admin/_cmdstat.jsp"

	body := fmt.Sprintf(`<ajax-request action="getstat" comp="stamgr" updater="dpsk-list.%s">
		<dpsklist/>
	</ajax-request>`, d.Client.getCurrentTimestamp())

	if d.Client.Debug {
		fmt.Println(body)
	}

	// Create the request object
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set the request headers
	req.Header.Set("X-CSRF-Token", d.Client.csrfToken)
	req.Header.Set("Cookie", d.Client.cookie)
	req.Header.Set("Content-Type", "text/xml")

	// Send the request
	resp, err := d.Client.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	entries, err := dpsk.FromXml(xmlData)
	if d.Client.Debug {
		fmt.Printf("Parsed DPSKs:\n")
		for _, dpsk := range entries {
			fmt.Printf("%v\n", dpsk)
		}
	}

	return entries, err
}

func (d *DpskService) Create(wlansvcID int, user string, dpskLen int) error {
	// Create the request URL
	url := d.Client.server + "/admin/_cmdstat.jsp"

	// Create the request body
	body := fmt.Sprintf(`<ajax-request action='docmd' checkAbility='2' updater='system.%s' comp='system'>
		<xcmd
			cmd='batch-dpsk'
			type='gen'
			num='1'
			max-num='2048'
			batch-dpsk=''
			wlansvc-id='%d'
			role-id=''
			dpsk-len='%d'
			dvlan-id=''
			user='%s'
		/>
	</ajax-request>`, d.Client.getCurrentTimestamp(), wlansvcID, dpskLen, user)

	if d.Client.Debug {
		fmt.Println(body)
	}

	// Create the request object
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the request headers
	req.Header.Set("X-CSRF-Token", d.Client.csrfToken)
	req.Header.Set("Cookie", d.Client.cookie)
	req.Header.Set("Content-Type", "text/xml")

	// Send the request
	resp, err := d.Client.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status code indicates success (e.g., 200 OK)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("create DPSK user failed with status code: %v", resp.StatusCode)
	}

	return nil
}

func (d *DpskService) Modify(dpskID int, fields map[string]string) error {
	// Create the request URL
	url := d.Client.server + "/admin/_conf.jsp"

	// Create the request body
	body := fmt.Sprintf(`<ajax-request action='updobj' updater='dpsk-list.%s'comp='dpsk-list'>
		<dpsk id='%d' name='dpsk%d' IS_PARTIAL='true' %s />
	</ajax-request>`, d.Client.getCurrentTimestamp(), dpskID, dpskID, fieldsToString(fields))

	if d.Client.Debug {
		fmt.Println(body)
	}

	// Create the request object
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the request headers
	req.Header.Set("X-CSRF-Token", d.Client.csrfToken)
	req.Header.Set("Cookie", d.Client.cookie)
	req.Header.Set("Content-Type", "text/xml")

	// Send the request
	resp, err := d.Client.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status code indicates success (e.g., 200 OK)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("create DPSK user failed with status code: %v", resp.StatusCode)
	}

	return nil
}

func fieldsToString(m map[string]string) string {
	pairs := make([]string, 0, len(m))
	for key, value := range m {
		pairs = append(pairs, fmt.Sprintf("%s='%s'", key, value))
	}
	return strings.Join(pairs, " ")
}
