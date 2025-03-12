package test

import (
	"fmt"
	"testing"

	"github.com/Filippo831/reverse_proxy/internal/read_configuration"
)

// missing key parameters but ssl activate
func TestMissingKeyParameters(t *testing.T) {
	err := readconfiguration.ReadConfiguration("configurations/missing_key_parameters.json")

	if fmt.Sprint(err) != "missing ssl parameter(s) in server 1" {
		t.Errorf("expected: missing ssl parameter(s) in server 1\n found: %s", err)
	}
}

// present key parameters but ssl not active
func TestSslDisabledWithDefinedKeys(t *testing.T) {
	err := readconfiguration.ReadConfiguration("configurations/https_disabled_but_keys_defined.json")

	if fmt.Sprint(err) != "ssl parameters set even if ssl is selected to false in server 0" {
		t.Errorf("expected: ssl parameters set even if ssl is selected to false in server 0\n found: %s", err)
	}
}

// test to verify that every subdomain shares the same domain name as the server name 
func TestDifferentSubdomains(t *testing.T) {
    err := readconfiguration.ReadConfiguration("configurations/different_subdomains.json")

	if fmt.Sprint(err) != "server number 0 domain: localhost\nlocation domain: gianni\n" {
		t.Errorf("expected: server number 0 domain: localhost\nlocation domain: gianni\n found: %s", err)
	}
}
