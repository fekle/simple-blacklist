package domainfilter

import (
	"regexp"
	"strings"

	"golang.org/x/net/publicsuffix"
)

var (
	commentExp    = regexp.MustCompile("#.*$")
	ipExp         = regexp.MustCompile("(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}")
	noSpaceExp    = regexp.MustCompile("\\S+")
	leadingDotExp = regexp.MustCompile("^\\.")
)

// TrimWhiteSpaceFilter trims leading and trailing whitespace
func TrimWhiteSpaceFilter(in string) []string {
	return []string{strings.Trim(in, " ")}
}

// DropCommentsFilter removes all comments and drops lines that are entirely comments
func DropCommentsFilter(in string) []string {
	if strings.HasPrefix(in, "#") {
		return []string{}
	}
	return []string{commentExp.ReplaceAllLiteralString(in, "")}
}

// DropIPAddressesFilter drops all IP addresses
func DropIPAddressesFilter(in string) []string {
	return []string{ipExp.ReplaceAllLiteralString(in, "")}
}

// ExtractNoSpaceGroupsFilter extracts all char sequences that don't contain a space
func ExtractNoSpaceGroupsFilter(in string) []string {
	return noSpaceExp.FindAllString(in, -1)
}

// HasValidTLDFilter drops all invalid domains
func HasValidTLDFilter(in string) []string {
	if strings.Contains(in, ".") {
		// check if domain has a valid tld
		suffix, _ := publicsuffix.PublicSuffix(in)

		if len(suffix) > 0 {
			return []string{leadingDotExp.ReplaceAllLiteralString(in, "*.")}
		}
	}
	return []string{}
}

// LeadingDotToWildcard replaces all trailing whitespaces with a wildcard
func LeadingDotToWildcard(in string) []string {
	return []string{leadingDotExp.ReplaceAllLiteralString(in, "*.")}
}
