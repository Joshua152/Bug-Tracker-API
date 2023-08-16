package data

import (
	"encoding/json"
	"io"
)

type Bug struct {
	ProjectID   int32   `json:"projectID"`
	BugID       int32   `json:"bugID"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	TimeAmt     float64 `json:"timeAmt"`
	Complexity  float64 `json:"complexity"`
}

func (b *Bug) FromJSON(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(b)
}
