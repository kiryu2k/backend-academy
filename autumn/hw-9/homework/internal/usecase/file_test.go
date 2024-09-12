package usecase_test

import (
	"context"
	"homework/internal/domain"
	"homework/internal/domain/mocks"
	"homework/internal/usecase"
	"testing"

	"github.com/stretchr/testify/suite"
)

type fileSuite struct {
	suite.Suite
	repo *mocks.FileRepo
	file *usecase.File
}

func (suite *fileSuite) SetupTest() {
	suite.repo = new(mocks.FileRepo)
	suite.file = usecase.New(suite.repo)
}

func (suite *fileSuite) TestFileGet() {
	testCases := []struct {
		name     string
		filename string
		exp      *domain.FileInfo
		err      error
	}{
		{
			name:     "OK",
			filename: "file.txt",
			exp: &domain.FileInfo{
				Name: "file.txt",
				Type: ".txt",
				Size: 22,
				Data: []byte("some~file~data"),
			},
		},
		{
			name:     "not found",
			filename: "not_exist_file.exe",
			err:      domain.ErrFileNotFound,
		},
	}
	const methodName = "Find"
	for _, test := range testCases {
		suite.repo.On(methodName, context.Background(), test.filename).Return(test.exp, test.err)
		actual, err := suite.repo.Find(context.Background(), test.filename)
		suite.Require().ErrorIs(test.err, err, test.name)
		suite.Require().Equal(test.exp, actual, test.name)
	}
}

func (suite *fileSuite) TestFileAll() {
	const methodName = "All"
	exp := []string{
		"file1.txt",
		"file2.jpeg",
		"file3.exe",
	}
	suite.repo.On(methodName, context.Background()).Return(exp)
	actual := suite.repo.All(context.Background())
	suite.Require().Equal(exp, actual)
}

func TestFile(t *testing.T) {
	suite.Run(t, new(fileSuite))
}
