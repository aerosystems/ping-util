package mx

import "net"

type DomainModel struct {
	Name string
	Type DomainType
	Ip   net.IP
}

type DomainType struct {
	slug string
}

func (d DomainType) String() string {
	return d.slug
}

var (
	UnknownDomainType   = DomainType{"unknown"}
	WhitelistDomainType = DomainType{"whitelist"}
	BlacklistDomainType = DomainType{"blacklist"}
)

func DomainTypeFromString(s string) DomainType {
	switch s {
	case WhitelistDomainType.String():
		return WhitelistDomainType
	case BlacklistDomainType.String():
		return BlacklistDomainType
	default:
		return UnknownDomainType
	}
}
