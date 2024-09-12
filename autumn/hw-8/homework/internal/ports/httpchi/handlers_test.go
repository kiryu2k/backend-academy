package httpchi_test

import (
	"bytes"
	"fmt"
	"homework/internal/adapters/devicerepo"
	"homework/internal/device"
	"homework/internal/device/mocks"
	"homework/internal/ports/httpchi"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
)

type handlerSuite struct {
	suite.Suite
	usecase *mocks.Usecase
	server  *httptest.Server
	client  *http.Client
}

func (suite *handlerSuite) SetupTest() {
	suite.usecase = new(mocks.Usecase)
	mux := chi.NewMux()
	httpchi.AppRouter(mux, suite.usecase)
	suite.server = httptest.NewServer(mux)
	suite.client = suite.server.Client()
}

func (suite *handlerSuite) TearDownTest() {
	suite.server.Close()
}

func (suite *handlerSuite) TestCreateDevice() {
	testCases := []struct {
		name      string
		inBody    string
		inDevice  *device.Device
		returnArg error
		expBody   string
		expStatus int
	}{
		{
			name:   "OK",
			inBody: `{"serial_num":"88","model":"some model","ip":"10.10.10.10"}`,
			inDevice: &device.Device{
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
`, httpchi.ErrInvalidDataFormat),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "empty fields",
			inBody: `{}`,
			expBody: fmt.Sprintf(`{"message":"%s: 'serial_num' cannot be empty"}
`, httpchi.ErrNotEnoughDeviceInfo),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "empty model field",
			inBody: `{"serial_num":"128","ip":"10.10.10.10"}`,
			expBody: fmt.Sprintf(`{"message":"%s: 'model' cannot be empty"}
`, httpchi.ErrNotEnoughDeviceInfo),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "empty ip field",
			inBody: `{"serial_num":"88","model":"some model"}`,
			expBody: fmt.Sprintf(`{"message":"%s: 'ip' cannot be empty"}
`, httpchi.ErrNotEnoughDeviceInfo),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "invalid ip format: octet out of range",
			inBody: `{"serial_num":"88","model":"some model","ip":"256.256.256.256"}`,
			expBody: fmt.Sprintf(`{"message":"%s: strconv.ParseUint: parsing \"256\": value out of range"}
`, httpchi.ErrInvalidIP),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "invalid ip format: incorrect count of octets",
			inBody: `{"serial_num":"88","model":"some model","ip":"224.1.0"}`,
			expBody: fmt.Sprintf(`{"message":"%s"}
`, httpchi.ErrInvalidIP),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "device already exists",
			inBody: `{"serial_num":"23","model":"some model","ip":"10.10.10.10"}`,
			inDevice: &device.Device{
				SerialNum: "23",
				Model:     "some model",
				IP:        "10.10.10.10",
			},
			returnArg: devicerepo.ErrDeviceExists,
			expBody: fmt.Sprintf(`{"message":"%s"}
`, devicerepo.ErrDeviceExists),
			expStatus: http.StatusBadRequest,
		},
	}
	url := suite.server.URL + "/device"
	const methodName = "CreateDevice"
	for _, test := range testCases {
		suite.usecase.On(methodName, test.inDevice).Return(test.returnArg)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(test.inBody))
		suite.NoError(err, test.name)
		resp, err := suite.client.Do(req)
		suite.NoError(err, test.name)
		respBody, err := io.ReadAll(resp.Body)
		suite.NoError(err, test.name)
		suite.Equal(test.expBody, string(respBody), test.name)
		suite.Equal(test.expStatus, resp.StatusCode, test.name)
		resp.Body.Close()
	}
}

func (suite *handlerSuite) TestGetDevice() {
	type args struct {
		device *device.Device
		err    error
	}
	testCases := []struct {
		name       string
		serialNum  string
		returnArgs args
		expBody    string
		expStatus  int
	}{
		{
			name:      "OK",
			serialNum: "1",
			returnArgs: args{
				device: &device.Device{
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
				err: devicerepo.ErrDeviceNotFound,
			},
			expBody: fmt.Sprintf(`{"message":"%s"}
`, devicerepo.ErrDeviceNotFound),
			expStatus: http.StatusNotFound,
		},
	}
	const methodName = "GetDevice"
	for _, test := range testCases {
		suite.usecase.On(methodName, test.serialNum).
			Return(test.returnArgs.device, test.returnArgs.err)
		url := fmt.Sprintf("%s/device/%s", suite.server.URL, test.serialNum)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		suite.NoError(err, test.name)
		resp, err := suite.client.Do(req)
		suite.NoError(err, test.name)
		respBody, err := io.ReadAll(resp.Body)
		suite.NoError(err, test.name)
		suite.Equal(test.expBody, string(respBody), test.name)
		suite.Equal(test.expStatus, resp.StatusCode, test.name)
		resp.Body.Close()
	}
}

func (suite *handlerSuite) TestUpdateDevice() {
	testCases := []struct {
		name      string
		inBody    string
		inDevice  *device.Device
		returnArg error
		expBody   string
		expStatus int
	}{
		{
			name:   "OK",
			inBody: `{"serial_num":"322","model":"some model","ip":"10.10.10.10"}`,
			inDevice: &device.Device{
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
`, httpchi.ErrInvalidDataFormat),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "empty model field",
			inBody: `{"serial_num":"128","ip":"10.10.10.10"}`,
			expBody: fmt.Sprintf(`{"message":"%s: 'model' cannot be empty"}
`, httpchi.ErrNotEnoughDeviceInfo),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "empty ip field",
			inBody: `{"serial_num":"88","model":"some model"}`,
			expBody: fmt.Sprintf(`{"message":"%s: 'ip' cannot be empty"}
`, httpchi.ErrNotEnoughDeviceInfo),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "invalid ip format: octet out of range",
			inBody: `{"serial_num":"88","model":"some model","ip":"256.256.256.256"}`,
			expBody: fmt.Sprintf(`{"message":"%s: strconv.ParseUint: parsing \"256\": value out of range"}
`, httpchi.ErrInvalidIP),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "invalid ip format: incorrect count of octets",
			inBody: `{"serial_num":"88","model":"some model","ip":"224.1.0"}`,
			expBody: fmt.Sprintf(`{"message":"%s"}
`, httpchi.ErrInvalidIP),
			expStatus: http.StatusBadRequest,
		},
		{
			name:   "device not found",
			inBody: `{"serial_num":"9999","model":"not found","ip":"10.10.10.10"}`,
			inDevice: &device.Device{
				SerialNum: "9999",
				Model:     "not found",
				IP:        "10.10.10.10",
			},
			returnArg: devicerepo.ErrDeviceNotFound,
			expBody: fmt.Sprintf(`{"message":"%s"}
`, devicerepo.ErrDeviceNotFound),
			expStatus: http.StatusNotFound,
		},
	}
	url := suite.server.URL + "/device"
	const methodName = "UpdateDevice"
	for _, test := range testCases {
		suite.usecase.On(methodName, test.inDevice).Return(test.returnArg)
		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(test.inBody))
		suite.NoError(err, test.name)
		resp, err := suite.client.Do(req)
		suite.NoError(err, test.name)
		respBody, err := io.ReadAll(resp.Body)
		suite.NoError(err, test.name)
		suite.Equal(test.expBody, string(respBody), test.name)
		suite.Equal(test.expStatus, resp.StatusCode, test.name)
		resp.Body.Close()
	}
}

