package ApiService

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out
}

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

	r := srv.Router()
	if r == nil {
		t.Errorf("expected an initialiced router for the Service")
	}

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

func TestListenAndServe(t *testing.T) {
	cnf := ServiceConfig{
		Interface:       "localhost",
		Port:            9999,
		ShutdownTimeout: 1 * time.Second,
	}

	srv := NewService("Test Service", "1.0.0", cnf)

	go func(s *Service) {
		err := srv.ListenAndServe()
		if err != nil {
			t.Errorf("Service Listeng return unexpected error:\ngot  %s\nwant NoError", err)
		}
	}(&srv)

	err := shutdown(context.Background(), srv.srv, srv.ShutdownTimeout)
	if err != ErrGracefullShutdown {
		t.Errorf("Service shutdown return unexpected error:\ngot  %s\nwant %s", err, ErrGracefullShutdown)
	}

	// syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}
