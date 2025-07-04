package utils

import (
	"fmt"
	"regexp"

	"github.com/weppos/publicsuffix-go/publicsuffix"
)

func ParseDomain(domain string) (*publicsuffix.DomainName, error) {
	regDomain := regexp.MustCompile(`^([\-a-zA-Z0-9]+\.){1,}[a-zA-Z0-9]+$`)

	if !regDomain.MatchString(domain) {
		return nil, fmt.Errorf("invalid domain: incorrect format")
	}

	parsedDomain, err := publicsuffix.Parse(domain)
	if err != nil {
		return nil, fmt.Errorf("invalid domain: %v", err)
	}

	return parsedDomain, nil
}
