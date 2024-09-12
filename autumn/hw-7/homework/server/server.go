package server

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"homework/models"
	"log"
	"net/http"
	"strings"
	"time"
)

type ServerConfig struct {
	Addr         string        `yaml:"address"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	Readimeout   time.Duration `yaml:"read_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type middleware func(h http.Handler) http.Handler

func NewServer(cfg *ServerConfig, mux http.Handler, mw ...middleware) *http.Server {
	for i := len(mw) - 1; i >= 0; i-- {
		mux = mw[i](mux)
	}
	return &http.Server{
		Addr:         cfg.Addr,
		Handler:      loggingMiddleware(basicAuthMiddleware(mux)),
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.Readimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := json.Marshal(models.Request{
			Method:      r.Method,
			Endpoint:    r.URL.Path,
			Header:      r.Header,
			QueryParams: r.URL.Query(),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("incoming request: %s", req)
		s := &responseSaver{
			status: 200,
			w:      w,
		}
		h.ServeHTTP(s, r)
		resp, err := json.Marshal(models.Response{
			Status:  s.status,
			Header:  s.Header(),
			Message: s.errMsg,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("outgoing response: %s", resp)
	})
}

var (
	errNoAuthHeader      = errors.New("no authorization header")
	errInvalidAuthHeader = errors.New("invalid authorization header")
)

func basicAuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			writeAuthErrorResponse(w, errNoAuthHeader)
			return
		}
		tokenParts := strings.Split(token, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Basic" {
			writeAuthErrorResponse(w, errInvalidAuthHeader)
			return
		}
		userData, err := base64.StdEncoding.DecodeString(tokenParts[1])
		if err != nil {
			writeAuthErrorResponse(w, errInvalidAuthHeader)
			return
		}
		if !strings.Contains(string(userData), ":") {
			writeAuthErrorResponse(w, errInvalidAuthHeader)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func writeAuthErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Add("WWW-Authenticate", `Basic realm="SomeRealm"`)
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(err.Error()))
}

type responseSaver struct {
	w      http.ResponseWriter
	status int
	errMsg string
}

func (r *responseSaver) Header() http.Header {
	return r.w.Header()
}

func (r *responseSaver) Write(bytes []byte) (int, error) {
	if r.status >= http.StatusBadRequest {
		r.errMsg = string(bytes)
	}
	n, err := r.w.Write(bytes)
	return n, err
}

func (r *responseSaver) WriteHeader(statusCode int) {
	r.status = statusCode
	r.w.WriteHeader(statusCode)
}
