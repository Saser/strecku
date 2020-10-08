package names

import (
	"regexp"
	"strings"
)

const (
	hex  = "[0-9a-fA-F]"
	UUID = hex + "{8}" + "-?" + hex + "{4}" + "-?" + hex + "{4}" + "-?" + hex + "{4}" + "-?" + hex + "{12}"
)

func Compile(parts ...string) (*regexp.Regexp, error) {
	return regexp.Compile(strings.Join(parts, "/"))
}

func MustCompile(parts ...string) *regexp.Regexp {
	re, err := Compile(parts...)
	if err != nil {
		panic(err)
	}
	return re
}
