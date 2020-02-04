package k8s

import (
	"context"
	"github.com/fpetkovski/newrelic-operator/pkg/apis/io/v1alpha1"
	"github.com/go-logr/logr"
	client_go "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Client struct {
	logr   logr.Logger
	client client_go.Client
}

func NewClient(logr logr.Logger, client client_go.Client) *Client {
	return &Client{
		logr:   logr,
		client: client,
	}
}

func (c *Client) GetPolicy(request reconcile.Request) (*v1alpha1.AlertPolicy, error) {
	var instance v1alpha1.AlertPolicy
	err := c.client.Get(context.TODO(), request.NamespacedName, &instance)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

func (c *Client) DeletePolicy(policy v1alpha1.AlertPolicy) error {
	policy.ObjectMeta.Finalizers = []string{}
	err := c.client.Update(context.TODO(), &policy)
	if err != nil {
		c.logr.Error(err, "Error updating resource")
		return err
	}
	return nil
}

func (c *Client) UpdatePolicy(policy v1alpha1.AlertPolicy) error {
	err := c.client.Status().Update(context.TODO(), &policy)
	if err != nil {
		c.logr.Error(err, "Error updating status")
		return err
	}

	policy.ObjectMeta.Finalizers = []string{"newrelic"}
	err = c.client.Update(context.TODO(), &policy)
	if err != nil {
		c.logr.Error(err, "Error updating status")
		return err
	}

	return nil
}
