package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"homework/models"
	"log"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

type ClientConfig struct {
	Timeout time.Duration
}

type RoundTripperFunc func(r *http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

type tripperware func(rt http.RoundTripper) http.RoundTripper

func NewClient(cfg *ClientConfig, tw ...tripperware) *http.Client {
	var transport = http.DefaultTransport
	for i := len(tw) - 1; i >= 0; i-- {
		transport = tw[i](transport)
	}
	return &http.Client{
		Transport: &loggingRoundTripper{
			next: newBreakerRoundTripper(transport),
		},
		Timeout: cfg.Timeout,
	}
}

type loggingRoundTripper struct {
	next http.RoundTripper
}

func (l *loggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	reqInfo, err := json.Marshal(models.Request{
		Method:      r.Method,
		Endpoint:    r.URL.Path,
		Header:      r.Header,
		QueryParams: r.URL.Query(),
	})
	if err != nil {
		return nil, err
	}
	log.Printf("outgoing request: %s", reqInfo)
	resp, err := l.next.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	respInfo, err := json.Marshal(models.Response{
		Status:  resp.StatusCode,
		Header:  resp.Header,
		Message: getErrMsg(resp),
	})
	if err != nil {
		return nil, err
	}
	log.Printf("incoming response: %s", respInfo)
	return resp, nil
}

func getErrMsg(resp *http.Response) string {
	if resp.StatusCode < http.StatusBadRequest {
		return ""
	}
	var msgBuf []byte
	defer resp.Body.Close()
	_, err := resp.Body.Read(msgBuf)
	if err != nil {
		return err.Error()
	}
	return string(msgBuf)
}

type breakerRoundTripper struct {
	cb   *gobreaker.CircuitBreaker
	next http.RoundTripper
}

var errFailedRequest = errors.New("failed request")

func newBreakerRoundTripper(next http.RoundTripper) *breakerRoundTripper {
	st := gobreaker.Settings{
		Name: "breaker round tripper",
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		Timeout: 30 * time.Second,
	}
	return &breakerRoundTripper{
		next: next,
		cb:   gobreaker.NewCircuitBreaker(st),
	}
}

func (b *breakerRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := b.cb.Execute(func() (any, error) {
		resp, err := b.next.RoundTrip(r)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode >= http.StatusInternalServerError {
			return nil, fmt.Errorf("%w: %s", errFailedRequest, http.StatusText(resp.StatusCode))
		}
		return resp, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*http.Response), nil
}
