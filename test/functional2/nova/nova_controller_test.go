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
package nova

import (
	"testing"
	"time"

	base "github.com/openstack-k8s-operators/nova-operator/test/functional2"
	"github.com/stretchr/testify/suite"
)

type NovaSuite struct {
	base.EnvTestSuite
}

func (suite *NovaSuite) TestCreate1() {
	suite.Log().Info("TestCreate1")
	time.Sleep(100 * time.Millisecond)
	suite.Log().Info("TestCreate1 after sleep")
}

func (suite *NovaSuite) TestCreate2() {
	suite.Log().Info("TestCreate2")
	time.Sleep(100 * time.Millisecond)
	suite.Log().Info("TestCreate2 after sleep")
}

func TestNovaSuite(t *testing.T) {
	suite.Run(t, new(NovaSuite))
}
