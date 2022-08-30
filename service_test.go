package rest

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {

	cnf := ServiceConfig{
		Interface: "localhost",
		Port:      9999,
	}

	srv := NewService(cnf)

	r := srv.Router()
	if r == nil {
		t.Errorf("expected an initialiced router for the Service")
	}

	expected := fmt.Sprintf("%s:%d", cnf.Interface, cnf.Port)
	if srv.Address != expected {
		t.Errorf("unexpected listen address: \n\t got %v\n\twant %v", srv.Address, expected)
	}

	if srv.router == nil {
		t.Errorf("uninitialized service Router")
	}

	cnf = ServiceConfig{}

	srv = NewService(cnf)

	expected = "127.0.0.1:1332"
	if srv.Address != expected {
		t.Errorf("unexpected listen address: \n\t got %v\n\twant %v", srv.Address, expected)
	}
}

func TestListenAndServe(t *testing.T) {
	cnf := ServiceConfig{
		Interface:       "localhost",
		Port:            9999,
		ShutdownTimeout: 1 * time.Second,
	}

	srv := NewService(cnf)

	go func(s *Service) {
		err := srv.ListenAndServe()
		if err != nil {
			t.Errorf("Service Listeng return unexpected error:\ngot  %s\nwant NoError", err)
		}
	}(&srv)

	err := shutdown(context.Background(), srv.srv, srv.ShutdownTimeout)
	if err != nil {
		t.Errorf("Service shutdown return unexpected error:\ngot  %s\nwant NoError", err)
	}

	// syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}
