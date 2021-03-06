package MethodCallRetrier

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RetrierTestSuite struct {
	suite.Suite

	retrier *MethodCallRetrier
}

func (s *RetrierTestSuite) SetupTest() {
	s.retrier = New(0, 1, nil)
}

func TestRetrierTestSuite(t *testing.T) {
	suite.Run(t, new(RetrierTestSuite))
}

func (s *RetrierTestSuite) TestRetrierWorksWithPointer() {
	arg := "TestArg"

	results, _ := s.retrier.ExecuteWithRetry(&RetryObject{}, "MethodReturningString", arg)

	s.Assert().EqualValues(results[0].String(), arg)
}

func (s *RetrierTestSuite) TestRetrierWorksWithObject() {
	arg := "TestArg"

	results, _ := s.retrier.ExecuteWithRetry(RetryObject{}, "MethodReturningString", arg)

	s.Assert().EqualValues(results[0].String(), arg)
}

func (s *RetrierTestSuite) TestRetrierThrowsErrorReturnsNilResults() {
	results, _ := s.retrier.ExecuteWithRetry(RetryObject{}, "MethodReturningError", "TestArg")

	s.Assert().Nil(results)
}

func (s *RetrierTestSuite) TestRetrierThrowsErrorReturnsErrors() {
	_, errs := s.retrier.ExecuteWithRetry(RetryObject{}, "MethodReturningError", "TestArg")

	s.Assert().IsType(errors.New(""), errs[0])
}

func (s *RetrierTestSuite) TestRetrierThrowsErrorReturnsCorrectNumberOfErrors() {
	_, errs := s.retrier.ExecuteWithRetry(RetryObject{}, "MethodReturningError", "TestArg")

	s.Assert().Len(errs, 2)
}

func (s *RetrierTestSuite) TestRetrierReturnsNilWhenGivenObjectWithNoReturnTypes() {
	results, _ := s.retrier.ExecuteWithRetry(RetryObject{}, "MethodReturningNoValues")

	s.Assert().Len(results, 0)
}

func (s *RetrierTestSuite) TestRetrierRetriesCorrectNumberOfTimes() {
	testObj := RetryMockObject{}
	methodName := "MethodReturningError"

	testObj.On(methodName, "").Return(errors.New(""))

	_, _ = New(0, 5, nil).ExecuteWithRetry(&testObj, methodName, "")

	testObj.AssertNumberOfCalls(s.T(), methodName, 5)

	testObj.AssertExpectations(s.T())
}

func (s *RetrierTestSuite) TestRetrierReturnsAllErrorsPlusOurError() {
	testObj := RetryMockObject{}
	methodName := "MethodReturningError"

	testObj.On(methodName, "").Return(errors.New(""))

	_, errs := New(0, 5, nil).ExecuteWithRetry(&testObj, methodName, "")

	s.Assert().Len(errs, 6)
}

func (s *RetrierTestSuite) TestRetrierWorksWhenErrorIsNotLastReturnParamOnObject() {
	testObj := RetryObject{}
	methodName := "MethodReturningErrorInRandomPosition"

	_, errs := New(0, 5, nil).ExecuteWithRetry(&testObj, methodName, "")

	s.Assert().IsType(errors.New(""), errs[0])
}

func (s *RetrierTestSuite) TestRetrierWorksWhenMultipleReturnParamsAreErrors() {
	testObj := RetryObject{}
	methodName := "MethodReturningMultipleErrors"

	_, errs := New(0, 5, nil).ExecuteWithRetry(&testObj, methodName, "")

	s.Assert().Len(errs, 11)
}

type RetryObject struct{}

func (m *RetryObject) MethodReturningNoValues() {}

func (m *RetryObject) MethodReturningString(anArgument string) string {
	return anArgument
}

func (m *RetryObject) MethodReturningError(anArgument string) error {
	return errors.New("")
}

func (m *RetryObject) MethodReturningErrorInRandomPosition() (string, error, string) {
	return "", errors.New(""), ""
}

func (m *RetryObject) MethodReturningMultipleErrors() (string, error, error) {
	return "", errors.New(""), errors.New("")
}

type RetryMockObject struct {
	mock.Mock
}

func (m *RetryMockObject) MethodReturningError(anArgument string) error {
	return m.Called(anArgument).Error(0)
}
