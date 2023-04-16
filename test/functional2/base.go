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

package functional2

import (
	"context"
	"path/filepath"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/openstack-k8s-operators/lib-common/modules/test"
	"github.com/openstack-k8s-operators/nova-operator/controllers"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zapcore"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	ctrl "sigs.k8s.io/controller-runtime"
	ctrl_log "sigs.k8s.io/controller-runtime/pkg/log"

	networkv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	rabbitmqv1 "github.com/openstack-k8s-operators/infra-operator/apis/rabbitmq/v1beta1"
	keystonev1 "github.com/openstack-k8s-operators/keystone-operator/api/v1beta1"
	mariadbv1 "github.com/openstack-k8s-operators/mariadb-operator/api/v1beta1"
	novav1beta1 "github.com/openstack-k8s-operators/nova-operator/api/v1beta1"
	aee "github.com/openstack-k8s-operators/openstack-ansibleee-operator/api/v1alpha1"
)

const (
	timeout = 10 * time.Second
	// have maximum 100 retries before the timeout hits
	interval = timeout / 100
	// consistencyTimeout is the amount of time we use to repeatedly check
	// that a condition is still valid. This is intended to be used in
	// asserts using `Consistently`.
	consistencyTimeout = timeout
)

type BaseSuite struct {
	suite.Suite
}

func (suite *BaseSuite) DeferCleanup(fn func()) {
	// TOOD(gibi): Propose this to testify upstream
	suite.T().Cleanup(fn)
}

type EnvTestSuite struct {
	BaseSuite
	*TestHelper

	cfg       *rest.Config
	K8sClient client.Client
	testEnv   *envtest.Environment
	Ctx       context.Context
	cancel    context.CancelFunc
	suiteLog  logr.Logger

	Namespace string
}

