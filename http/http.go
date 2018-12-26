package http

import (
	"context"
	"net/http"
	"net/url"
)

type (
	// Pinger is HTTP request pinger
	Pinger struct {
		url               *url.URL
		method            string
		client            *http.Client
		headerMiddlewares []HeaderMiddleware
	}

	// PingerOption use for Pinger
	PingerOption func(*pingerOptions)

	// HeaderMiddleware can injection ping header
	HeaderMiddleware func(*http.Header)

	pingerOptions struct {
		method            string
		client            *http.Client
		headerMiddlewares []HeaderMiddleware
	}
)

// NewPinger returns new configured Pinger
func NewPinger(u *url.URL, opts ...PingerOption) *Pinger {
	options := &pingerOptions{
		method: "OPTION",
		client: http.DefaultClient,
	}
	for _, o := range opts {
		o(options)
	}

	return &Pinger{
		url:               u,
		method:            options.method,
		client:            options.client,
		headerMiddlewares: options.headerMiddlewares,
	}
}

// WithMethod is configure using HTTP method when ping request
func WithMethod(method string) PingerOption {
	return func(h *pingerOptions) {
		h.method = method
	}
}

// WithClient is configure custom your http.Client
func WithClient(c *http.Client) PingerOption {
	return func(h *pingerOptions) {
		h.client = c
	}
}

// WithHeader is configure using HTTP header when ping request
func WithHeader(key string, value string) PingerOption {
	return func(h *pingerOptions) {
		h.headerMiddlewares = append(
			h.headerMiddlewares,
			func(h *http.Header) {
				h.Add(key, value)
			},
		)
	}
}

// PingContext to any url with HTTP method
func (p *Pinger) PingContext(ctx context.Context) error {
	req, err := http.NewRequest(p.method, p.url.String(), nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	for _, o := range p.headerMiddlewares {
		o(&req.Header)
	}

	resp, doErr := p.client.Do(req)
	if doErr != nil {
		return doErr
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}

	return nil
}
