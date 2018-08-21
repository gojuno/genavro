package types

import (
	"encoding/json"
	"fmt"
	"time"

	"junolab.net/lib_api/core"
)

type AuthInfo struct {
	UserId    core.ID       `json:"user_id,omitempty"`
	Group     SecurityGroup `json:"group,omitempty"`
	ExpiredAt time.Time     `json:"expired_at,omitempty"`
}

type MsgMeta struct {
	AuthInfo    *AuthInfo         `json:"auth,omitempty"`
	HTTPHeaders map[string]string `json:"headers,omitempty"`
	Vars        map[string]string `json:"vars,omitempty"`

	Fields map[string]string `json:"fields,omitempty"`

	AppID      AppID      `json:"app_id,omitempty"`
	AppVersion AppVersion `json:"app_version,omitempty"`

	// Pager logic.
	Page        *Pager `json:"page,omitempty"`
	RequestSubj string `json:"subj,omitempty"`

	// Timing logic.
	TimingInfo *TimingInfo `json:"timing_info,omitempty"`
}

//SetIP set IP in meta
func (meta *MsgMeta) SetIP(ip string) {
	meta.Fields["ip"] = ip
}

//GetIP returns IP from meta
func (meta *MsgMeta) GetIP() string {
	return meta.Fields["ip"]
}

func (meta *MsgMeta) SetSystemUserAgent(sua string) {
	meta.Fields["system_user_agent"] = sua
}

func (meta *MsgMeta) GetSystemUserAgent() string {
	return meta.Fields["system_user_agent"]
}

func (meta *MsgMeta) SetTimingInfo(t *TimingInfo) {
	meta.TimingInfo = t
}

// GetTimingInfo returns copy of timing info!
func (meta *MsgMeta) GetTimingInfo() TimingInfo {
	if meta.TimingInfo == nil {
		return *NewTimingInfo()
	}
	return *meta.TimingInfo
}

func (meta *MsgMeta) MakeCheckPoint() {
	if meta.TimingInfo == nil {
		return
	}
	meta.TimingInfo.CheckPoint()
}

func (meta *MsgMeta) TimingInfoAsString() string {
	if meta.TimingInfo == nil {
		return ""
	}
	return meta.TimingInfo.String()
}

// IsDeadline returns flag which indicates whether we can continue execution or stop it because of timeouts.
func (meta *MsgMeta) IsDeadline() bool {
	if meta.TimingInfo == nil {
		return false
	}
	return meta.TimingInfo.IsDeadline()
}

// DurationUntilDeadline returns duration before deadline
func (meta *MsgMeta) DurationUntilDeadline() (time.Duration, bool) {
	if meta.TimingInfo == nil {
		return 0, false
	}
	return meta.TimingInfo.DurationUntilDeadline()
}

// MergeMeta merges metainfo.
func (ctx *MsgMeta) MergeMeta(meta *MsgMeta) {
	ctx.Page = meta.Page
}

func (meta *MsgMeta) SetPage(page *Pager) {
	meta.Page = page
}

func (meta *MsgMeta) GetPage() *Pager {
	if meta.Page == nil {
		return &Pager{}
	}
	return meta.Page
}

func (meta *MsgMeta) FinishPager() {
	// if pager is not used remove it from meta.
	if meta.Page != nil && meta.Page.FirstValue == "" {
		meta.Page = nil
	}
	if meta.Page != nil {
		meta.Page.ToResponse()
	}
}

func (meta *MsgMeta) SetAuth(auth *AuthInfo) error {
	meta.AuthInfo = auth
	return nil
}

func (meta *MsgMeta) SetRequestSubj(subj string) {
	meta.RequestSubj = subj
}

func (meta *MsgMeta) CleanSubj() {
	meta.RequestSubj = ""
}

func (meta *MsgMeta) Auth() *AuthInfo {
	if meta.AuthInfo == nil {
		return new(AuthInfo)
	}
	return meta.AuthInfo
}

// Intentionally defining this with non-pointer receiver - this way it works both for pointer and non-pointer cases when formatting with fmt package;
// this method is necessary for convenient meta logging - otherwise we could see in-memory pointers addresses instead of some useful info (like e.g. authenticated user id, etc).
func (meta MsgMeta) String() string {
	b, err := json.Marshal(meta)
	if err != nil {
		return fmt.Sprintf("%+v (marshal error: %v)", meta, err)
	}
	return string(b)
}

// SetHeader sets header
func (meta *MsgMeta) SetHeader(name string, value string) error {
	meta.HTTPHeaders[name] = value
	return nil
}

// Header gets header
func (meta *MsgMeta) Header(name string) (string, bool) {
	v, ok := meta.HTTPHeaders[name]
	return v, ok
}

// SetVar sets Var
func (meta *MsgMeta) SetVar(name string, value string) error {
	meta.Vars[name] = value
	return nil
}

// Var gets Var
func (meta *MsgMeta) Var(name string) (string, bool) {
	v, ok := meta.Vars[name]
	return v, ok
}

// SetField sets header
func (meta *MsgMeta) SetField(name string, value string) error {
	meta.Fields[name] = value
	return nil
}

// Field gets header
func (meta *MsgMeta) Field(name string) (string, bool) {
	v, ok := meta.Fields[name]
	return v, ok
}

// Language gets language
func (meta *MsgMeta) Language() string {
	v, ok := meta.Fields["language"]
	if !ok || v == "" {
		// todo: EN value is invalid, it must be en-US
		v = "EN" //default language
	}

	return v
}
