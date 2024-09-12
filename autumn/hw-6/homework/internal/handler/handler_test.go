package handler_test

import (
	"bytes"
	"fmt"
	"homework/internal/handler"
	"homework/internal/handler/mocks"
	"homework/internal/repository"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
)

type handlerSuite struct {
	suite.Suite
	service *mocks.Service
	handler *handler.Handler
	server  *httptest.Server
	client  *http.Client
}

func (suite *handlerSuite) SetupTest() {
	suite.service = new(mocks.Service)
	suite.handler = handler.New(suite.service)
	mux := chi.NewMux()
	mux.Post("/device", suite.handler.CreateDevice)
	mux.Get("/device/{serialNum}", suite.handler.GetDevice)
	mux.Put("/device", suite.handler.UpdateDevice)
	mux.Delete("/device/{serialNum}", suite.handler.DeleteDevice)
	suite.server = httptest.NewServer(mux)
	suite.client = suite.server.Client()
}

func (suite *handlerSuite) TearDownTest() {
	suite.server.Close()
}

func (suite *handlerSuite) TestCreateDevice() {
	url := fmt.Sprintf("%s/device", suite.server.URL)
	const methodName = "CreateDevice"
	type Test struct {
		name      string
		inBody    string
		inDevice  handler.Device
		returnArg error
		expBody   string
		expStatus int
	}
	testCases := []Test{
		{
			name:   "OK",
			inBody: `{"serial_num":"88","model":"some model","ip":"10.10.10.10"}`,
			inDevice: handler.Device{
				SerialNum: "88",
				Model:     "some model",
				IP:        "10.10.10.10",
			},
			expStatus: http.StatusOK,
		},
		{
			name:   "invalid body",
			inBody: `{"serial_""""}`,
			expBody: fmt.Sprintf(`{"message":"%s"}
`, handler.ErrInvalidDataFormat),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "empty fields",
			inBody: `{}`,
			expBody: fmt.Sprintf(`{"message":"%s: 'serial_num' cannot be empty"}
`, handler.ErrNotEnoughDeviceInfo),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "empty model field",
			inBody: `{"serial_num":"128","ip":"10.10.10.10"}`,
			expBody: fmt.Sprintf(`{"message":"%s: 'model' cannot be empty"}
`, handler.ErrNotEnoughDeviceInfo),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "empty ip field",
			inBody: `{"serial_num":"88","model":"some model"}`,
			expBody: fmt.Sprintf(`{"message":"%s: 'ip' cannot be empty"}
`, handler.ErrNotEnoughDeviceInfo),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "invalid ip format: octet out of range",
			inBody: `{"serial_num":"88","model":"some model","ip":"256.256.256.256"}`,
			expBody: fmt.Sprintf(`{"message":"%s: strconv.ParseUint: parsing \"256\": value out of range"}
`, handler.ErrInvalidIP),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "invalid ip format: incorrect count of octets",
			inBody: `{"serial_num":"88","model":"some model","ip":"224.1.0"}`,
			expBody: fmt.Sprintf(`{"message":"%s"}
`, handler.ErrInvalidIP),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "device already exists",
			inBody: `{"serial_num":"23","model":"some model","ip":"10.10.10.10"}`,
			inDevice: handler.Device{
				SerialNum: "23",
				Model:     "some model",
				IP:        "10.10.10.10",
			},
			returnArg: repository.ErrDeviceExists,
			expBody: fmt.Sprintf(`{"message":"%s"}
`, repository.ErrDeviceExists),
			expStatus: http.StatusBadRequest,
		},
	}
	for _, test := range testCases {
		suite.service.On(methodName, test.inDevice).Return(test.returnArg)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(test.inBody))
		suite.NoError(err, test.name)
		resp, err := suite.client.Do(req)
		suite.NoError(err, test.name)
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		suite.NoError(err, test.name)
		suite.Equal(test.expBody, string(respBody), test.name)
		suite.Equal(test.expStatus, resp.StatusCode, test.name)
	}
}

