package data

import (
	"encoding/json"
	"io"
	"time"
)

type Bug struct {
	ProjectID   int32     `json:"projectID"`
	BugID       int32     `json:"bugID"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TimeAmt     float64   `json:"timeAmt"`
	Complexity  float64   `json:"complexity"`
	CreatedOn   time.Time `json:"createdOn"`
	LastUpdated time.Time `json:"lastUpdated"`
}

type Bugs []Bug

func (b *Bug) FromJSON(r io.Reader) error {
	// https://pkg.go.dev/encoding/json#Decoder.Decode
	dec := json.NewDecoder(r)
	return dec.Decode(b)
}
