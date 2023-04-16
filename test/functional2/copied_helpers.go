package functional2

import (
	"github.com/google/uuid"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	novav1 "github.com/openstack-k8s-operators/nova-operator/api/v1beta1"
)

// Below are duplicates from functional/base_test.go but we cannot meaningfully
// import them from there as they are depend on package globals that are not
// intialized if the functiona/suite_test is not run

func (h *TestHelper) CreateSecret(name types.NamespacedName, data map[string][]byte) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name.Name,
			Namespace: name.Namespace,
		},
		Data: data,
	}
	Expect(h.k8sClient.Create(h.ctx, secret)).Should(Succeed())
	return secret
}

func (h *TestHelper) CreateNovaMessageBusSecret(namespace string, name string) *corev1.Secret {
	return h.CreateSecret(
		types.NamespacedName{Namespace: namespace, Name: name},
		map[string][]byte{
			"transport_url": []byte("rabbit://fake"),
		},
	)
}

func (h *TestHelper) CreateUnstructured(rawObj map[string]interface{}) *unstructured.Unstructured {
	h.logger.Info("Creating", "raw", rawObj)
	unstructuredObj := &unstructured.Unstructured{Object: rawObj}
	_, err := controllerutil.CreateOrPatch(
		h.ctx, h.k8sClient, unstructuredObj, func() error { return nil })
	Expect(err).ShouldNot(HaveOccurred())
	return unstructuredObj
}

func (h *TestHelper) GetDefaultNovaAPISpec() map[string]interface{} {
	return map[string]interface{}{
		"secret":                  SecretName,
		"apiDatabaseHostname":     "nova-api-db-hostname",
		"apiMessageBusSecretName": MessageBusSecretName,
		"cell0DatabaseHostname":   "nova-cell0-db-hostname",
		"keystoneAuthURL":         "keystone-auth-url",
		"containerImage":          ContainerImage,
	}
}

func (h *TestHelper) CreateNovaAPI(namespace string, spec map[string]interface{}) client.Object {
	novaAPIName := uuid.New().String()

	raw := map[string]interface{}{
		"apiVersion": "nova.openstack.org/v1beta1",
		"kind":       "NovaAPI",
		"metadata": map[string]interface{}{
			"name":      novaAPIName,
			"namespace": namespace,
		},
		"spec": spec,
	}
	return h.CreateUnstructured(raw)

}

func (h *TestHelper) GetNovaAPI(name types.NamespacedName) *novav1.NovaAPI {
	instance := &novav1.NovaAPI{}
	Eventually(func(g Gomega) {
		g.Expect(h.k8sClient.Get(h.ctx, name, instance)).Should(Succeed())
	}, timeout, interval).Should(Succeed())
	return instance
}

func (h *TestHelper) DeleteInstance(instance client.Object) {
	// We have to wait for the controller to fully delete the instance
	h.logger.Info("Deleting", "Name", instance.GetName(), "Namespace", instance.GetNamespace(), "Kind", instance.GetObjectKind().GroupVersionKind().Kind)
	Eventually(func(g Gomega) {
		name := types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()}
		err := h.k8sClient.Get(h.ctx, name, instance)
		// if it is already gone that is OK
		if k8s_errors.IsNotFound(err) {
			return
		}
		g.Expect(err).ShouldNot(HaveOccurred())

		g.Expect(h.k8sClient.Delete(h.ctx, instance)).Should(Succeed())

		err = h.k8sClient.Get(h.ctx, name, instance)
		g.Expect(k8s_errors.IsNotFound(err)).To(BeTrue())
	}, timeout, interval).Should(Succeed())
}

func (h *TestHelper) NovaAPIConditionGetter(name types.NamespacedName) condition.Conditions {
	instance := h.GetNovaAPI(name)
	return instance.Status.Conditions
}
