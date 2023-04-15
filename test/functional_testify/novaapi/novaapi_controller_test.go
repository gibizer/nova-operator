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
package novaapi

import (
	"testing"
	"time"

	functional "github.com/openstack-k8s-operators/nova-operator/test/functional_testify"
	"github.com/stretchr/testify/suite"
)

type NovaAPISuite struct {
	functional.EnvTestSuite
}

func (suite *NovaAPISuite) SetupSuite() {
	suite.T().Log("!!!SetupSuite!!")
	suite.EnvTestSuite.SetupSuite()
}

func (suite *NovaAPISuite) SetupTest() {
	suite.T().Log("SetupTest")
}

func (suite *NovaAPISuite) TearDownTest() {
	suite.T().Log("TearDownTest")
}

func (suite *NovaAPISuite) TestCreate() {
	suite.Log().Info("TestCreate1")
	time.Sleep(100 * time.Millisecond)
	suite.Log().Info("TestCreate1 after sleep")
}

func (suite *NovaAPISuite) TestCreate2() {
	suite.Log().Info("TestCreate2")
	time.Sleep(100 * time.Millisecond)
	suite.Log().Info("TestCreate2 after sleep")
	suite.Fail("NovaAPI.TestCreate failed")
}

func TestNovaAPISuite(t *testing.T) {
	suite.Run(t, new(NovaAPISuite))
}
