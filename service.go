package ApiService

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

var ErrGracefullShutdown = errors.New("service stopped gracefully")

type Service struct {
	Version         string
	Name            string
	Address         string
	WriteTimeout    time.Duration
	ReadTimeout     time.Duration
	ShutdownTimeout time.Duration
	srv             *http.Server
	router          *mux.Router
}

type ServiceConfig struct {
	Interface       string
	Port            int
	ShutdownTimeout time.Duration
	WriteTimeout    time.Duration
	ReadTimeout     time.Duration
}

func NewService(name string, version string, cnf ServiceConfig) Service {

	if cnf.ShutdownTimeout == 0 {
		cnf.ShutdownTimeout = time.Duration(30) * time.Second
	}

	if cnf.WriteTimeout == 0 {
		cnf.WriteTimeout = time.Duration(30) * time.Second
	}

	if cnf.ReadTimeout == 0 {
		cnf.ReadTimeout = time.Duration(30) * time.Second
	}

	if cnf.Interface == "" {
		cnf.Interface = "127.0.0.1"
	}

	if cnf.Port == 0 {
		cnf.Port = 1332
	}

	srv := Service{
		Name:            name,
		Version:         version,
		Address:         fmt.Sprintf("%v:%v", cnf.Interface, cnf.Port),
		ShutdownTimeout: cnf.ShutdownTimeout,
		WriteTimeout:    cnf.WriteTimeout,
		ReadTimeout:     cnf.ReadTimeout,
		router:          mux.NewRouter(),
	}

	srv.srv = &http.Server{
		Handler:      srv.router,
		Addr:         srv.Address,
		WriteTimeout: srv.WriteTimeout,
		ReadTimeout:  srv.ReadTimeout,
	}

	return srv
}

func (s *Service) Router() *mux.Router {
	return s.router
}

func (s *Service) PrintWelcome() {
	fmt.Println("------------------------------------------------------------")
	fmt.Println(s.Name + " v" + s.Version)
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("* Address: %v\n", s.Address)
	fmt.Println("------------------------------------------------------------")
}

func stopChannel() (chan os.Signal, func()) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	return stopCh, func() {
		close(stopCh)
	}
}

func shutdown(ctx context.Context, server *http.Server, timeout time.Duration) error {
	var err error
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err == nil {
		err = ErrGracefullShutdown
	}

	return err
}

func (s *Service) ListenAndServe() error {
	go func(srv *http.Server) {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}(s.srv)
	stopCh, closeCh := stopChannel()
	defer closeCh()
	log.Println("Shutdown Notified (Timeout 30sec):", <-stopCh)

	return shutdown(context.Background(), s.srv, s.ShutdownTimeout)
}
