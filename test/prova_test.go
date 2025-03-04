package test 

import (
	"testing"
    "github.com/Filippo831/reverse_proxy/internal/prove"
)

func TestFirst(t *testing.T) {
	if 40 != prove.PrintHello() {
		t.Errorf("not good")
	}
}
