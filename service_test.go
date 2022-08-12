package ApiService

import (
	"fmt"
	"testing"
)

func TestServiceWelcommeMessage(t *testing.T) {

	output := captureOutput(func() {
		srv := NewService("Test Service", "1.0.0", ServiceConfig{
			Interface: "localhost",
			Port:      1234,
		})
		srv.PrintWelcome()
	})

	expected := `------------------------------------------------------------
Test Service v1.0.0
------------------------------------------------------------
* Address: localhost:1234
------------------------------------------------------------
`
	if output != expected {
		t.Errorf("function returned unexpected message: \n\t got -->\n%v\n\twant -->\n%v", output, expected)
	}
}

func TestNewService(t *testing.T) {

	cnf := ServiceConfig{
		Interface: "localhost",
		Port:      9999,
	}

	srv := NewService("Test Service", "1.0.0", cnf)

	expected := fmt.Sprintf("%s:%d", cnf.Interface, cnf.Port)
	if srv.Address != expected {
		t.Errorf("unexpected listen address: \n\t got %v\n\twant %v", srv.Address, expected)
	}

	expected = "Test Service"
	if srv.Name != expected {
		t.Errorf("unexpected service Name: \n\t got %v\n\twant %v", srv.Name, expected)
	}

	expected = "1.0.0"
	if srv.Version != expected {
		t.Errorf("unexpected service Version: \n\t got %v\n\twant %v", srv.Name, expected)
	}

	if srv.router == nil {
		t.Errorf("uninitialized service Router")
	}

	cnf = ServiceConfig{}

	srv = NewService("Test Service", "1.0.0", cnf)

	expected = "127.0.0.1:1332"
	if srv.Address != expected {
		t.Errorf("unexpected listen address: \n\t got %v\n\twant %v", srv.Address, expected)
	}
}