func (suite *handlerSuite) TestDeleteDevice() {
	testCases := []struct {
		name      string
		serialNum string
		returnArg error
		expBody   string
		expStatus int
	}{
		{
			name:      "OK",
			serialNum: "22",
			expStatus: http.StatusOK,
		},
		{
			name:      "device not found",
			serialNum: "13",
			returnArg: devicerepo.ErrDeviceNotFound,
			expBody: fmt.Sprintf(`{"message":"%s"}
`, devicerepo.ErrDeviceNotFound),
			expStatus: http.StatusNotFound,
		},
	}
	const methodName = "DeleteDevice"
	for _, test := range testCases {
		suite.usecase.On(methodName, test.serialNum).Return(test.returnArg)
		req, err := http.NewRequest(http.MethodDelete,
			fmt.Sprintf("%s/device/%s", suite.server.URL, test.serialNum), nil)
		suite.NoError(err, test.name)
		resp, err := suite.client.Do(req)
		suite.NoError(err, test.name)
		respBody, err := io.ReadAll(resp.Body)
		suite.NoError(err, test.name)
		suite.Equal(test.expBody, string(respBody), test.name)
		suite.Equal(test.expStatus, resp.StatusCode, test.name)
		resp.Body.Close()
	}
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(handlerSuite))
}
