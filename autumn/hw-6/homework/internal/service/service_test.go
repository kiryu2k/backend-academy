package service_test

import (
	"homework/internal/handler"
	"homework/internal/service"
	"homework/internal/service/mocks"
	"testing"

	"github.com/stretchr/testify/suite"
)

type serviceSuite struct {
	suite.Suite
	repo    *mocks.Repository
	service *service.Service
}

func (suite *serviceSuite) SetupTest() {
	suite.repo = new(mocks.Repository)
	suite.service = service.New(suite.repo)
}

func (suite *serviceSuite) TestCreateDevice() {
	const methodName = "CreateDevice"
	device := handler.Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}
	suite.repo.On(methodName, device).Return(nil)
	err := suite.service.CreateDevice(device)
	suite.NoError(err)
}

func (suite *serviceSuite) TestGetDevice() {
	const methodName = "GetDevice"
	device := handler.Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}
	suite.repo.On(methodName, device.SerialNum).Return(device, nil)
	gotDevice, err := suite.service.GetDevice(device.SerialNum)
	suite.NoError(err)
	suite.Equal(device, gotDevice)
}

func (suite *serviceSuite) TestDeleteDevice() {
	const methodName = "DeleteDevice"
	serialNum := "888"
	suite.repo.On(methodName, serialNum).Return(nil)
	err := suite.service.DeleteDevice(serialNum)
	suite.NoError(err)
}

func (suite *serviceSuite) TestUpdateDevice() {
	const methodName = "UpdateDevice"
	device := handler.Device{
		SerialNum: "123",
		Model:     "model322",
		IP:        "1.1.1.1",
	}
	suite.repo.On(methodName, device).Return(nil)
	err := suite.service.UpdateDevice(device)
	suite.NoError(err)
}

func TestService(t *testing.T) {
	suite.Run(t, new(serviceSuite))
}
