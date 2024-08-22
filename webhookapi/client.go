package webhookapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// Client is the client used to subscribe to the Logs API
type Client struct {
	httpClient     *http.Client
	logsApiBaseUrl string
}

// NewClient returns a new Client with the given URL
func NewClient(logsApiBaseUrl string) (*Client, error) {
	return &Client{
		httpClient:     &http.Client{},
		logsApiBaseUrl: logsApiBaseUrl,
	}, nil
}

// URI is used to set the endpoint where the logs will be sent to
type URI string

// HttpMethod represents the HTTP method used to receive logs from Logs API
type HttpMethod string

const (
	// HttpPost is to receive logs through POST.
	HttpPost HttpMethod = "POST"
	// HttpPut is to receive logs through PUT.
	HttpPut HttpMethod = "PUT"
)

// HttpProtocol is used to specify the protocol when subscribing to Logs API for HTTP
type HttpProtocol string

const (
	HttpProto HttpProtocol = "HTTP"
)

// HttpEncoding denotes what the content is encoded in
type HttpEncoding string

const (
	JSON HttpEncoding = "JSON"
)

// Destination is the configuration for listeners who would like to receive logs with HTTP
type Destination struct {
	Protocol   HttpProtocol `json:"protocol"`
	URI        URI          `json:"URI"`
	HttpMethod HttpMethod   `json:"method"`
	Encoding   HttpEncoding `json:"encoding"`
}

// SendRequest is the request body that is sent to Logs API on subscribe
type SendRequest struct {
	Destination Destination `json:"destination"`
}

// SendResponse is the response body that is received from Logs API on subscribe
type SendResponse struct {
	Body string `json:"body"`
}

// Send sends the configuration to the Logs API and returns the response
func (c *Client) Send(destination Destination, extensionId string) (*SendResponse, error) {
	data, err := json.Marshal(&SendRequest{
		Destination: destination,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to marshal SendRequest")
	}

	headers := map[string]string{"X-Extension-Id": extensionId}
	url := c.logsApiBaseUrl

	resp, err := httpPostWithHeaders(c.httpClient, url, data, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted {
		return nil, errors.New("Logs API is not supported! Is this extension running in a local sandbox?")
	} else if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Errorf("%s failed: %d[%s]", url, resp.StatusCode, resp.Status)
		}
		return nil, errors.Errorf("%s failed: %d[%s] %s", url, resp.StatusCode, resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read response body")
	}

	return &SendResponse{Body: string(body)}, nil
}

func httpPostWithHeaders(client *http.Client, url string, data []byte, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
