package data

import (
	"encoding/json"
	"io"
)

func ToJSON(i interface{}, w io.Writer) error {
	b, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return nil
}
