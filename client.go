package main

import (
    "crypto/tls"
    "fmt"
    "net/http"
    "net/url"
    "strings"
	"io/ioutil"
	"os"
	"io"
	"time"
)

type RuckusClient struct {
	Debug     bool
    client    *http.Client
	server		  string // Add a field to store the server address
    csrfToken string // Add a field to store the CSRF token
    cookie    string // Add a field to store the cookie
}

func NewRuckusClient(server string) *RuckusClient {
	// Create a custom transport that skips certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create an HTTP client with the custom transport
	httpClient := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &RuckusClient{client: httpClient, server: server, Debug: false}
}

func (rc *RuckusClient) Login(username, password string) error {
    // Login URL
    loginURL := rc.server + "/admin/login.jsp"
    loginData := url.Values{
        "username": {username},
        "password": {password},
        "ok":       {"Log In"},
    }

    // Send the login request
    loginResp, err := rc.client.PostForm(loginURL, loginData)
    if err != nil {
        return fmt.Errorf("error sending login request: %v", err)
    }
    defer loginResp.Body.Close()

    // Check if login was successful
    if loginResp.StatusCode != http.StatusOK && loginResp.StatusCode != http.StatusFound {
        return fmt.Errorf("login failed with status code: %v", loginResp.StatusCode)
    }

    rc.cookie = loginResp.Header.Get("Set-Cookie")
    rc.csrfToken = loginResp.Header.Get("HTTP_X_CSRF_TOKEN")

    return nil
}

func (rc *RuckusClient) getCurrentTimestamp() string {
	// Get the current time in UnixNano format (nanoseconds since epoch)
	currentTime := time.Now().UnixNano()

	// Convert nanoseconds to milliseconds (1 nanosecond = 1e-6 milliseconds)
	milliseconds := currentTime / int64(time.Millisecond)

	// Convert nanoseconds to microseconds (1 nanosecond = 1e-3 microseconds)
	microseconds := currentTime / int64(time.Microsecond)

	// Calculate the fractional part (microseconds) by subtracting milliseconds
	fractionalPart := microseconds - (milliseconds * 1000)

	// Format the timestamp in the desired format
	// Use %d for milliseconds and %04d for microseconds with leading zeros
	timestamp := fmt.Sprintf("%d.%04d", milliseconds, fractionalPart)

	return timestamp
}

func (rc *RuckusClient) SaveBackup(outputFile string) error {
	// Define the URL for saving the backup
	saveBackupURL := rc.server + "/admin/webPage/system/admin/_savebackup.jsp"

	// Create an HTTP GET request to the save backup URL
	req, err := http.NewRequest("GET", saveBackupURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the necessary headers
	req.Header.Set("Accept", "application/octet-stream") // Specify the desired content type
	req.Header.Set("Cookie", rc.cookie) // Use the stored cookie

	// Send the GET request
	resp, err := rc.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status code indicates success (e.g., 200 OK)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("save backup failed with status code: %v", resp.StatusCode)
	}

	// Create or open the outputFile for writing
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outFile.Close()

	// Copy the response body (binary file) to the outputFile
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("error copying response body to output file: %v", err)
	}

	return nil
}

func (rc *RuckusClient) GetDpskData() (*DpskData, error) {
	// Create the request URL
	url := rc.server + "/admin/_cmdstat.jsp"

	// Create the request body
	body := `<ajax-request action="getstat" comp="stamgr" updater="dpsk-list.`+rc.getCurrentTimestamp()+`">
		<dpsklist/>
	</ajax-request>`

	if rc.Debug {
		fmt.Println(body)
	}

	// Create the request object
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set the request headers
	req.Header.Set("X-CSRF-Token", rc.csrfToken)
	req.Header.Set("Cookie", rc.cookie)
	req.Header.Set("Content-Type", "text/xml")

	// Send the request
	resp, err := rc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	xmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	dpskData, err := NewDpskData(xmlData)
	if rc.Debug {
		fmt.Printf("Parsed DPSKs:\n")
		for _, dpsk := range dpskData.Entries {
			fmt.Printf("%v\n", dpsk)
		}
	}

	return dpskData, err
}

func (rc *RuckusClient) CreateDpskUser(wlanID int, username string) error {
	// Create the request URL
	url := rc.server + "/admin/_cmdstat.jsp"

	// Create the request body
	body := fmt.Sprintf(`<ajax-request action='docmd' checkAbility='2' updater='system.`+rc.getCurrentTimestamp()+`' comp='system'>
		<xcmd
			cmd='batch-dpsk'
			type='gen'
			num='1'
			max-num='2048'
			batch-dpsk=''
			wlansvc-id='%d'
			role-id=''
			dpsk-len='8'
			dvlan-id=''
			user='%s'
		/>
	</ajax-request>`, wlanID, username)

	if rc.Debug {
		fmt.Println(body)
	}

	// Create the request object
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the request headers
	req.Header.Set("X-CSRF-Token", rc.csrfToken)
	req.Header.Set("Cookie", rc.cookie)
	req.Header.Set("Content-Type", "text/xml")

	// Send the request
	resp, err := rc.client.Do(req)
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