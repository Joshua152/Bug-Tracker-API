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

type Bugs []Bug

// read in bytes instead? and unmarshall
// pass in decoder where first token is already read? --> can't because can't reverse reading after dec.Token()
func (b *Bug) FromJSON(r io.Reader) error {
	// https://pkg.go.dev/encoding/json#Decoder.Decode
	dec := json.NewDecoder(r)
	return dec.Decode(b)
}

func (b *Bugs) ToInterface() {

}
