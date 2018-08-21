package utils

import (
	"encoding/json"

	"junolab.net/ms_core/types"
)

// MergeJSONData - take list of JSONs, merge it into one object
// example:
//    out := struct {
//    		F1 int `json:"f1"`
//    		F2 int `json:"f2"`
//    }{}
//    err := utils.MergeJSONData(&out, []byte(`{"f1": 345, "f2": 567}`), []byte(`{"f1": 123}`})
//    out is: {F1: 123, F2: 567}
func MergeJSONData(outPointer interface{}, dataList ...types.RawMessage) error {
	for _, data := range dataList {
		if len(data) == 0 {
			continue
		}

		if err := json.Unmarshal(data, outPointer); err != nil {
			return err
		}
	}

	return nil
}
