package pinger

import (
	"context"
	"time"
)

type (
	// Pinger is behavior of ping
	Pinger interface {
		PingContext(ctx context.Context) error
	}

	// WorkerOption use for Worker
	WorkerOption func(*workerOptions)

	errorHandler  func(error)
	workerOptions struct {
		interval     time.Duration
		errorHandler errorHandler
	}
)

const defaultCheckInterval = 10 * time.Second

// NewWorker returns new configured Worker
func NewWorker(client Pinger, opts ...WorkerOption) *Worker {
	options := &workerOptions{
		interval:     defaultCheckInterval,
		errorHandler: nil,
	}
	for _, o := range opts {
		o(options)
	}

	return &Worker{
		client:       client,
		interval:     options.interval,
		errorHandler: options.errorHandler,
	}
}

// WithInterval is configure ping interval
func WithInterval(t time.Duration) WorkerOption {
	return func(w *workerOptions) {
		w.interval = t
	}
}

// WithErrorHandler is configure error handling for Pinger.PingContext
func WithErrorHandler(handler func(error)) WorkerOption {
	return func(w *workerOptions) {
		w.errorHandler = handler
	}
}

// Worker for pinger
type Worker struct {
	client       Pinger
	interval     time.Duration
	errorHandler errorHandler

	ctx  context.Context
	stop context.CancelFunc
}

// Run pinger worker
func (d *Worker) Run(ctx context.Context) error {
	d.ctx, d.stop = context.WithCancel(ctx)

	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err := d.client.PingContext(ctx)
			if err != nil && d.errorHandler != nil {
				d.errorHandler(err)
			}
		}
	}
}

// Stop pinger worker
func (d *Worker) Stop() {
	if d.stop != nil {
		d.stop()
	}
}