func (suite *EnvTestSuite) setupEnvtest() {
	ctrl_log.SetLogger(zap.New(zap.UseDevMode(true), func(o *zap.Options) {
		o.Development = true
		o.TimeEncoder = zapcore.ISO8601TimeEncoder
	}))

	suite.Ctx, suite.cancel = context.WithCancel(context.TODO())

	const gomod = "../../../go.mod"

	keystoneCRDs, err := test.GetCRDDirFromModule(
		"github.com/openstack-k8s-operators/keystone-operator/api", gomod, "bases")
	suite.Assert().NoError(err)
	mariadbCRDs, err := test.GetCRDDirFromModule(
		"github.com/openstack-k8s-operators/mariadb-operator/api", gomod, "bases")
	suite.Assert().NoError(err)
	rabbitCRDs, err := test.GetCRDDirFromModule(
		"github.com/openstack-k8s-operators/infra-operator/apis", gomod, "bases")
	suite.Assert().NoError(err)
	routev1CRDs, err := test.GetOpenShiftCRDDir("route/v1", gomod)
	suite.Assert().NoError(err)
	// NOTE(gibi): there are packages where the CRD directory has other
	// yamls files as well, then we need to specify the extac file to load
	networkv1CRD, err := test.GetCRDDirFromModule(
		"github.com/k8snetworkplumbingwg/network-attachment-definition-client", gomod, "artifacts/networks-crd.yaml")
	suite.Assert().NoError(err)
	aeeCRDs, err := test.GetCRDDirFromModule(
		"github.com/openstack-k8s-operators/openstack-ansibleee-operator/api", gomod, "bases")
	suite.Assert().NoError(err)

	//	By("bootstrapping test environment")
	suite.testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "config", "crd", "bases"),
			// NOTE(gibi): we need to list all the external CRDs our operator depends on
			mariadbCRDs,
			keystoneCRDs,
			rabbitCRDs,
			routev1CRDs,
			aeeCRDs,
		},
		CRDInstallOptions: envtest.CRDInstallOptions{
			Paths: []string{
				networkv1CRD,
			},
		},
		ErrorIfCRDPathMissing: true,
	}

	// cfg is defined in this file globally.
	suite.cfg, err = suite.testEnv.Start()
	suite.Assert().NoError(err)
	suite.Assert().NotNil(suite.cfg)

	// NOTE(gibi): Need to add all API schemas our operator can own.
	// this includes external scheme lke mariadb otherwise the
	// reconciler loop will silently not start
	// TODO(sean): factor this out to a common function.
	err = novav1beta1.AddToScheme(scheme.Scheme)
	suite.Assert().NoError(err)
	err = mariadbv1.AddToScheme(scheme.Scheme)
	suite.Assert().NoError(err)
	err = keystonev1.AddToScheme(scheme.Scheme)
	suite.Assert().NoError(err)
	err = corev1.AddToScheme(scheme.Scheme)
	suite.Assert().NoError(err)
	err = appsv1.AddToScheme(scheme.Scheme)
	suite.Assert().NoError(err)
	err = routev1.AddToScheme(scheme.Scheme)
	suite.Assert().NoError(err)
	err = rabbitmqv1.AddToScheme(scheme.Scheme)
	suite.Assert().NoError(err)
	err = networkv1.AddToScheme(scheme.Scheme)
	suite.Assert().NoError(err)
	err = aee.AddToScheme(scheme.Scheme)
	suite.Assert().NoError(err)

	//+kubebuilder:scaffold:scheme

	suite.suiteLog = ctrl.Log.WithName("---Test")

	suite.K8sClient, err = client.New(suite.cfg, client.Options{Scheme: scheme.Scheme})
	suite.Assert().NoError(err)
	suite.Assert().NotNil(suite.K8sClient)
	suite.TestHelper = NewTestHelper(suite.Ctx, suite.K8sClient, timeout, interval, suite.suiteLog)
	suite.Assert().NotNil(suite.TestHelper)

	// Start the controller-manager in a goroutine
	k8sManager, err := ctrl.NewManager(suite.cfg, ctrl.Options{
		Scheme: scheme.Scheme,
		// NOTE(gibi): disable metrics reporting in test to allow
		// parallel test execution. Otherwise each instance would like to
		// bind to the same port
		MetricsBindAddress: "0",
	})
	suite.Assert().NoError(err)

	kclient, err := kubernetes.NewForConfig(suite.cfg)
	suite.Assert().NoError(err)

	reconcilers := controllers.NewReconcilers(k8sManager, kclient)
	// NOTE(gibi): During envtest we simulate success of tasks (e.g Job,
	// Deployment, DB) so we can speed up the test execution by reducing the
	// time we wait before we reconcile when a task is running.
	reconcilers.OverrideRequeueTimeout(time.Duration(10) * time.Millisecond)
	err = reconcilers.Setup(k8sManager, ctrl.Log.WithName("testSetup"))
	suite.Assert().NoError(err)

	go func() {
		//		defer GinkgoRecover()
		err = k8sManager.Start(suite.Ctx)
		suite.Assert().NoError(err)
	}()

}

func (suite *EnvTestSuite) tearDownEnvTest() {
	suite.cancel()
	err := suite.testEnv.Stop()
	suite.Assert().NoError(err)
}

func (suite *EnvTestSuite) SetupSuite() {
	// NOTE(gibi) To support Gomega we need to define for Gomega how to
	// report an assertion failure to the test library.
	gomega.RegisterFailHandler(func(message string, callerSkip ...int) {
		suite.Assert().Fail(message)
	})
	suite.setupEnvtest()
}

func (suite *EnvTestSuite) TearDownSuite() {
	suite.tearDownEnvTest()
}

func (suite *EnvTestSuite) SetupTest() {
	suite.T().Log("SetupTest")
	// NOTE(gibi): We need to create a unique namespace for each test run
	// as namespaces cannot be deleted in a locally running envtest. See
	// https://book.kubebuilder.io/reference/envtest.html#namespace-usage-limitation
	suite.Namespace = uuid.New().String()
	suite.CreateNamespace(suite.Namespace)
	// We still request the delete of the Namespace to properly cleanup if
	// we run the test in an existing cluster.
	suite.DeferCleanup(func() { suite.DeleteNamespace(suite.Namespace) })
}

func (suite *EnvTestSuite) Log() logr.Logger {
	return suite.suiteLog.WithName(suite.T().Name())
}
