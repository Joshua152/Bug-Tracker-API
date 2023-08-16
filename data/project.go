package data

import (
	"encoding/json"
	"io"
)

type Project struct {
	ProjectID int32  `json:"projectID"`
	Name      string `json:"name"`
}

func (p *Project) FromJSON(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(p)
}
