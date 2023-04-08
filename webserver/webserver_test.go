package webserver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello World")
}
func TestServerRun(t *testing.T) {
	s, _ := NewSetup()
	s.Add("/", hello)
	server, _ := s.NewServer()
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return server.Run(ctx)
	})
	rsp1, err := http.Get("http://localhost:8080/")
	if err != nil {
		t.Errorf("%+v", err)
	}

	got1, err := io.ReadAll((rsp1.Body))
	if err != nil {
		t.Errorf("%+v", err)
	}
	if string(got1) != "hello World" {
		t.Errorf("message : %v", string(got1))
	}
	defer rsp1.Body.Close()
	cancel()

}

func TestServerSetupRun(t *testing.T) {
	s, _ := NewSetup()
	c := []WebConfig{
		{Pass: "/", Handler: hello},
	}
	if err := Config(s, c); err != nil {
		t.Errorf("%+v", err)

	}
	server, _ := s.NewServer()
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return server.Run(ctx)
	})
	rsp1, err := http.Get("http://localhost:8080/")
	if err != nil {
		t.Errorf("%+v", err)
	}

	got1, err := io.ReadAll((rsp1.Body))
	if err != nil {
		t.Errorf("%+v", err)
	}
	if string(got1) != "hello World" {
		t.Errorf("message : %v", string(got1))
	}
	defer rsp1.Body.Close()
	cancel()

}
