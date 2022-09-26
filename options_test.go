package reverseip

import (
	"net/url"
	"reflect"
	"testing"
)

// TestOptions tests the Options functions.
func TestOptions(t *testing.T) {
	tests := []struct {
		name   string
		values url.Values
		option Option
		want   string
	}{
		{
			name:   "output format",
			values: url.Values{},
			option: OptionOutputFormat("JSON"),
			want:   "outputFormat=JSON",
		},
		{
			name:   "from",
			values: url.Values{},
			option: OptionFrom("test.com"),
			want:   "from=test.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.option(tt.values)
			if got := tt.values.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Option() = %v, want %v", got, tt.want)
			}
		})
	}
}
