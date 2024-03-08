package config

import (
	"log"
	"net"
	"strings"
	"fmt"
)

// checks if the given string is a valid IP address.
func IsValidIPAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}

// checks if the given string is a valid domain name.
func IsValidDomain(domain string) bool {
	parts := strings.Split(domain, ".")
	for _, part := range parts {
		if len(part) == 0 {
			return false
		}
	}
	return true
}

// isAddressValid checks if the address is a valid IP address or a domain.
// If it's a domain, it resolves it to its public IP address and compares with the provided value.
func IsAddressValid(address, value string) bool {
	if IsValidIPAddress(address) {
		log.Printf("Valid IP address: %s", address)
		return address == value
	}
	if IsValidDomain(address) {
		log.Printf("Valid domain: %s", address)
		ipAddr, err := GetPublicIP(address)
		if err != nil {
			log.Printf("Error resolving domain: %s", err)
			return false
		}
		log.Printf("Resolved IP address for domain %s: %s", address, ipAddr)
		return ipAddr == value
	}
	log.Printf("Invalid address format: %s", address)
	return false
}

func GetPublicIP(domain string) (string, error) {
	addrs, err := net.LookupIP(domain)
	if err != nil {
		return "", err
	}

	// Find the IPv4 address
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4.String(), nil
		}
	}

	return "", fmt.Errorf("no IPv4 address found for domain %s", domain)
}