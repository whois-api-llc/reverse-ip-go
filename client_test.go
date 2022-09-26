package reverseip

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

const (
	pathReverseIPResponseOK         = "/ReverseIP/ok"
	pathReverseIPResponseError      = "/ReverseIP/error"
	pathReverseIPResponse500        = "/ReverseIP/500"
	pathReverseIPResponsePartial1   = "/ReverseIP/partial"
	pathReverseIPResponsePartial2   = "/ReverseIP/partial2"
	pathReverseIPResponseUnparsable = "/ReverseIP/unparsable"
)

const apiKey = "at_LoremIpsumDolorSitAmetConsect"

// dummyServer is the sample of the Reverse IP/DNS API server for testing.
func dummyServer(resp, respUnparsable string, respErr string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var response string

		response = resp

		switch req.URL.Path {
		case pathReverseIPResponseOK:
		case pathReverseIPResponseError:
			w.WriteHeader(499)
			response = respErr
		case pathReverseIPResponse500:
			w.WriteHeader(500)
			response = respUnparsable
		case pathReverseIPResponsePartial1:
			response = response[:len(response)-10]
		case pathReverseIPResponsePartial2:
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			response = response[:len(response)-10]
		case pathReverseIPResponseUnparsable:
			response = respUnparsable
		default:
			panic(req.URL.Path)
		}
		_, err := w.Write([]byte(response))
		if err != nil {
			panic(err)
		}
	}))

	return server
}

// newAPI returns new Reverse IP/DNS API client for testing.
func newAPI(apiServer *httptest.Server, link string) *Client {
	apiURL, err := url.Parse(apiServer.URL)
	if err != nil {
		panic(err)
	}

	apiURL.Path = link

	params := ClientParams{
		HTTPClient:       apiServer.Client(),
		ReverseIPBaseURL: apiURL,
	}

	return NewClient(apiKey, params)
}

// TestReverseIPGet tests the Get function.
func TestReverseIPGet(t *testing.T) {
	checkResultRec := func(res *ReverseIPResponse) bool {
		return res != nil
	}

	ctx := context.Background()

	const resp = `{"current_page":"0","size":3,"result":[{"name":"iana.com","first_seen":1570492800,
"last_visit":1657756800},{"name":"iana.net","first_seen":1571097600,"last_visit":1660780800},
{"name":"iana.org","first_seen":1570147200,"last_visit":1657843200}]}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory net.IP
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathReverseIPResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionFrom("1"),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathReverseIPResponse500,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionFrom("1"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathReverseIPResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionFrom("1"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathReverseIPResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionFrom("1"),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathReverseIPResponseError,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionFrom("1"),
				},
			},
			want:    false,
			wantErr: "API error: [499] Test error message.",
		},
		{
			name: "unparsable response",
			path: pathReverseIPResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionFrom("1"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "invalid argument1",
			path: pathReverseIPResponseError,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8},
					OptionFrom("1"),
				},
			},
			want:    false,
			wantErr: `API error: [499] Test error message.`,
		},
		{
			name: "invalid argument2",
			path: pathReverseIPResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{},
					OptionFrom("1"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "ip" can not be empty`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			gotRec, _, err := api.Get(tt.args.ctx, tt.args.options.mandatory, tt.args.options.option)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("ReverseIP.Get() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("ReverseIP.Get() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != nil {
					t.Errorf("ReverseIP.Get() got = %v, expected nil", gotRec)
				}
			}
		})
	}
}

// TestReverseIPGetRaw tests the GetRaw function.
func TestReverseIPGetRaw(t *testing.T) {
	checkResultRaw := func(res []byte) bool {
		return len(res) != 0
	}

	ctx := context.Background()

	const resp = `{"current_page":"0","size":3,"result":[{"name":"iana.com","first_seen":1570492800,
"last_visit":1657756800},{"name":"iana.net","first_seen":1571097600,"last_visit":1660780800},
{"name":"iana.org","first_seen":1570147200,"last_visit":1657843200}]}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory net.IP
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		wantErr string
	}{
		{
			name: "successful request",
			path: pathReverseIPResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionFrom("1"),
				},
			},
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathReverseIPResponse500,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionFrom("1"),
				},
			},
			wantErr: "API failed with status code: 500",
		},
		{
			name: "partial response 1",
			path: pathReverseIPResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionFrom("1"),
				},
			},
			wantErr: "",
		},
		{
			name: "partial response 2",
			path: pathReverseIPResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionFrom("1"),
				},
			},
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "unparsable response",
			path: pathReverseIPResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionFrom("1"),
				},
			},
			wantErr: "",
		},
		{
			name: "could not process request",
			path: pathReverseIPResponseError,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionFrom("1"),
				},
			},
			wantErr: "API failed with status code: 499",
		},
		{
			name: "invalid argument1",
			path: pathReverseIPResponseError,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8},
					OptionFrom("1"),
				},
			},
			wantErr: "API failed with status code: 499",
		},
		{
			name: "invalid argument2",
			path: pathReverseIPResponseError,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{},
					OptionFrom("1"),
				},
			},
			wantErr: `invalid argument: "ip" can not be empty`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			resp, err := api.GetRaw(tt.args.ctx, tt.args.options.mandatory)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("ReverseIP.GetRaw() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if resp != nil && !checkResultRaw(resp.Body) {
				t.Errorf("ReverseIP.GetRaw() got = %v, expected something else", string(resp.Body))
			}
		})
	}
}
