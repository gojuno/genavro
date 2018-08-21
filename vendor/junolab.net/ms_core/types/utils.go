package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func UnmarshalBytes(data []byte) (interface{}, error) {
	type dummyUnmarshalObject struct {
		Value interface{} `json:"value"`
	}

	tmp := []byte(fmt.Sprintf(`{"value": %s}`, data))
	tmpV := dummyUnmarshalObject{}
	err := json.Unmarshal(tmp, &tmpV) // string(data)

	if err != nil {
		return "", err
	}
	return tmpV.Value, nil
}

func StreamToBytes(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func Unquote(value interface{}) (string, error) {
	var bytes []byte

	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return "", fmt.Errorf("Could not convert value '%+v' to byte array", value)
	}

	// If the value is quoted, strip the quotes
	if len(bytes) > 2 && bytes[0] == '"' && bytes[len(bytes)-1] == '"' {
		bytes = bytes[1 : len(bytes)-1]
	}
	return string(bytes), nil
}

func CheckValidRunes(validRunes, str string) error {
	for _, runeValue := range str {
		if !strings.ContainsRune(validRunes, runeValue) {
			return fmt.Errorf("invalid character found: %s", string(runeValue))
		}
	}

	return nil
}
