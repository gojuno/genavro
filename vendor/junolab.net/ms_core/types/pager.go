package types

import (
	"encoding/base64"
	"strings"
)

type Pager struct {
	Limit      int64  `json:"limit,omitempty"`
	Direction  string `json:"direction,omitempty"`
	Token      string `json:"token,omitempty"`
	Id         string `json:"id,omitempty"`
	InputValue string `json:"input_value,omitempty"`
	FirstValue string `json:"first_token,omitempty"`
	LastValue  string `json:"last_token,omitempty"`
}

type PagerDefaultSet struct {
	LimitMax     int64  `key:"limit_max"`
	LimitDefault int64  `key:"limit_default"`
	Direction    string `key:"direction"`
}

type PagerConfig struct {
	Limits map[string]PagerDefaultSet `key:"limits"`
}

const (
	_NEXT              = "next"
	_PREV              = "prev"
	_DEFAULT_LIMIT     = 10
	_DEFAULT_DIRECTION = _NEXT
)

var (
	// should be init on ms startup
	PagerDefaults  = map[string]PagerDefaultSet{}
	PagerLimitsCfg = PagerConfig{}
)

func (page *Pager) ToResponse() {
	page.Limit = 0
	page.Token = ""
	page.Id = ""
	page.InputValue = ""
}

// format of token: id$real_value
func (page *Pager) Decode(key string) {
	defaults, ok := PagerDefaults[key]
	if !ok {
		defaults = PagerDefaultSet{_DEFAULT_LIMIT, _DEFAULT_LIMIT, _DEFAULT_DIRECTION}
	}

	if page.Limit > defaults.LimitMax {
		page.Limit = defaults.LimitMax
	}
	if page.Limit == 0 {
		page.Limit = defaults.LimitDefault
	}

	if page.Direction == "" {
		page.Direction = defaults.Direction
	}

	if page.Token == "" {
		return
	}

	dst := []byte{}
	dst, err := base64.StdEncoding.DecodeString(page.Token)
	if err != nil {
		return
	}

	decoded := strings.Split(string(dst), "$")
	if len(decoded) != 2 {
		return
	}

	page.Id = decoded[0]
	page.InputValue = decoded[1]
}

func (page *Pager) EncodeToken(id, val string, first bool) {
	src := []byte((id + "$" + val))
	dst := base64.StdEncoding.EncodeToString(src)
	if first {
		page.FirstValue = dst
	} else {
		page.LastValue = dst
	}
}
