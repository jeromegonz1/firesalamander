package config

import "fmt"

const (
	VersionMajor = 1
	VersionMinor = 0
	VersionPatch = 0
	VersionPre   = "" // "alpha", "beta", "rc1"
)

func Version() string {
	v := fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	if VersionPre != "" {
		v += "-" + VersionPre
	}
	return v
}

func VersionShort() string {
	return fmt.Sprintf("v%s", Version())
}

func VersionFull() string {
	return fmt.Sprintf("Fire Salamander %s", VersionShort())
}