package reverseip

import (
	"fmt"
)

// Result is a part of the Reverse IP/DNS API response.
type Result struct {
	// Name is the domain name.
	Name string `json:"name"`

	// FirstSeen is the timestamp of the first time that the record was seen.
	FirstSeen int64 `json:"first_seen"`

	// LastVisit is the timestamp of the last update for this record.
	LastVisit int64 `json:"last_visit"`
}

// ReverseIPResponse is a response of Reverse IP/DNS API.
type ReverseIPResponse struct {
	// Result is a segment that contains info about the resulting data.
	Result []Result `json:"result"`

	// CurrentPage is the selected page.
	CurrentPage string `json:"current_page"`

	// Size is the number of records in the Result segment.
	Size int `json:"size"`
}

// ErrorMessage is the error message.
type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"messages"`
}

// Error returns error message as a string.
func (e *ErrorMessage) Error() string {
	return fmt.Sprintf("API error: [%d] %s", e.Code, e.Message)
}
