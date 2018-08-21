package utils

import (
	"bytes"
	"encoding/json"
	"io"
)

// SafeUnmarshal preserves numbers not making everything float when converting to interface{} or map[]
func SafeUnmarshal(input []byte, output interface{}) error {
	return SafeUnmarshalReader(bytes.NewReader(input), output)
}

// SafeUnmarshalReader preserves numbers not making everything float when converting to interface{} or map[]
func SafeUnmarshalReader(r io.Reader, output interface{}) error {
	decoder := json.NewDecoder(r)
	decoder.UseNumber()
	return decoder.Decode(output)
}
