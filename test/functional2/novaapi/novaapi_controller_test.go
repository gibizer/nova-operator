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

	"github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	common_helpers "github.com/openstack-k8s-operators/lib-common/modules/test/helpers"

	base "github.com/openstack-k8s-operators/nova-operator/test/functional2"
	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	. "github.com/onsi/gomega"
)

type NovaAPISuite struct {
	base.EnvTestSuite
}

func (suite *NovaAPISuite) SetupTest() {
	suite.EnvTestSuite.SetupTest()
	suite.DeferCleanup(func() {
		suite.K8sClient.Delete(
			suite.Ctx, suite.CreateNovaMessageBusSecret(suite.Namespace, base.MessageBusSecretName))
	})

}

func (suite *NovaAPISuite) TestCreateMissingSecret() {
	// Create NovaAPI
	spec := suite.GetDefaultNovaAPISpec()
	spec["secret"] = "non-existent-secret"
	raw := suite.CreateNovaAPI(suite.Namespace, spec)
	novaAPIName := types.NamespacedName{Name: raw.GetName(), Namespace: raw.GetNamespace()}
	suite.DeferCleanup(func() { suite.DeleteInstance(raw) })

	// Read the NovaAPI back and assert that it reports the missing secret
	instance := suite.GetNovaAPI(novaAPIName)
	suite.ExpectConditionWithDetails(
		novaAPIName,
		common_helpers.ConditionGetterFunc(suite.NovaAPIConditionGetter),
		condition.InputReadyCondition,
		corev1.ConditionFalse,
		condition.RequestedReason,
		"Input data resources missing: secret/non-existent-secret",
	)
	suite.ExpectCondition(
		novaAPIName,
		common_helpers.ConditionGetterFunc(suite.NovaAPIConditionGetter),
		condition.ReadyCondition,
		corev1.ConditionFalse,
	)

	// Also assert that it does not report any other status

	// NOTE(gibi): Hash and Endpoints have `omitempty` tags so while
	// they are initialized to {} that value is omitted from the output
	// when sent to the client. So we see nils here.
	Expect(instance.Status.Hash).To(BeEmpty())
	Expect(instance.Status.APIEndpoints).To(BeEmpty())
	Expect(instance.Status.ReadyCount).To(Equal(int32(0)))
	Expect(instance.Status.ServiceID).To(Equal(""))
}

func (suite *NovaAPISuite) TestCreateSecretMissingField() {
	suite.Log().Info("Namespace:" + suite.Namespace)
}

func (suite *NovaAPISuite) TestCreate() {
	suite.Log().Info("Namespace:" + suite.Namespace)
}

func (suite *NovaAPISuite) TestDelete() {
	suite.Log().Info("Namespace:" + suite.Namespace)
}

func (suite *NovaAPISuite) TestCreateWithNetworkAttachments() {
	suite.Log().Info("Namespace:" + suite.Namespace)
}

func (suite *NovaAPISuite) TestCreateWithNetworkAttachmentsDefinitionMissing() {
	suite.Log().Info("Namespace:" + suite.Namespace)
}

func (suite *NovaAPISuite) TestCreateWithNetworkAttachmentsMissing() {
	suite.Log().Info("Namespace:" + suite.Namespace)
}

func (suite *NovaAPISuite) TestCreateWithNetworkAttachmentsIPMissing() {
	suite.Log().Info("Namespace:" + suite.Namespace)
}

func (suite *NovaAPISuite) TestCreateWithExternalEndpoints() {
	suite.Log().Info("Namespace:" + suite.Namespace)
}

func (suite *NovaAPISuite) TestReconfigureToHaveNetworkAttachment() {
	suite.Log().Info("Namespace:" + suite.Namespace)
}

func TestNovaAPISuite(t *testing.T) {
	suite.Run(t, new(NovaAPISuite))
}