func (suite *handlerSuite) TestGetDevice() {
	const methodName = "GetDevice"
	type args struct {
		device handler.Device
		err    error
	}
	type Test struct {
		name       string
		serialNum  string
		returnArgs args
		expBody    string
		expStatus  int
	}
	testCases := []Test{
		{
			name:      "OK",
			serialNum: "1",
			returnArgs: args{
				device: handler.Device{
					SerialNum: "1",
					Model:     "some model",
					IP:        "10.10.10.10",
				},
			},
			expBody: `{"serial_num":"1","model":"some model","ip":"10.10.10.10"}
`,
			expStatus: http.StatusOK,
		},
		{
			name:      "device not found",
			serialNum: "13",
			returnArgs: args{
				err: repository.ErrDeviceNotFound,
			},
			expBody: fmt.Sprintf(`{"message":"%s"}
`, repository.ErrDeviceNotFound),
			expStatus: http.StatusNotFound,
		},
	}
	for _, test := range testCases {
		suite.service.On(methodName, test.serialNum).
			Return(test.returnArgs.device, test.returnArgs.err)
		url := fmt.Sprintf("%s/device/%s", suite.server.URL, test.serialNum)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		suite.NoError(err, test.name)
		resp, err := suite.client.Do(req)
		suite.NoError(err, test.name)
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		suite.NoError(err, test.name)
		suite.Equal(test.expBody, string(respBody), test.name)
		suite.Equal(test.expStatus, resp.StatusCode, test.name)
	}
}

func (suite *handlerSuite) TestUpdateDevice() {
	url := fmt.Sprintf("%s/device", suite.server.URL)
	const methodName = "UpdateDevice"
	type Test struct {
		name      string
		inBody    string
		inDevice  handler.Device
		returnArg error
		expBody   string
		expStatus int
	}
	testCases := []Test{
		{
			name:   "OK",
			inBody: `{"serial_num":"322","model":"some model","ip":"10.10.10.10"}`,
			inDevice: handler.Device{
				SerialNum: "322",
				Model:     "some model",
				IP:        "10.10.10.10",
			},
			expStatus: http.StatusOK,
		},
		{
			name:   "invalid body",
			inBody: `{"serial_""""}`,
			expBody: fmt.Sprintf(`{"message":"%s"}
`, handler.ErrInvalidDataFormat),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "empty model field",
			inBody: `{"serial_num":"128","ip":"10.10.10.10"}`,
			expBody: fmt.Sprintf(`{"message":"%s: 'model' cannot be empty"}
`, handler.ErrNotEnoughDeviceInfo),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "empty ip field",
			inBody: `{"serial_num":"88","model":"some model"}`,
			expBody: fmt.Sprintf(`{"message":"%s: 'ip' cannot be empty"}
`, handler.ErrNotEnoughDeviceInfo),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "invalid ip format: octet out of range",
			inBody: `{"serial_num":"88","model":"some model","ip":"256.256.256.256"}`,
			expBody: fmt.Sprintf(`{"message":"%s: strconv.ParseUint: parsing \"256\": value out of range"}
`, handler.ErrInvalidIP),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "invalid ip format: incorrect count of octets",
			inBody: `{"serial_num":"88","model":"some model","ip":"224.1.0"}`,
			expBody: fmt.Sprintf(`{"message":"%s"}
`, handler.ErrInvalidIP),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "device not found",
			inBody: `{"serial_num":"9999","model":"not found","ip":"10.10.10.10"}`,
			inDevice: handler.Device{
				SerialNum: "9999",
				Model:     "not found",
				IP:        "10.10.10.10",
			},
			returnArg: repository.ErrDeviceNotFound,
			expBody: fmt.Sprintf(`{"message":"%s"}
`, repository.ErrDeviceNotFound),
			expStatus: http.StatusNotFound,
		},
	}
	for _, test := range testCases {
		suite.service.On(methodName, test.inDevice).Return(test.returnArg)
		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(test.inBody))
		suite.NoError(err, test.name)
		resp, err := suite.client.Do(req)
		suite.NoError(err, test.name)
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		suite.NoError(err, test.name)
		suite.Equal(test.expBody, string(respBody), test.name)
		suite.Equal(test.expStatus, resp.StatusCode, test.name)
	}
}

func (suite *handlerSuite) TestDeleteDevice() {
	const methodName = "DeleteDevice"
	type Test struct {
		name      string
		serialNum string
		returnArg error
		expBody   string
		expStatus int
	}
	testCases := []Test{
		{
			name:      "OK",
			serialNum: "22",
			expStatus: http.StatusOK,
		},
		{
			name:      "device not found",
			serialNum: "13",
			returnArg: repository.ErrDeviceNotFound,
			expBody: fmt.Sprintf(`{"message":"%s"}
`, repository.ErrDeviceNotFound),
			expStatus: http.StatusNotFound,
		},
	}
	for _, test := range testCases {
		suite.service.On(methodName, test.serialNum).Return(test.returnArg)
		req, err := http.NewRequest(http.MethodDelete,
			fmt.Sprintf("%s/device/%s", suite.server.URL, test.serialNum), nil)
		suite.NoError(err, test.name)
		resp, err := suite.client.Do(req)
		suite.NoError(err, test.name)
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		suite.NoError(err, test.name)
		suite.Equal(test.expBody, string(respBody), test.name)
		suite.Equal(test.expStatus, resp.StatusCode, test.name)
	}
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(handlerSuite))
}
