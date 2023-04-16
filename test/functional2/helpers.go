package functional2

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	//. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	common_helpers "github.com/openstack-k8s-operators/lib-common/modules/test/helpers"
)

const (
	SecretName           = "test-secret"
	MessageBusSecretName = "rabbitmq-secret"
	ContainerImage       = "test://nova"
)

type TestHelper struct {
	common_helpers.TestHelper

	// TODO(gibi): make these accessible via the common_helpers.TestHelper
	// so that we don't have to store them here as well
	k8sClient client.Client
	ctx       context.Context
	timeout   time.Duration
	interval  time.Duration
	logger    logr.Logger
}

func NewTestHelper(
	ctx context.Context,
	k8sClient client.Client,
	timeout time.Duration,
	interval time.Duration,
	logger logr.Logger,
) *TestHelper {
	return &TestHelper{
		TestHelper: *common_helpers.NewTestHelper(ctx, k8sClient, timeout, interval, logger),
		ctx:        ctx,
		k8sClient:  k8sClient,
		timeout:    timeout,
		interval:   interval,
		logger:     logger,
	}
}
