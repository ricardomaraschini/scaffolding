package httpcli

import (
	"encoding/json"
	"net/http"
)

// Client wraps a regular http client and provides json unmarshaler out of the box. It also
// provides prometheus instrumentation and a customized http.Transport with timeouts.
type Client struct {
	*http.Client
}

// UnmarshalDo executes an http request and unmarshals the output into provided object.
func (c *Client) UnmarshalDo(req *http.Request, obj interface{}) error {
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(obj)
}

// New returns a new HTTP Client.
func New() *Client {
	return &Client{
		&http.Client{
			Transport: defaultTransport(),
		},
	}
}
