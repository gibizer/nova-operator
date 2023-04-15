/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package functionalpure_test

import (
	"reflect"
	"testing"
)

// This is the basic test that works
// func TestGlobalFunc(t *testing.T) {
// 	fmt.Print("TestGlobalFunc")
// 	t.Error("TestGlobalFunc")
// }

// This won't work as adding to the signature is not allowed
// func TestGlobalFuncWithParam(t *testing.T, bt BaseTest) {
// 	fmt.Print("TestGlobalFunc")
// 	t.Error("TestGlobalFunc")
// }

// type BaseTest struct {
// }
// This is not discovered as it is not a top level function
// func (bt *BaseTest) TestFuncWithReceiver(t *testing.T) {
// 	fmt.Print("TestFunctionWithReceiver")
// 	t.Error("TestFunctionWithReceiver")
// }

type TestSuite interface {
	SetupSuite(t *testing.T)
	CleanupSuite(t *testing.T)
	Setup(t *testing.T)
	Cleanup(t *testing.T)
	Run(t *testing.T)
}

type BaseSuite struct {
}

func (s *BaseSuite) SetupSuite(t *testing.T) {
	t.Log("BaseSuite.SetupSuite")
}

func (s *BaseSuite) CleanupSuite(t *testing.T) {
	t.Log("BaseSuite.CleanupSuite")
}

func (s *BaseSuite) Setup(t *testing.T) {
	t.Log("BaseSuite.Setup")
}

func (s *BaseSuite) Cleanup(t *testing.T) {
	t.Log("BaseSuite.Cleanup")
}

func (s *BaseSuite) Run(t *testing.T) {
	// run all test cases
	// at this point we don't have access to the actual suite to type hust BaseSuite
	t.Log("BaseSuite.Run", reflect.TypeOf(s))
}

type SuiteA struct {
	BaseSuite
}

func (s *SuiteA) SetupSuite(t *testing.T) {
	t.Log("SuiteA.SetupSuite")
}

func (s *SuiteA) CleanupSuite(t *testing.T) {
	t.Log("SuiteA.CleanupSuite")
}

func (s *SuiteA) Setup(t *testing.T) {
	t.Log("SuiteA.Setup")
}

func (s *SuiteA) Cleanup(t *testing.T) {
	t.Log("SuiteA.Cleanup")
}

func (s *SuiteA) Test1(t *testing.T) {
	t.Log("SuiteA.Test1")
	// t.Error("TestSuiteA.Test1")
}

func (s *SuiteA) Test2(t *testing.T) {
	t.Log("SuiteA.Test2")
	// t.Error("TestSuiteA.Test2")
}

func TestSuites(t *testing.T) {
	RunTestSuite(t, &SuiteA{})
}

func RunTestSuite(t *testing.T, suite TestSuite) {
	defer suite.Cleanup(t)
	suite.SetupSuite(t)
	suite.Run(t)
}
