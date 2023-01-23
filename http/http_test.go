package http

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/sters/pinger/pinger"
)

func TestWithClient(t *testing.T) {
	t.Parallel()

	h := &pingerOptions{}
	if h.client != nil {
		t.Errorf("already configured http client")
	}

	WithClient(http.DefaultClient)(h)
	if h.client == nil {
		t.Errorf("not configured http client")
	}
}

func TestHTTPPinger(t *testing.T) {
	t.Parallel()

	testWantMethod := "GET"
	handledEndpoint := false
	mux := sync.Mutex{}

	handler := http.NewServeMux()
	handler.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != testWantMethod {
			t.Fatalf("request method want = %s, got = %s", testWantMethod, r.Method)
		}
		if got := r.Header.Get("X-ADDITIONAL-HEADER-1"); got != "foo" {
			t.Fatalf("request header want = foo, got = %s", got)
		}
		if got := r.Header.Get("X-ADDITIONAL-HEADER-2"); got != "bar" {
			t.Fatalf("request header want = bar, got = %s", got)
		}

		mux.Lock()
		defer mux.Unlock()
		handledEndpoint = true
	})
	s := &http.Server{
		Addr:              "localhost:8080",
		Handler:           handler,
		ReadHeaderTimeout: time.Second,
	}
	go func() { _ = s.ListenAndServe() }()
	defer s.Close()

	u, _ := url.Parse("http://localhost:8080/bar")
	p := NewPinger(
		u,
		WithMethod(testWantMethod),
		WithHeader("X-ADDITIONAL-HEADER-1", "foo"),
		WithHeader("X-ADDITIONAL-HEADER-2", "bar"),
	)
	worker := pinger.NewWorker(p, pinger.WithInterval(time.Millisecond))

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Millisecond))
	defer cancel()
	if err := worker.Run(ctx); err != nil {
		t.Fatalf("worker.Run want no error, got error: %+v", err)
	}

	mux.Lock()
	defer mux.Unlock()
	if handledEndpoint == false {
		t.Fatalf("worker.Run was not pinged")
	}
}
