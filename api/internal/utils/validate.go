package utils

import (
	"regexp"
	"strings"
)

var (
	reEmail = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	reUsername = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]{3,32}$`)
	reSubLabel = regexp.MustCompile(`^[a-z0-9]([a-z0-9\-]{0,61}[a-z0-9])?$`)
)

func IsValidEmail(s string) bool {
	return reEmail.MatchString(s)
}

func IsValidUsername(s string) bool {
	return reUsername.MatchString(s)
}

func NormalizeSubdomain(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.Trim(s, ".")
	return s
}

func IsValidSubdomainLabel(s string) bool {
	return reSubLabel.MatchString(s)
}
