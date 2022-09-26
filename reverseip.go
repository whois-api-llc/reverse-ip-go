package reverseip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

// ReverseIP is an interface for Reverse IP/DNS API.
type ReverseIP interface {
	// Get returns parsed Reverse IP/DNS API response.
	Get(ctx context.Context, ip net.IP, opts ...Option) (*ReverseIPResponse, *Response, error)

	// GetRaw returns raw Reverse IP/DNS API response as the Response struct with Body saved as a byte slice.
	GetRaw(ctx context.Context, ip net.IP, opts ...Option) (*Response, error)
}

// Response is the http.Response wrapper with Body saved as a byte slice.
type Response struct {
	*http.Response

	// Body is the byte slice representation of http.Response Body
	Body []byte
}

// reverseIPServiceOp is the type implementing the ReverseIP interface.
type reverseIPServiceOp struct {
	client  *Client
	baseURL *url.URL
}

var _ ReverseIP = &reverseIPServiceOp{}

// newRequest creates the API request with default parameters and the specified apiKey.
func (service reverseIPServiceOp) newRequest() (*http.Request, error) {
	req, err := service.client.NewRequest(http.MethodGet, service.baseURL, nil)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("apiKey", service.client.apiKey)

	req.URL.RawQuery = query.Encode()

	return req, nil
}

// apiResponse is used for parsing Reverse IP/DNS API response as a model instance.
type apiResponse struct {
	ReverseIPResponse
	ErrorMessage
}

// request returns intermediate API response for further actions.
func (service reverseIPServiceOp) request(ctx context.Context, ip string, opts ...Option) (*Response, error) {
	var err error

	req, err := service.newRequest()
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()

	for _, opt := range opts {
		opt(q)
	}
	q.Set("ip", ip)

	req.URL.RawQuery = q.Encode()

	var b bytes.Buffer

	resp, err := service.client.Do(ctx, req, &b)
	if err != nil {
		return &Response{
			Response: resp,
			Body:     b.Bytes(),
		}, err
	}

	return &Response{
		Response: resp,
		Body:     b.Bytes(),
	}, nil
}

// parse parses raw Reverse IP/DNS API response.
func parse(raw []byte) (*apiResponse, error) {
	var response apiResponse

	err := json.NewDecoder(bytes.NewReader(raw)).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("cannot parse response: %w", err)
	}

	return &response, nil
}

// Get returns parsed Reverse IP/DNS API response.
func (service reverseIPServiceOp) Get(
	ctx context.Context,
	ip net.IP,
	opts ...Option,
) (reverseIPResponse *ReverseIPResponse, resp *Response, err error) {
	ipString := ip.String()
	if ipString == "<nil>" {
		return nil, nil, &ArgError{"ip", "can not be empty"}
	}

	optsJSON := make([]Option, 0, len(opts)+1)
	optsJSON = append(optsJSON, opts...)
	optsJSON = append(optsJSON, OptionOutputFormat("JSON"))

	resp, err = service.request(ctx, ipString, optsJSON...)
	if err != nil {
		return nil, resp, err
	}

	reverseIPResp, err := parse(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	if reverseIPResp.Message != "" || reverseIPResp.Code != 0 {
		return nil, nil, &ErrorMessage{
			Code:    reverseIPResp.Code,
			Message: reverseIPResp.Message,
		}
	}

	return &reverseIPResp.ReverseIPResponse, resp, nil
}

// GetRaw returns raw Reverse IP/DNS API response as the Response struct with Body saved as a byte slice.
func (service reverseIPServiceOp) GetRaw(
	ctx context.Context,
	ip net.IP,
	opts ...Option,
) (resp *Response, err error) {
	ipString := ip.String()
	if ipString == "<nil>" {
		return nil, &ArgError{"ip", "can not be empty"}
	}

	resp, err = service.request(ctx, ipString, opts...)
	if err != nil {
		return resp, err
	}

	if respErr := checkResponse(resp.Response); respErr != nil {
		return resp, respErr
	}

	return resp, nil
}

// ArgError is the argument error.
type ArgError struct {
	Name    string
	Message string
}

// Error returns error message as a string.
func (a *ArgError) Error() string {
	return `invalid argument: "` + a.Name + `" ` + a.Message
}
