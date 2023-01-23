package pinger

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type fakePinger struct {
	handler func() error
}

func (m *fakePinger) PingContext(ctx context.Context) error {
	return m.handler()
}

func TestWithInterval(t *testing.T) {
	t.Parallel()

	want := 10 * time.Second

	w := &workerOptions{}
	if w.interval == want {
		t.Errorf("already configured interval, want = %d, got = %d", want, w.interval)
	}

	WithInterval(want)(w)
	if w.interval != want {
		t.Errorf("not configured interval, want = %d, got = %d", want, w.interval)
	}
}

func TestWithErrorHandler(t *testing.T) {
	t.Parallel()

	w := &workerOptions{}
	if w.errorHandler != nil {
		t.Errorf("already configured error handler")
	}

	WithErrorHandler(func(error) {})(w)
	if w.errorHandler == nil {
		t.Errorf("not configured error handler")
	}
}

func TestNewWorker(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                    string
		pingError               bool
		wantError               bool
		wantTriggerErrorHandler bool
	}{
		{
			"ping noerror",
			false,
			false,
			false,
		},
		{
			"ping error",
			true,
			false,
			true,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			triggeredErrorHandler := false

			worker := NewWorker(
				&fakePinger{
					func() error {
						if test.pingError {
							return fmt.Errorf("dummy")
						}

						return nil
					},
				},
				WithInterval(time.Millisecond),
				WithErrorHandler(func(err error) {
					triggeredErrorHandler = true
				}),
			)
			defer worker.Stop()

			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Millisecond))
			defer cancel()
			if err := worker.Run(ctx); (err != nil) != test.wantError {
				t.Fatalf("worker.Run want error = %v, got error: %+v", test.wantError, err)
			}

			if test.wantTriggerErrorHandler != triggeredErrorHandler {
				t.Fatalf(
					"worker.Run want error handler trigger = %v, got error handler trigger = %v",
					test.wantTriggerErrorHandler,
					triggeredErrorHandler,
				)
			}
		})
	}
}
