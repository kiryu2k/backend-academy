package models

import (
	"net/http"
	"net/url"
)

type Request struct {
	Method      string      `json:"method"`
	Endpoint    string      `json:"endpoint"`
	Header      http.Header `json:"header"`
	QueryParams url.Values  `json:"query_params"`
}

type Response struct {
	Status  int         `json:"status"`
	Header  http.Header `json:"header"`
	Message string      `json:"message"`
}
