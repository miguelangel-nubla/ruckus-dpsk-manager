package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Client struct {
	Debug     bool
	client    *http.Client
	server    string // Add a field to store the server address
	csrfToken string // Add a field to store the CSRF token
	cookie    string // Add a field to store the cookie
}

func New(server string, caCertPath string) (*Client, error) {
	// Base transport
	tr := &http.Transport{}

	// If caCertPath is specified, add the custom CA cert
	if caCertPath != "" {
		caCert, err := os.ReadFile(caCertPath)
		if err != nil {
			return nil, fmt.Errorf("error reading custom CA certificate: %v", err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to append custom CA certificate")
		}
		tr.TLSClientConfig = &tls.Config{
			RootCAs: caCertPool,
		}
	}

	httpClient := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &Client{client: httpClient, server: server, Debug: false}, nil
}

func (rc *Client) Login(username, password string) error {
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
	if loginResp.Header.Get("HTTP_X_CSRF_TOKEN") == "" {
		if rc.Debug {
			body, err := io.ReadAll(loginResp.Body)
			if err != nil {
				return fmt.Errorf("error reading response body: %v", err)
			}
			fmt.Println(string(body))
		}
		return fmt.Errorf("check user and password, status code: %v", loginResp.StatusCode)
	}

	rc.cookie = loginResp.Header.Get("Set-Cookie")
	rc.csrfToken = loginResp.Header.Get("HTTP_X_CSRF_TOKEN")

	return nil
}

func (rc *Client) getCurrentTimestamp() string {
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

func (rc *Client) Backup(outputFile string) error {
	// Define the URL for saving the backup
	saveBackupURL := rc.server + "/admin/webPage/system/admin/_savebackup.jsp"

	// Create an HTTP GET request to the save backup URL
	req, err := http.NewRequest("GET", saveBackupURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the necessary headers
	req.Header.Set("Accept", "application/octet-stream") // Specify the desired content type
	req.Header.Set("Cookie", rc.cookie)                  // Use the stored cookie

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
