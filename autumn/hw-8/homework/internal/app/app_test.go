package app_test

import (
	"homework/internal/app"
	"homework/internal/device"
	"homework/internal/device/mocks"
	"testing"

	"github.com/stretchr/testify/suite"
)

type appSuite struct {
	suite.Suite
	repo *mocks.Repository
	app  *app.App
}

func (suite *appSuite) SetupTest() {
	suite.repo = new(mocks.Repository)
	suite.app = app.New(suite.repo)
}

func (suite *appSuite) TestCreateDevice() {
	const methodName = "CreateDevice"
	device := &device.Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}
	suite.repo.On(methodName, device).Return(nil)
	err := suite.app.CreateDevice(device)
	suite.NoError(err)
}

func (suite *appSuite) TestGetDevice() {
	const methodName = "GetDevice"
	device := &device.Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}
	suite.repo.On(methodName, device.SerialNum).Return(device, nil)
	gotDevice, err := suite.app.GetDevice(device.SerialNum)
	suite.NoError(err)
	suite.Equal(device, gotDevice)
}

func (suite *appSuite) TestDeleteDevice() {
	const methodName = "DeleteDevice"
	serialNum := "888"
	suite.repo.On(methodName, serialNum).Return(nil)
	err := suite.app.DeleteDevice(serialNum)
	suite.NoError(err)
}

func (suite *appSuite) TestUpdateDevice() {
	const methodName = "UpdateDevice"
	device := &device.Device{
		SerialNum: "123",
		Model:     "model322",
		IP:        "1.1.1.1",
	}
	suite.repo.On(methodName, device).Return(nil)
	err := suite.app.UpdateDevice(device)
	suite.NoError(err)
}

func TestService(t *testing.T) {
	suite.Run(t, new(appSuite))
}
