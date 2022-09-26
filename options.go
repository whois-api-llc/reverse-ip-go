package reverseip

import (
	"net/url"
	"strings"
)

// Option adds parameters to the query.
type Option func(v url.Values)

var _ = []Option{
	OptionOutputFormat("JSON"),
	OptionFrom("whoisxmlapi.com"),
}

// OptionOutputFormat sets Response output format JSON | XML. Default: JSON.
func OptionOutputFormat(outputFormat string) Option {
	return func(v url.Values) {
		v.Set("outputFormat", strings.ToUpper(outputFormat))
	}
}

// OptionFrom sets the domain name which is used as an offset for the results returned.
func OptionFrom(value string) Option {
	return func(v url.Values) {
		v.Set("from", value)
	}
}
