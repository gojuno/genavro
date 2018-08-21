package types

import (
	"database/sql/driver"
	"fmt"
)

type SecurityGroup string

type securityGroupDescr struct {
	name string
	flag int
}

const (
	SecurityGroupUnknown SecurityGroup = "unknown"
)

var securityGroups map[SecurityGroup]securityGroupDescr

func RegisterSecurityGroup(sg SecurityGroup, flag int) {
	securityGroups[sg] = securityGroupDescr{
		name: sg.String(),
		flag: flag,
	}
}

func (group SecurityGroup) String() string {
	return string(group)
}

func (group SecurityGroup) Eq(other SecurityGroup) bool {
	return group.String() == other.String()
}

func (group SecurityGroup) Flag() int {
	if groupDescr, ok := securityGroups[group]; ok {
		return groupDescr.flag
	}
	return -1
}

func (group SecurityGroup) Value() (driver.Value, error) {
	value := group
	if _, ok := securityGroups[group]; !ok {
		value = SecurityGroupUnknown
	}
	return value.String(), nil
}

// Scan implements the sql.Scanner interface for database deserialization.
func (group *SecurityGroup) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
		*group = SecurityGroupUnknown
	case string:
		*group = strToSG(v)
	case []byte:
		*group = strToSG(string(v))
	default:
		return fmt.Errorf("Not supported Security Group type: %v", v)
	}

	return nil
}

func (group SecurityGroup) Validate() error {
	_, registered := securityGroups[group]
	if !registered {
		return fmt.Errorf("unregistered security group '%s'", group)
	}
	return nil
}

func strToSG(s string) (group SecurityGroup) {
	sg := SecurityGroup(s)
	if _, ok := securityGroups[sg]; ok {
		return sg
	}

	return SecurityGroupUnknown
}

func init() {
	securityGroups = map[SecurityGroup]securityGroupDescr{
		SecurityGroupUnknown: {
			name: SecurityGroupUnknown.String(),
			flag: -1,
		},
	}
}
