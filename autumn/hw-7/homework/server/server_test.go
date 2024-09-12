package server

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	testCases := []struct {
		name     string
		method   string
		endpoint string
		handler  func(w http.ResponseWriter, r *http.Request)
		expReq   string
		expResp  string
	}{
		{
			name:     "get request without query params",
			method:   http.MethodGet,
			endpoint: "/who",
			handler:  func(w http.ResponseWriter, r *http.Request) {},
			expReq:   `{"method":"GET","endpoint":"/who","header":{"Accept-Encoding":["gzip"],"User-Agent":["Go-http-client/1.1"]},"query_params":{}}`,
			expResp: `{"status":200,"header":{},"message":""}
`,
		},
		{
			name:     "post request with query params",
			method:   http.MethodPost,
			endpoint: "/requestik/woah?data=some_data&test=true",
			handler:  func(w http.ResponseWriter, r *http.Request) {},
			expReq:   `{"method":"POST","endpoint":"/requestik/woah","header":{"Accept-Encoding":["gzip"],"Content-Length":["0"],"User-Agent":["Go-http-client/1.1"]},"query_params":{"data":["some_data"],"test":["true"]}}`,
			expResp: `{"status":200,"header":{},"message":""}
`,
		},
		{
			name:     "bad request",
			method:   http.MethodDelete,
			endpoint: "/endpoint",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("some error message"))
			},
			expReq: `{"method":"DELETE","endpoint":"/endpoint","header":{"Accept-Encoding":["gzip"],"User-Agent":["Go-http-client/1.1"]},"query_params":{}}`,
			expResp: `{"status":400,"header":{},"message":"some error message"}
`,
		},
		{
			name:     "response with json content-type",
			method:   http.MethodPost,
			endpoint: "/json4ik",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
			},
			expReq: `{"method":"POST","endpoint":"/json4ik","header":{"Accept-Encoding":["gzip"],"Content-Length":["0"],"User-Agent":["Go-http-client/1.1"]},"query_params":{}}`,
			expResp: `{"status":200,"header":{"Content-Type":["application/json"]},"message":""}
`,
		},
	}
	buf := bytes.NewBuffer(nil)
	log.SetOutput(buf)
	for _, test := range testCases {
		server := httptest.NewServer(loggingMiddleware(http.HandlerFunc(test.handler)))
		defer server.Close()
		req, err := http.NewRequest(test.method, server.URL+test.endpoint, http.NoBody)
		assert.NoError(t, err, test.name)
		_, err = server.Client().Do(req)
		assert.NoError(t, err, test.name)

		logs := strings.SplitN(buf.String(), "\n", 2)
		startIdx := strings.Index(logs[0], "{")
		assert.NotEqual(t, -1, startIdx, test.name)
		assert.Equal(t, test.expReq, logs[0][startIdx:], test.name)
		startIdx = strings.Index(logs[1], "{")
		assert.NotEqual(t, -1, startIdx, test.name)
		assert.Equal(t, test.expResp, logs[1][startIdx:], test.name)
		buf.Reset()
	}
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {}

const correctAuthHeaderExample = "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" // Aladdin:open sesame (base64 encoded)

func TestBasicAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name       string
		authHeader string
		expStatus  int
		expHeader  string
		expBody    string
	}{
		{
			name:      "empty auth header",
			expStatus: http.StatusUnauthorized,
			expHeader: `Basic realm="SomeRealm"`,
			expBody:   errNoAuthHeader.Error(),
		},
		{
			name:       "invalid auth header format",
			authHeader: "invalid header",
			expStatus:  http.StatusUnauthorized,
			expHeader:  `Basic realm="SomeRealm"`,
			expBody:    errInvalidAuthHeader.Error(),
		},
		{
			name:       "auth header without scheme",
			authHeader: "QWxhZGRpbjpvcGVuIHNlc2FtZQ==",
			expStatus:  http.StatusUnauthorized,
			expHeader:  `Basic realm="SomeRealm"`,
			expBody:    errInvalidAuthHeader.Error(),
		},
		{
			name:       "username and password aren't combined with colon",
			authHeader: "Basic QWxhZGRpbm9wZW4gc2VzYW1l",
			expStatus:  http.StatusUnauthorized,
			expHeader:  `Basic realm="SomeRealm"`,
			expBody:    errInvalidAuthHeader.Error(),
		},
		{
			name:       "ok",
			authHeader: correctAuthHeaderExample,
			expStatus:  http.StatusOK,
		},
	}
	server := httptest.NewServer(basicAuthMiddleware(http.HandlerFunc(dummyHandler)))
	defer server.Close()
	for _, test := range testCases {
		req, err := http.NewRequest(http.MethodGet, server.URL, http.NoBody)
		assert.NoError(t, err, test.name)
		if test.authHeader != "" {
			req.Header.Add("Authorization", test.authHeader)
		}
		resp, err := server.Client().Do(req)
		assert.NoError(t, err, test.name)
		assert.Equal(t, test.expStatus, resp.StatusCode, test.name)
		assert.Equal(t, test.expHeader, resp.Header.Get("WWW-Authenticate"), test.name)
		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		assert.NoError(t, err, test.name)
		assert.Equal(t, test.expBody, string(body), test.name)
	}
}

func TestAddCustomMiddleware(t *testing.T) {
	firstMw := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("first: entering...\n"))
			h.ServeHTTP(w, r)
			_, _ = w.Write([]byte("first: quiting...\n"))
		})
	}
	secondMw := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("second: entering...\n"))
			h.ServeHTTP(w, r)
			_, _ = w.Write([]byte("second: quiting...\n"))
		})
	}
	var (
		cfg     = &ServerConfig{}
		handler = NewServer(cfg, http.HandlerFunc(dummyHandler), firstMw, secondMw).Handler
		server  = httptest.NewServer(handler)
	)
	defer server.Close()
	req, err := http.NewRequest(http.MethodGet, server.URL, http.NoBody)
	assert.NoError(t, err)
	req.Header.Add("Authorization", correctAuthHeaderExample)
	resp, err := server.Client().Do(req)
	assert.NoError(t, err)
	exp := `first: entering...
second: entering...
second: quiting...
first: quiting...
`
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, exp, string(body))
}

func TestServer(t *testing.T) {
	expBody := `
	Lorem ipsum dolor sit amet, consectetur adipiscing elit,
	sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
	`
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, expBody, string(body))
		_, _ = w.Write(body)
	}
	server := NewServer(&ServerConfig{}, http.HandlerFunc(handler))
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, server.Addr, strings.NewReader(expBody))
	assert.NoError(t, err)
	req.Header.Add("Authorization", correctAuthHeaderExample)
	server.Handler.ServeHTTP(recorder, req)
	assert.Equal(t, expBody, recorder.Body.String())
}
