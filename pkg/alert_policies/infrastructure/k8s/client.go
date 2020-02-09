package k8s

import (
	"context"
	"github.com/fpetkovski/newrelic-operator/pkg/apis/io/v1alpha1"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	client_go "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Client struct {
	client client_go.Client
}

func NewClient(client client_go.Client) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) GetPolicy(ctx context.Context, request reconcile.Request) (*v1alpha1.AlertPolicy, error) {
	sp, _ := opentracing.StartSpanFromContext(ctx, "get-alert-policy-k8s")
	defer sp.Finish()

	var instance v1alpha1.AlertPolicy
	err := c.client.Get(context.TODO(), request.NamespacedName, &instance)
	if err != nil {
		sp.LogFields(log.Error(err))
		return nil, err
	}

	return &instance, nil
}

func (c *Client) DeletePolicy(ctx context.Context, policy v1alpha1.AlertPolicy) error {
	sp, _ := opentracing.StartSpanFromContext(ctx, "delete-alert-policy-k8s")
	defer sp.Finish()

	policy.ObjectMeta.Finalizers = []string{}
	err := c.client.Update(context.TODO(), &policy)
	if err != nil {
		sp.LogFields(log.Error(err))

		return err
	}
	return nil
}

func (c *Client) UpdatePolicy(ctx context.Context, policy v1alpha1.AlertPolicy) error {
	sp, _ := opentracing.StartSpanFromContext(ctx, "update-alert-policy-k8s")
	defer sp.Finish()

	err := c.client.Status().Update(context.TODO(), &policy)
	if err != nil {
		sp.LogFields(log.Error(err))
		return err
	}

	policy.ObjectMeta.Finalizers = []string{"newrelic"}
	err = c.client.Update(context.TODO(), &policy)
	if err != nil {
		sp.LogFields(log.Error(err))
		return err
	}

	return nil
}
