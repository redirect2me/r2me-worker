package main

import (
	"context"
	"errors"
	"net"
	"strings"
)

func IsCertEligible(ctx context.Context, domain string) error {

	if domain == "" {
		return errors.New("domain cannot be empty")
	}

	if strings.HasSuffix(strings.ToLower(domain), "example.com") {
		return errors.New("example.com is not eligible for certificate issuance")
	}

	if net.ParseIP(domain) != nil && domain != Config.AdminIP {
		return errors.New("IP addresses are not eligible for certificate issuance")
	}

	return nil
}
