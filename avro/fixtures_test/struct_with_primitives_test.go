package fixtures_test

import (
	"junolab.net/lib_api/core"
	"junolab.net/lib_api/timeapi"
)

const minorVersionPrimitivesV1 = "1"

type PrimitivesV1 struct {
	// comment here
	Int       int               `json:"int"`
	Int64     int64             `json:"int_64"`
	Float32   float32           `json:"float_32"`
	Float64   float64           `json:"float_64"`
	Bool      bool              `json:"bool"`
	String    string            `json:"string"`
	Map       map[string]string `json:"map"`
	Slice     []int             `json:"slice"`
	MapOpt    map[string]string `json:"map_opt,omitempty"`
	SliceOpt  []int             `json:"slice_opt,omitempty"`
	Omitempty int               `json:"omitempty,omitempty"`
	Ptr       *int              `json:"ptr,omitempty"`
	ID        core.ID           `json:"id"`
	Time      timeapi.Time      `json:"time"`
	Duration  timeapi.Duration  `json:"duration"`
}
