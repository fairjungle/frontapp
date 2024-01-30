package frontapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"time"
)

const (
	defaultAPIURL  = "https://api2.frontapp.com/"
	defaultTimeout = 80 * time.Second
)

// Option is a function that overwrites default client behaviour
type Option func(*Client)

// Client is a frontapp API client
type Client struct {
	apiToken string
	apiURL   string
	client   *http.Client
	log      *slog.Logger
	timeout  time.Duration
}

type errorResp struct {
	Error Error `json:"_error"`
}

// NewClient instantiates a new client
func NewClient(apiToken string, options ...Option) *Client {
	// default client
	result := &Client{
		apiToken: apiToken,
		apiURL:   defaultAPIURL,
		log:      slog.Default(),
		timeout:  defaultTimeout,
	}

	// overwrite default values
	for _, option := range options {
		option(result)
	}

	result.client = &http.Client{
		Timeout: result.timeout,
	}

	return result
}

func (c *Client) sendReq(method string, urlPath string, body interface{}, resp interface{}) error {
	// compute URL
	u, err := url.Parse(c.apiURL)
	if err != nil {
		return fmt.Errorf("failed to parse api url: %w", err)
	}
	u.Path = path.Join(u.Path, urlPath)
	url := u.String()

	// setup request body reader
	var reader io.Reader
	if body != nil {
		j, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshall request body: %w", err)
		}
		reader = bytes.NewReader(j)
	}

	c.log.Debug(fmt.Sprintf("Sending request: %s", url))

	// create request
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))

	// // DEBUG
	// dump, _ := httputil.DumpRequestOut(req, true)
	// c.log.Debug(fmt.Sprintf("Sending request: %s", string(dump))

	// send request
	response, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer response.Body.Close()

	// // DEBUG
	// raw, _ := ioutil.ReadAll(response.Body)
	// c.log.Debug(fmt.Sprintf("Received response: %s", string(raw))

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		// NOOP
	default:
		// error returned
		errResp := &errorResp{}
		if err := json.NewDecoder(response.Body).Decode(errResp); err != nil {
			return fmt.Errorf("failed to unmarshall error response (status: %d): %w", response.StatusCode, err)
		}
		return errResp.Error
	}

	if resp != nil {
		// decode response
		if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
			return fmt.Errorf("failed to unmarshall response: %w", err)
		}
	}

	return nil
}

//
// Options
//

// APIURL sets client API URL
func APIURL(val string) Option {
	return func(client *Client) {
		client.apiURL = val
	}
}

// Logger sets client logger
func Logger(val *slog.Logger) Option {
	return func(client *Client) {
		client.log = val
	}
}

// Logger sets client timeout
func Timeout(val time.Duration) Option {
	return func(client *Client) {
		client.timeout = val
	}
}
