package gcurl

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		name     string
		given    string
		expected *Request
	}{
		{
			"simple get",
			"curl https://api.sloths.com",
			&Request{
				Method: http.MethodGet,
				URL:    "https://api.sloths.com",
				Header: map[string]string{},
			},
		},
		{
			"simple put",
			"curl -XPUT https://api.sloths.com/sloth/4",
			&Request{
				Method: http.MethodPut,
				URL:    "https://api.sloths.com/sloth/4",
				Header: map[string]string{},
			},
		},
		{
			"encoding gzip",
			`curl -H "Accept-Encoding: gzip" --compressed http://api.sloths.com`,
			&Request{
				Method: http.MethodGet,
				URL:    "http://api.sloths.com",
				Header: map[string]string{
					"accept-encoding": "gzip",
				},
			},
		},
		{
			"delete sloth",
			"curl -X DELETE https://api.sloths.com/sloth/4",
			&Request{
				Method: http.MethodDelete,
				URL:    "https://api.sloths.com/sloth/4",
				Header: map[string]string{},
			},
		},
		{
			"url encoded data",
			`curl -d "foo=bar" https://api.sloths.com/sloth/4`,
			&Request{
				Method: http.MethodPost,
				URL:    "https://api.sloths.com/sloth/4",
				Header: map[string]string{"content-type": "application/x-www-form-urlencoded"},
				Body:   "foo=bar",
			},
		},
		{
			"JSON",
			`curl -d '{"hello": "world"}' -H 'content-type: application/json' https://api.sloths.com/sloth/4`,
			&Request{
				Method: http.MethodPost,
				URL:    "https://api.sloths.com/sloth/4",
				Header: map[string]string{"content-type": "application/json"},
				Body:   `{"hello":"world"}`,
			},
		},
		{
			"user agent",
			`curl -H "Accept: text/plain" --header "User-Agent: slothy" https://api.sloths.com`,
			&Request{
				Method: http.MethodGet,
				URL:    "https://api.sloths.com",
				Header: map[string]string{
					"accept":     "text/plain",
					"user-agent": "slothy",
				},
			},
		},
		{
			"cookie",
			`curl --cookie 'species=sloth;type=galactic' slothy https://api.sloths.com`,
			&Request{
				Method: http.MethodGet,
				URL:    "https://api.sloths.com",
				Header: map[string]string{
					"cookie": "species=sloth;type=galactic",
				},
			},
		},
		{
			"location",
			`curl --location --request GET 'https://api.sloths.com/users?token=admin'`,
			&Request{
				Method: http.MethodGet,
				URL:    "https://api.sloths.com/users?token=admin",
				Header: map[string]string{},
			},
		},
		{
			"timeout and skip TLS",
			`curl --max-time 30 -k 'https://api.sloths.com/users?token=admin'`,
			&Request{
				Method:  http.MethodGet,
				URL:     "https://api.sloths.com/users?token=admin",
				Header:  map[string]string{},
				Timeout: "30",
				SkipTLS: true,
			},
		},
		{
			"repeated data fields",
			`curl -d 'foo=bar&bar=foo' -d 'q=GoogleQuery' https://api.sloths.com/sloth/4`,
			&Request{
				Method: http.MethodPost,
				URL:    "https://api.sloths.com/sloth/4",
				Header: map[string]string{"content-type": "application/x-www-form-urlencoded"},
				Body:   "foo=bar&bar=foo&q=GoogleQuery",
			},
		},
		{
			"custom authorization",
			`curl -H 'Authorization: Token some-custom-auth' https://api.sloths.com/sloth/4`,
			&Request{
				Method: http.MethodGet,
				URL:    "https://api.sloths.com/sloth/4",
				Header: map[string]string{"authorization": "Token some-custom-auth"},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual, err := Parse(tt.given)
			require.NoError(t, err)
			require.Equal(t, tt.expected, actual)
		})
	}
}
