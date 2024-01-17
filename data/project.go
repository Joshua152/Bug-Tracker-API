package data

import (
	"encoding/json"
	"io"
	"time"
)

type Project struct {
	ProjectID   int32     `json:"projectID"`
	Name        string    `json:"name"`
	CreatedOn   time.Time `json:"createdOn"`
	LastUpdated time.Time `json:"lastUpdated"`
}

func (p *Project) FromJSON(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(p)
}
