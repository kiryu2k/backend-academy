package client

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {}

func TestLoggingRoundTripper(t *testing.T) {
	testCases := []struct {
		name     string
		method   string
		endpoint string
		expReq   string
		expResp  string
	}{
		{
			name:     "get request without query params",
			method:   http.MethodGet,
			endpoint: "/who",
			expReq:   `{"method":"GET","endpoint":"/who","header":{},"query_params":{}}`,
			expResp: `{"status":200,"header":{"Content-Length":["0"]},"message":""}
`,
		},
		{
			name:     "post request with query params",
			method:   http.MethodPost,
			endpoint: "/requestik/woah?data=some_data&test=true",
			expReq:   `{"method":"POST","endpoint":"/requestik/woah","header":{},"query_params":{"data":["some_data"],"test":["true"]}}`,
			expResp: `{"status":200,"header":{"Content-Length":["0"]},"message":""}
`,
		},
	}
	client := &http.Client{
		Transport: &loggingRoundTripper{
			next: http.DefaultTransport,
		},
	}
	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()
	buf := bytes.NewBuffer(nil)
	log.SetOutput(buf)
	for _, test := range testCases {
		req, err := http.NewRequest(test.method, server.URL+test.endpoint, http.NoBody)
		assert.NoError(t, err, test.name)
		_, err = client.Do(req)
		assert.NoError(t, err, test.name)

		logs := strings.SplitN(buf.String(), "\n", 2)
		startIdx := strings.Index(logs[0], "{")
		assert.NotEqual(t, -1, startIdx, test.name)
		assert.Equal(t, test.expReq, logs[0][startIdx:], test.name)
		startIdx = strings.Index(logs[1], "{")
		assert.NotEqual(t, -1, startIdx, test.name)
		withoutDate := regexp.MustCompile(`,"Date":\[.+\]`).ReplaceAllString(logs[1][startIdx:], "")
		assert.Equal(t, test.expResp, withoutDate, test.name)
		buf.Reset()
	}
}

func TestBreakerRoundTripper(t *testing.T) {
	handlerCallCount := 0
	brokenHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		handlerCallCount++
	}
	server := httptest.NewServer(http.HandlerFunc(brokenHandler))
	defer server.Close()
	client := &http.Client{
		Transport: newBreakerRoundTripper(http.DefaultTransport),
	}
	req, err := http.NewRequest(http.MethodGet, server.URL, http.NoBody)
	assert.NoError(t, err)
	for i := 0; i < 3; i++ {
		_, err = client.Do(req)
		assert.ErrorIs(t, err, errFailedRequest)
		assert.Equal(t, i+1, handlerCallCount)
	}
	for i := 0; i < 3; i++ {
		_, err = client.Do(req)
		assert.ErrorIs(t, err, gobreaker.ErrOpenState)
		assert.Equal(t, 3, handlerCallCount)
	}
}

func TestAddCustomRoundTripper(t *testing.T) {
	expKey := "round-tripper-test"
	expVal := "chillin~"
	firstRt := func(rt http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			r.Header.Add(expKey, expVal)
			return rt.RoundTrip(r)
		})
	}
	secondRt := func(rt http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			testMsg := r.Header.Get(expKey)
			r.Body = io.NopCloser(strings.NewReader(testMsg))
			return rt.RoundTrip(r)
		})
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, expVal, string(body))
		assert.Equal(t, expVal, r.Header.Get(expKey))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	client := NewClient(&ClientConfig{}, firstRt, secondRt)
	req, err := http.NewRequest(http.MethodGet, server.URL, http.NoBody)
	assert.NoError(t, err)
	_, err = client.Do(req)
	assert.NoError(t, err)
}

func TestClient(t *testing.T) {
	expBody := `
Lorem ipsum dolor sit amet, consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
`
	client := NewClient(&ClientConfig{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, expBody, string(body))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	req, err := http.NewRequest(http.MethodGet, server.URL, strings.NewReader(expBody))
	assert.NoError(t, err)
	_, err = client.Do(req)
	assert.NoError(t, err)
}
